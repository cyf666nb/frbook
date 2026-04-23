package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bookshare/internal/config"
	"bookshare/internal/database"
	"bookshare/internal/model"
	"bookshare/internal/repository"
)

type RentalService struct {
	orderRepo   *repository.RentalOrderRepository
	bookRepo    *repository.BookRepository
	userRepo    *repository.UserRepository
	walletRepo  *repository.WalletRepository
	txRepo      *repository.TransactionRepository
	redis       *database.RedisClient
	platform    *config.PlatformConfig
}

func NewRentalService(
	orderRepo *repository.RentalOrderRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	walletRepo *repository.WalletRepository,
	txRepo *repository.TransactionRepository,
	redis *database.RedisClient,
	platform *config.PlatformConfig,
) *RentalService {
	return &RentalService{
		orderRepo:  orderRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
		walletRepo: walletRepo,
		txRepo:     txRepo,
		redis:      redis,
		platform:   platform,
	}
}

func (s *RentalService) CreateOrder(ctx context.Context, renterID uint64, req *model.CreateRentalRequest) (*model.RentalOrder, error) {
	book, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		return nil, errors.New("图书不存在")
	}

	if book.Mode != model.BookModeRental {
		return nil, errors.New("该图书不支持租赁")
	}

	if book.Status != model.BookStatusOnline {
		return nil, errors.New("图书已下架")
	}

	if book.UserID == renterID {
		return nil, errors.New("不能租赁自己的图书")
	}

	if req.RentDays < book.MinRentDays {
		return nil, fmt.Errorf("最短租期为%d天", book.MinRentDays)
	}

	lockKey := fmt.Sprintf("lock:rental:book:%d", book.ID)
	orderNo := generateOrderNo("RN")
	locked, err := s.redis.Lock(ctx, lockKey, orderNo, 10*time.Minute)
	if err != nil || !locked {
		return nil, errors.New("图书已被其他人锁定,请稍后重试")
	}

	totalRent := *book.DailyRent * float64(req.RentDays)
	totalAmount := *book.Deposit + totalRent

	order := &model.RentalOrder{
		OrderNo:     orderNo,
		BookID:      book.ID,
		BookTitle:   book.Title,
		OwnerID:     book.UserID,
		RenterID:    renterID,
		DailyRent:   *book.DailyRent,
		Deposit:     *book.Deposit,
		RentDays:    req.RentDays,
		TotalRent:   totalRent,
		TotalAmount: totalAmount,
		Status:      model.RentalStatusPendingPayment,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		s.redis.Unlock(ctx, lockKey, orderNo)
		return nil, err
	}

	newLockKey := fmt.Sprintf("lock:rental:pay:%d", order.ID)
	newLockValue := fmt.Sprintf("%d", order.ID)
	payLocked, _ := s.redis.Lock(ctx, newLockKey, newLockValue, s.platform.PayTimeout)
	if !payLocked {
		_ = s.orderRepo.UpdateStatus(ctx, order.ID, model.RentalStatusCancelled)
		s.redis.Unlock(ctx, lockKey, orderNo)
		return nil, errors.New("锁定失败")
	}

	s.redis.Unlock(ctx, lockKey, orderNo)
	return order, nil
}

func (s *RentalService) PayOrder(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.RenterID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.RentalStatusPendingPayment {
		return errors.New("订单状态错误")
	}

	wallet, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet.Balance < order.TotalAmount {
		return errors.New("余额不足")
	}

	if err := s.walletRepo.Freeze(ctx, userID, order.TotalAmount); err != nil {
		return err
	}

	tx := &model.Transaction{
		TransactionNo: generateOrderNo("TX"),
		UserID:        userID,
		RelatedUserID: &order.OwnerID,
		Type:          model.TransactionTypeRentalDeposit,
		OrderType:     model.TransactionTypeOrderRental,
		OrderID:       order.ID,
		OrderNo:       order.OrderNo,
		Amount:        -order.TotalAmount,
		Status:        1,
		Remark:        "租赁支付押金和租金",
	}
	_ = s.txRepo.Create(ctx, tx)

	order.Status = model.RentalStatusPaidWaitingPick
	return s.orderRepo.Update(ctx, order)
}

func (s *RentalService) ConfirmPickup(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.OwnerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.RentalStatusPaidWaitingPick {
		return errors.New("订单状态错误")
	}

	now := time.Now()
	endTime := now.AddDate(0, 0, order.RentDays)
	order.StartTime = &now
	order.EndTime = &endTime
	order.Status = model.RentalStatusRenting

	countdownKey := fmt.Sprintf("rental:countdown:%d", order.ID)
	_ = s.redis.Set(ctx, countdownKey, endTime.Unix(), time.Duration(order.RentDays)*24*time.Hour)

	return s.orderRepo.Update(ctx, order)
}

func (s *RentalService) ReturnBook(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.RenterID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.RentalStatusRenting && order.Status != model.RentalStatusOverdue {
		return errors.New("订单状态错误")
	}

	now := time.Now()
	order.ActualReturnTime = &now
	order.Status = model.RentalStatusWaitingInspection

	if order.Status == model.RentalStatusOverdue && order.EndTime != nil {
		overdueDays := int(now.Sub(*order.EndTime).Hours() / 24)
		if overdueDays > 0 {
			order.OverdueFee = order.DailyRent * 1.5 * float64(overdueDays)
		}
	}

	return s.orderRepo.Update(ctx, order)
}

func (s *RentalService) Inspect(ctx context.Context, orderID, userID uint64, req *model.InspectRentalRequest) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.OwnerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.RentalStatusWaitingInspection {
		return errors.New("订单状态错误")
	}

	refundAmount := order.Deposit - order.OverdueFee
	if refundAmount < 0 {
		refundAmount = 0
	}

	settledAmount := order.TotalRent

	order.RefundAmount = &refundAmount
	order.SettledAmount = &settledAmount
	order.Status = model.RentalStatusCompleted

	if err := s.walletRepo.Unfreeze(ctx, order.RenterID, order.TotalAmount); err != nil {
		return err
	}

	if refundAmount > 0 {
		if err := s.walletRepo.AddBalance(ctx, order.RenterID, refundAmount); err != nil {
			return err
		}
		tx := &model.Transaction{
			TransactionNo: generateOrderNo("TX"),
			UserID:        order.RenterID,
			RelatedUserID: &order.OwnerID,
			Type:          model.TransactionTypeRentalDepositRefund,
			OrderType:     model.TransactionTypeOrderRental,
			OrderID:       order.ID,
			OrderNo:       order.OrderNo,
			Amount:        refundAmount,
			Status:        1,
			Remark:        "租赁押金退还",
		}
		_ = s.txRepo.Create(ctx, tx)
	}

	if err := s.walletRepo.AddBalance(ctx, order.OwnerID, settledAmount); err != nil {
		return err
	}
	tx := &model.Transaction{
		TransactionNo: generateOrderNo("TX"),
		UserID:        order.OwnerID,
		RelatedUserID: &order.RenterID,
		Type:          model.TransactionTypeRentalRent,
		OrderType:     model.TransactionTypeOrderRental,
		OrderID:       order.ID,
		OrderNo:       order.OrderNo,
		Amount:        settledAmount,
		Status:        1,
		Remark:        "租赁收入",
	}
	_ = s.txRepo.Create(ctx, tx)

	lockKey := fmt.Sprintf("lock:rental:book:%d", order.BookID)
	_ = s.redis.Unlock(ctx, lockKey, fmt.Sprintf("%d", order.ID))

	return s.orderRepo.Update(ctx, order)
}

func (s *RentalService) CancelOrder(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.RenterID != userID && order.OwnerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.RentalStatusPendingPayment && order.Status != model.RentalStatusPaidWaitingPick {
		return errors.New("订单状态错误,无法取消")
	}

	if order.Status == model.RentalStatusPaidWaitingPick {
		if err := s.walletRepo.Unfreeze(ctx, order.RenterID, order.TotalAmount); err != nil {
			return err
		}
		if err := s.walletRepo.AddBalance(ctx, order.RenterID, order.TotalAmount); err != nil {
			return err
		}
	}

	order.Status = model.RentalStatusCancelled

	lockKey := fmt.Sprintf("lock:rental:book:%d", order.BookID)
	_ = s.redis.Unlock(ctx, lockKey, fmt.Sprintf("%d", order.ID))

	return s.orderRepo.Update(ctx, order)
}

func (s *RentalService) GetByID(ctx context.Context, orderID uint64) (*model.RentalOrder, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

func (s *RentalService) ListByUserID(ctx context.Context, userID uint64, query *model.RentalListQuery) ([]model.RentalOrder, int64, error) {
	return s.orderRepo.ListByUserID(ctx, userID, query.Role, query.Status, query.Page, query.PageSize)
}

func (s *RentalService) Rate(ctx context.Context, orderID, userID uint64, req *model.RateRequest) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != model.RentalStatusCompleted {
		return errors.New("订单未完成")
	}

	if userID == order.OwnerID {
		order.RenterRating = &req.Rating
		order.RenterComment = req.Comment
	} else if userID == order.RenterID {
		order.OwnerRating = &req.Rating
		order.OwnerComment = req.Comment
	} else {
		return errors.New("无权限评价")
	}

	return s.orderRepo.Update(ctx, order)
}

func (s *RentalService) ProcessOverdueOrders(ctx context.Context) error {
	orders, err := s.orderRepo.GetOverdueOrders(ctx)
	if err != nil {
		return err
	}

	for i := range orders {
		orders[i].Status = model.RentalStatusOverdue
		_ = s.orderRepo.Update(ctx, &orders[i])
	}

	return nil
}
