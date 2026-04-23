package repository

import (
	"context"
	"errors"
	"time"

	"bookshare/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("account = ?", account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) UpdateCreditScore(ctx context.Context, userID uint64, scoreChange int, reason string, orderType *int8, orderID *uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		oldScore := user.CreditScore
		newScore := oldScore + scoreChange
		if newScore < 0 {
			newScore = 0
		}
		if newScore > 100 {
			newScore = 100
		}

		if err := tx.Model(&user).Update("credit_score", newScore).Error; err != nil {
			return err
		}

		creditRecord := &model.CreditRecord{
			UserID:          userID,
			ScoreChange:     scoreChange,
			ScoreBefore:     oldScore,
			ScoreAfter:      newScore,
			Reason:          reason,
			RelatedOrderType: orderType,
			RelatedOrderID:  orderID,
		}
		return tx.Create(creditRecord).Error
	})
}

func (r *UserRepository) ExistsByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("account = ?", account).Count(&count).Error
	return count > 0, err
}

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(ctx context.Context, book *model.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *BookRepository) FindByID(ctx context.Context, id uint64) (*model.Book, error) {
	var book model.Book
	err := r.db.WithContext(ctx).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindByIDWithUser(ctx context.Context, id uint64) (*model.Book, error) {
	var book model.Book
	err := r.db.WithContext(ctx).Preload("User").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Update(ctx context.Context, book *model.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *BookRepository) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.Book{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Book{}, id).Error
}

func (r *BookRepository) IncrementViewCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Book{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *BookRepository) List(ctx context.Context, query *model.BookListQuery) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Book{}).Where("status = ?", model.BookStatusOnline)

	if query.Mode > 0 {
		db = db.Where("mode = ?", query.Mode)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR author LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.School != "" {
		db = db.Joins("JOIN users ON users.id = books.user_id").Where("users.school = ?", query.School)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "DESC"
	if query.Order == "asc" {
		order = "ASC"
	}

	sortField := "created_at"
	switch query.Sort {
	case "view_count":
		sortField = "view_count"
	case "price":
		sortField = "COALESCE(daily_rent, sell_price, 0)"
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order(sortField + " " + order).Offset(offset).Limit(pageSize).Find(&books).Error
	return books, total, err
}

func (r *BookRepository) ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Book{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&books).Error
	return books, total, err
}

func (r *BookRepository) GetHotBooks(ctx context.Context, limit int) ([]model.Book, error) {
	var books []model.Book
	err := r.db.WithContext(ctx).Where("status = ?", model.BookStatusOnline).
		Order("view_count DESC").Limit(limit).Find(&books).Error
	return books, err
}

func (r *BookRepository) GetLatestGiftBooks(ctx context.Context, limit int) ([]model.Book, error) {
	var books []model.Book
	err := r.db.WithContext(ctx).Where("status = ? AND mode = ?", model.BookStatusOnline, model.BookModeGift).
		Order("created_at DESC").Limit(limit).Find(&books).Error
	return books, err
}

type RentalOrderRepository struct {
	db *gorm.DB
}

func NewRentalOrderRepository(db *gorm.DB) *RentalOrderRepository {
	return &RentalOrderRepository{db: db}
}

func (r *RentalOrderRepository) Create(ctx context.Context, order *model.RentalOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *RentalOrderRepository) FindByID(ctx context.Context, id uint64) (*model.RentalOrder, error) {
	var order model.RentalOrder
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *RentalOrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.RentalOrder, error) {
	var order model.RentalOrder
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *RentalOrderRepository) Update(ctx context.Context, order *model.RentalOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *RentalOrderRepository) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.RentalOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *RentalOrderRepository) ListByUserID(ctx context.Context, userID uint64, role int8, status *int8, page, pageSize int) ([]model.RentalOrder, int64, error) {
	var orders []model.RentalOrder
	var total int64

	db := r.db.WithContext(ctx).Model(&model.RentalOrder{})

	if role == 1 {
		db = db.Where("owner_id = ?", userID)
	} else if role == 2 {
		db = db.Where("renter_id = ?", userID)
	} else {
		db = db.Where("owner_id = ? OR renter_id = ?", userID, userID)
	}

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *RentalOrderRepository) GetOverdueOrders(ctx context.Context) ([]model.RentalOrder, error) {
	var orders []model.RentalOrder
	now := time.Now()
	err := r.db.WithContext(ctx).Where("status = ? AND end_time < ?", model.RentalStatusRenting, now).Find(&orders).Error
	return orders, err
}

type SellOrderRepository struct {
	db *gorm.DB
}

func NewSellOrderRepository(db *gorm.DB) *SellOrderRepository {
	return &SellOrderRepository{db: db}
}

func (r *SellOrderRepository) Create(ctx context.Context, order *model.SellOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *SellOrderRepository) FindByID(ctx context.Context, id uint64) (*model.SellOrder, error) {
	var order model.SellOrder
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *SellOrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*model.SellOrder, error) {
	var order model.SellOrder
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *SellOrderRepository) Update(ctx context.Context, order *model.SellOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *SellOrderRepository) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.SellOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *SellOrderRepository) ListByUserID(ctx context.Context, userID uint64, role int8, status *int8, page, pageSize int) ([]model.SellOrder, int64, error) {
	var orders []model.SellOrder
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SellOrder{})

	if role == 1 {
		db = db.Where("seller_id = ?", userID)
	} else if role == 2 {
		db = db.Where("buyer_id = ?", userID)
	} else {
		db = db.Where("seller_id = ? OR buyer_id = ?", userID, userID)
	}

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *SellOrderRepository) GetPendingAutoSettle(ctx context.Context) ([]model.SellOrder, error) {
	var orders []model.SellOrder
	now := time.Now()
	err := r.db.WithContext(ctx).Where("status = ? AND auto_settle_time <= ?", model.SellStatusShipped, now).Find(&orders).Error
	return orders, err
}

type GiftRecordRepository struct {
	db *gorm.DB
}

func NewGiftRecordRepository(db *gorm.DB) *GiftRecordRepository {
	return &GiftRecordRepository{db: db}
}

func (r *GiftRecordRepository) Create(ctx context.Context, record *model.GiftRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *GiftRecordRepository) FindByID(ctx context.Context, id uint64) (*model.GiftRecord, error) {
	var record model.GiftRecord
	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *GiftRecordRepository) FindByRecordNo(ctx context.Context, recordNo string) (*model.GiftRecord, error) {
	var record model.GiftRecord
	err := r.db.WithContext(ctx).Where("record_no = ?", recordNo).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *GiftRecordRepository) Update(ctx context.Context, record *model.GiftRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *GiftRecordRepository) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.GiftRecord{}).Where("id = ?", id).Update("status", status).Error
}

func (r *GiftRecordRepository) ListByUserID(ctx context.Context, userID uint64, role int8, status *int8, page, pageSize int) ([]model.GiftRecord, int64, error) {
	var records []model.GiftRecord
	var total int64

	db := r.db.WithContext(ctx).Model(&model.GiftRecord{})

	if role == 1 {
		db = db.Where("giver_id = ?", userID)
	} else if role == 2 {
		db = db.Where("receiver_id = ?", userID)
	} else {
		db = db.Where("giver_id = ? OR receiver_id = ?", userID, userID)
	}

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	return records, total, err
}

func (r *GiftRecordRepository) GetPendingByBookID(ctx context.Context, bookID uint64) ([]model.GiftRecord, error) {
	var records []model.GiftRecord
	err := r.db.WithContext(ctx).Where("book_id = ? AND status = ?", bookID, model.GiftStatusPending).Find(&records).Error
	return records, err
}

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) FindByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = model.Wallet{UserID: userID}
			if createErr := r.db.Create(&wallet).Error; createErr != nil {
				return nil, createErr
			}
			return &wallet, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Create(ctx context.Context, wallet *model.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *WalletRepository) Freeze(ctx context.Context, userID uint64, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}
		if wallet.Balance < amount {
			return errors.New("余额不足")
		}
		return tx.Model(&wallet).Updates(map[string]interface{}{
			"balance":        gorm.Expr("balance - ?", amount),
			"frozen_balance": gorm.Expr("frozen_balance + ?", amount),
		}).Error
	})
}

func (r *WalletRepository) Unfreeze(ctx context.Context, userID uint64, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Wallet{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
			"balance":        gorm.Expr("balance + ?", amount),
			"frozen_balance": gorm.Expr("frozen_balance - ?", amount),
		}).Error
	})
}

func (r *WalletRepository) DeductFrozen(ctx context.Context, userID uint64, amount float64) error {
	return r.db.Model(&model.Wallet{}).Where("user_id = ?", userID).
		Update("frozen_balance", gorm.Expr("frozen_balance - ?", amount)).Error
}

func (r *WalletRepository) AddBalance(ctx context.Context, userID uint64, amount float64) error {
	return r.db.Model(&model.Wallet{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"balance":      gorm.Expr("balance + ?", amount),
		"total_income": gorm.Expr("total_income + ?", amount),
	}).Error
}

func (r *WalletRepository) DeductBalance(ctx context.Context, userID uint64, amount float64) error {
	return r.db.Model(&model.Wallet{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"balance":       gorm.Expr("balance - ?", amount),
		"total_expense": gorm.Expr("total_expense + ?", amount),
	}).Error
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *TransactionRepository) ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Transaction{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error
	return transactions, total, err
}

type BookFavoriteRepository struct {
	db *gorm.DB
}

func NewBookFavoriteRepository(db *gorm.DB) *BookFavoriteRepository {
	return &BookFavoriteRepository{db: db}
}

func (r *BookFavoriteRepository) Create(ctx context.Context, favorite *model.BookFavorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *BookFavoriteRepository) Delete(ctx context.Context, userID, bookID uint64) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&model.BookFavorite{}).Error
}

func (r *BookFavoriteRepository) Exists(ctx context.Context, userID, bookID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BookFavorite{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count).Error
	return count > 0, err
}

func (r *BookFavoriteRepository) ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.BookFavorite, int64, error) {
	var favorites []model.BookFavorite
	var total int64

	db := r.db.WithContext(ctx).Model(&model.BookFavorite{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Preload("Book").Find(&favorites).Error
	return favorites, total, err
}

type VerificationCodeRepository struct {
	db *gorm.DB
}

func NewVerificationCodeRepository(db *gorm.DB) *VerificationCodeRepository {
	return &VerificationCodeRepository{db: db}
}

func (r *VerificationCodeRepository) Create(ctx context.Context, code *model.VerificationCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *VerificationCodeRepository) FindLatest(ctx context.Context, target string, codeType int8) (*model.VerificationCode, error) {
	var vc model.VerificationCode
	err := r.db.WithContext(ctx).Where("target = ? AND type = ? AND used = 0 AND expires_at > ?", target, codeType, time.Now()).
		Order("created_at DESC").First(&vc).Error
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

func (r *VerificationCodeRepository) MarkUsed(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.VerificationCode{}).Where("id = ?", id).Update("used", 1).Error
}
