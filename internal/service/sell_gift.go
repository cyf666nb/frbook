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

type SellService struct {
	orderRepo   *repository.SellOrderRepository
	bookRepo    *repository.BookRepository
	userRepo    *repository.UserRepository
	walletRepo  *repository.WalletRepository
	txRepo      *repository.TransactionRepository
	redis       *database.RedisClient
	platform    *config.PlatformConfig
}

func NewSellService(
	orderRepo *repository.SellOrderRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	walletRepo *repository.WalletRepository,
	txRepo *repository.TransactionRepository,
	redis *database.RedisClient,
	platform *config.PlatformConfig,
) *SellService {
	return &SellService{
		orderRepo:  orderRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
		walletRepo: walletRepo,
		txRepo:     txRepo,
		redis:      redis,
		platform:   platform,
	}
}

func (s *SellService) CreateOrder(ctx context.Context, buyerID uint64, req *model.CreateSellRequest) (*model.SellOrder, error) {
	book, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		return nil, errors.New("图书不存在")
	}

	if book.Mode != model.BookModeSell {
		return nil, errors.New("该图书不支持购买")
	}

	if book.Status != model.BookStatusOnline {
		return nil, errors.New("图书已下架")
	}

	if book.UserID == buyerID {
		return nil, errors.New("不能购买自己的图书")
	}

	lockKey := fmt.Sprintf("lock:sell:book:%d", book.ID)
	orderNo := generateOrderNo("SL")
	locked, err := s.redis.Lock(ctx, lockKey, orderNo, 15*time.Minute)
	if err != nil || !locked {
		return nil, errors.New("图书已被其他人锁定,请稍后重试")
	}

	defer func() {
		s.redis.Unlock(ctx, lockKey, orderNo)
	}()

	platformFee := *book.SellPrice * s.platform.FeeRate

	order := &model.SellOrder{
		OrderNo:      orderNo,
		BookID:       book.ID,
		BookTitle:    book.Title,
		SellerID:     book.UserID,
		BuyerID:      buyerID,
		Price:        *book.SellPrice,
		PlatformFee:  platformFee,
		Status:       model.SellStatusPendingPayment,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	locked, _ = s.redis.Lock(ctx, lockKey, fmt.Sprintf("%d", order.ID), s.platform.PayTimeout)
	if !locked {
		_ = s.orderRepo.UpdateStatus(ctx, order.ID, model.SellStatusCancelled)
		return nil, errors.New("锁定失败")
	}

	return order, nil
}

func (s *SellService) PayOrder(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.SellStatusPendingPayment {
		return errors.New("订单状态错误")
	}

	wallet, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet.Balance < order.Price {
		return errors.New("余额不足")
	}

	if err := s.walletRepo.Freeze(ctx, userID, order.Price); err != nil {
		return err
	}

	tx := &model.Transaction{
		TransactionNo: generateOrderNo("TX"),
		UserID:        userID,
		RelatedUserID: &order.SellerID,
		Type:          model.TransactionTypeSellPayment,
		OrderType:     model.TransactionTypeOrderSell,
		OrderID:       order.ID,
		OrderNo:       order.OrderNo,
		Amount:        -order.Price,
		Status:        1,
		Remark:        "购买图书支付",
	}
	_ = s.txRepo.Create(ctx, tx)

	order.Status = model.SellStatusPaidWaitingShip
	return s.orderRepo.Update(ctx, order)
}

func (s *SellService) Ship(ctx context.Context, orderID, userID uint64, req *model.ShipSellRequest) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.SellerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.SellStatusPaidWaitingShip {
		return errors.New("订单状态错误")
	}

	now := time.Now()
	autoSettleTime := now.AddDate(0, 0, s.platform.AutoSettleDay)

	order.DeliveryType = &req.DeliveryType
	order.DeliveryNo = req.DeliveryNo
	order.DeliveryCompany = req.DeliveryCompany
	order.DeliveryTime = &now
	order.AutoSettleTime = &autoSettleTime
	order.Status = model.SellStatusShipped

	delayKey := fmt.Sprintf("sell:settle:delay:%d", order.ID)
	_ = s.redis.Set(ctx, delayKey, autoSettleTime.Unix(), time.Duration(s.platform.AutoSettleDay)*24*time.Hour)

	return s.orderRepo.Update(ctx, order)
}

func (s *SellService) ConfirmPickup(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.SellerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.SellStatusPaidWaitingShip {
		return errors.New("订单状态错误")
	}

	now := time.Now()
	autoSettleTime := now.AddDate(0, 0, s.platform.AutoSettleDay)

	deliveryType := int8(model.DeliveryTypePickup)
	order.DeliveryType = &deliveryType
	order.DeliveryTime = &now
	order.AutoSettleTime = &autoSettleTime
	order.Status = model.SellStatusShipped

	return s.orderRepo.Update(ctx, order)
}

func (s *SellService) ConfirmReceive(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.SellStatusShipped {
		return errors.New("订单状态错误")
	}

	return s.settleOrder(ctx, order)
}

func (s *SellService) settleOrder(ctx context.Context, order *model.SellOrder) error {
	if err := s.walletRepo.DeductFrozen(ctx, order.BuyerID, order.Price); err != nil {
		return err
	}

	settleAmount := order.Price - order.PlatformFee
	order.SettleAmount = &settleAmount

	if err := s.walletRepo.AddBalance(ctx, order.SellerID, settleAmount); err != nil {
		return err
	}

	tx := &model.Transaction{
		TransactionNo: generateOrderNo("TX"),
		UserID:        order.SellerID,
		RelatedUserID: &order.BuyerID,
		Type:          model.TransactionTypeSellSettlement,
		OrderType:     model.TransactionTypeOrderSell,
		OrderID:       order.ID,
		OrderNo:       order.OrderNo,
		Amount:        settleAmount,
		Status:        1,
		Remark:        "出售收入",
	}
	_ = s.txRepo.Create(ctx, tx)

	now := time.Now()
	order.ReceiveTime = &now
	order.Status = model.SellStatusCompleted

	if err := s.bookRepo.UpdateStatus(ctx, order.BookID, model.BookStatusTraded); err != nil {
		return err
	}

	lockKey := fmt.Sprintf("lock:sell:book:%d", order.BookID)
	_ = s.redis.Unlock(ctx, lockKey, fmt.Sprintf("%d", order.ID))

	return s.orderRepo.Update(ctx, order)
}

func (s *SellService) CancelOrder(ctx context.Context, orderID, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return errors.New("无权限操作")
	}

	if order.Status != model.SellStatusPendingPayment && order.Status != model.SellStatusPaidWaitingShip {
		return errors.New("订单状态错误,无法取消")
	}

	if order.Status == model.SellStatusPaidWaitingShip {
		if err := s.walletRepo.Unfreeze(ctx, order.BuyerID, order.Price); err != nil {
			return err
		}
		if err := s.walletRepo.AddBalance(ctx, order.BuyerID, order.Price); err != nil {
			return err
		}
	}

	order.Status = model.SellStatusCancelled

	lockKey := fmt.Sprintf("lock:sell:book:%d", order.BookID)
	_ = s.redis.Unlock(ctx, lockKey, fmt.Sprintf("%d", order.ID))

	return s.orderRepo.Update(ctx, order)
}

func (s *SellService) GetByID(ctx context.Context, orderID uint64) (*model.SellOrder, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

func (s *SellService) ListByUserID(ctx context.Context, userID uint64, query *model.SellListQuery) ([]model.SellOrder, int64, error) {
	return s.orderRepo.ListByUserID(ctx, userID, query.Role, query.Status, query.Page, query.PageSize)
}

func (s *SellService) Rate(ctx context.Context, orderID, userID uint64, req *model.RateRequest) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != model.SellStatusCompleted {
		return errors.New("订单未完成")
	}

	if userID == order.SellerID {
		order.BuyerRating = &req.Rating
		order.BuyerComment = req.Comment
	} else if userID == order.BuyerID {
		order.SellerRating = &req.Rating
		order.SellerComment = req.Comment
	} else {
		return errors.New("无权限评价")
	}

	return s.orderRepo.Update(ctx, order)
}

func (s *SellService) ProcessAutoSettle(ctx context.Context) error {
	orders, err := s.orderRepo.GetPendingAutoSettle(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		_ = s.settleOrder(ctx, &order)
	}

	return nil
}

type GiftService struct {
	recordRepo *repository.GiftRecordRepository
	bookRepo   *repository.BookRepository
	userRepo   *repository.UserRepository
	redis      *database.RedisClient
}

func NewGiftService(
	recordRepo *repository.GiftRecordRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	redis *database.RedisClient,
) *GiftService {
	return &GiftService{
		recordRepo: recordRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
		redis:      redis,
	}
}

func (s *GiftService) CreateRecord(ctx context.Context, receiverID uint64, req *model.CreateGiftRequest) (*model.GiftRecord, error) {
	book, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		return nil, errors.New("图书不存在")
	}

	if book.Mode != model.BookModeGift {
		return nil, errors.New("该图书不支持赠送")
	}

	if book.Status != model.BookStatusOnline {
		return nil, errors.New("图书已下架")
	}

	if book.UserID == receiverID {
		return nil, errors.New("不能领取自己的图书")
	}

	lockKey := fmt.Sprintf("lock:gift:book:%d", book.ID)
	recordNo := generateOrderNo("GF")
	locked, err := s.redis.Lock(ctx, lockKey, recordNo, 5*time.Minute)
	if err != nil || !locked {
		return nil, errors.New("图书已被其他人锁定,请稍后重试")
	}

	defer func() {
		s.redis.Unlock(ctx, lockKey, recordNo)
	}()

	record := &model.GiftRecord{
		RecordNo:   recordNo,
		BookID:     book.ID,
		BookTitle:  book.Title,
		GiverID:    book.UserID,
		ReceiverID: receiverID,
		Status:     model.GiftStatusPending,
	}

	if err := s.recordRepo.Create(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *GiftService) Confirm(ctx context.Context, recordID, userID uint64) error {
	record, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return errors.New("记录不存在")
	}

	if record.GiverID != userID {
		return errors.New("无权限操作")
	}

	if record.Status != model.GiftStatusPending {
		return errors.New("状态错误")
	}

	record.Status = model.GiftStatusConfirmed
	return s.recordRepo.Update(ctx, record)
}

func (s *GiftService) Reject(ctx context.Context, recordID, userID uint64) error {
	record, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return errors.New("记录不存在")
	}

	if record.GiverID != userID {
		return errors.New("无权限操作")
	}

	if record.Status != model.GiftStatusPending {
		return errors.New("状态错误")
	}

	record.Status = model.GiftStatusCancelled
	return s.recordRepo.Update(ctx, record)
}

func (s *GiftService) Cancel(ctx context.Context, recordID, userID uint64) error {
	record, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return errors.New("记录不存在")
	}

	if record.ReceiverID != userID {
		return errors.New("无权限操作")
	}

	if record.Status != model.GiftStatusPending && record.Status != model.GiftStatusConfirmed {
		return errors.New("状态错误")
	}

	record.Status = model.GiftStatusCancelled
	return s.recordRepo.Update(ctx, record)
}

func (s *GiftService) Deliver(ctx context.Context, recordID, userID uint64) error {
	record, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return errors.New("记录不存在")
	}

	if record.GiverID != userID && record.ReceiverID != userID {
		return errors.New("无权限操作")
	}

	if record.Status != model.GiftStatusConfirmed {
		return errors.New("状态错误")
	}

	now := time.Now()
	record.DeliveryTime = &now
	record.Status = model.GiftStatusCompleted

	if err := s.bookRepo.UpdateStatus(ctx, record.BookID, model.BookStatusTraded); err != nil {
		return err
	}

	return s.recordRepo.Update(ctx, record)
}

func (s *GiftService) GetByID(ctx context.Context, recordID uint64) (*model.GiftRecord, error) {
	return s.recordRepo.FindByID(ctx, recordID)
}

func (s *GiftService) ListByUserID(ctx context.Context, userID uint64, query *model.GiftListQuery) ([]model.GiftRecord, int64, error) {
	return s.recordRepo.ListByUserID(ctx, userID, query.Role, query.Status, query.Page, query.PageSize)
}

func (s *GiftService) GetApplicationsByBookID(ctx context.Context, bookID, userID uint64) ([]model.GiftRecord, error) {
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return nil, errors.New("图书不存在")
	}

	if book.UserID != userID {
		return nil, errors.New("无权限查看")
	}

	return s.recordRepo.GetPendingByBookID(ctx, bookID)
}

func (s *GiftService) Rate(ctx context.Context, recordID, userID uint64, req *model.RateRequest) error {
	record, err := s.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return errors.New("记录不存在")
	}

	if record.Status != model.GiftStatusCompleted {
		return errors.New("未完成")
	}

	if userID == record.GiverID {
		record.ReceiverRating = &req.Rating
		record.ReceiverComment = req.Comment
	} else if userID == record.ReceiverID {
		record.GiverRating = &req.Rating
		record.GiverComment = req.Comment
	} else {
		return errors.New("无权限评价")
	}

	return s.recordRepo.Update(ctx, record)
}
