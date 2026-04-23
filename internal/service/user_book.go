package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bookshare/internal/config"
	"bookshare/internal/model"
	"bookshare/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repository.UserRepository
	walletRepo *repository.WalletRepository
	jwtConfig *config.JWTConfig
	platform  *config.PlatformConfig
}

func NewUserService(userRepo *repository.UserRepository, walletRepo *repository.WalletRepository, jwtConfig *config.JWTConfig, platform *config.PlatformConfig) *UserService {
	return &UserService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		jwtConfig:  jwtConfig,
		platform:   platform,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	exists, err := s.userRepo.ExistsByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户已存在")
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("出生日期格式错误")
	}

	user := &model.User{
		Account:   req.Account,
		BirthDate: birthDate,
	}

	if !user.IsUnderAge(s.platform.MaxAge) {
		return nil, errors.New("年龄不符合要求,需小于30岁")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	wallet := &model.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}
	_ = s.walletRepo.Create(ctx, wallet)

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	token, err := generateToken(s.jwtConfig, user.ID)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:    token,
		ExpireAt: time.Now().Add(s.jwtConfig.ExpireTime).Unix(),
		User:     user,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint64, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.School != "" {
		user.School = req.School
	}
	if req.Campus != "" {
		user.Campus = req.Campus
	}
	if req.Grade != "" {
		user.Grade = req.Grade
	}
	if req.InterestTags != nil {
		user.InterestTags = req.InterestTags
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateCreditScore(ctx context.Context, userID uint64, scoreChange int, reason string, orderType *int8, orderID *uint64) error {
	return s.userRepo.UpdateCreditScore(ctx, userID, scoreChange, reason, orderType, orderID)
}

func generateToken(cfg *config.JWTConfig, userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(cfg.ExpireTime).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     cfg.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

type BookService struct {
	bookRepo     *repository.BookRepository
	userRepo     *repository.UserRepository
	favoriteRepo *repository.BookFavoriteRepository
}

func NewBookService(bookRepo *repository.BookRepository, userRepo *repository.UserRepository, favoriteRepo *repository.BookFavoriteRepository) *BookService {
	return &BookService{
		bookRepo:     bookRepo,
		userRepo:     userRepo,
		favoriteRepo: favoriteRepo,
	}
}

func (s *BookService) Create(ctx context.Context, userID uint64, req *model.CreateBookRequest) (*model.Book, error) {
	book := &model.Book{
		UserID:         userID,
		ISBN:           req.ISBN,
		Title:          req.Title,
		Author:         req.Author,
		Publisher:      req.Publisher,
		CoverImage:     req.CoverImage,
		Description:    req.Description,
		Category:       req.Category,
		Mode:           req.Mode,
		Status:         model.BookStatusOnline,
		DailyRent:      req.DailyRent,
		WeeklyRent:     req.WeeklyRent,
		Deposit:        req.Deposit,
		MinRentDays:    req.MinRentDays,
		SellPrice:      req.SellPrice,
		Images:         req.Images,
		PickupLocation: req.PickupLocation,
	}

	if book.Mode == model.BookModeGift {
		zero := float64(0)
		book.SellPrice = &zero
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetByID(ctx context.Context, id uint64, userID uint64) (*model.BookDetail, error) {
	book, err := s.bookRepo.FindByIDWithUser(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.bookRepo.IncrementViewCount(ctx, id)

	isFavorited := false
	if userID > 0 {
		isFavorited, _ = s.favoriteRepo.Exists(ctx, userID, id)
	}

	ownerName := ""
	ownerSchool := ""
	if book.User != nil {
		ownerName = book.User.Nickname
		if ownerName == "" {
			ownerName = book.User.Account
		}
		ownerSchool = book.User.School
	}

	return &model.BookDetail{
		Book:        book,
		IsFavorited: isFavorited,
		OwnerName:   ownerName,
		OwnerSchool: ownerSchool,
	}, nil
}

func (s *BookService) Update(ctx context.Context, bookID, userID uint64, req *model.UpdateBookRequest) (*model.Book, error) {
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if book.UserID != userID {
		return nil, errors.New("无权限修改")
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	if req.CoverImage != "" {
		book.CoverImage = req.CoverImage
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Category != "" {
		book.Category = req.Category
	}
	if req.DailyRent != nil {
		book.DailyRent = req.DailyRent
	}
	if req.WeeklyRent != nil {
		book.WeeklyRent = req.WeeklyRent
	}
	if req.Deposit != nil {
		book.Deposit = req.Deposit
	}
	if req.MinRentDays > 0 {
		book.MinRentDays = req.MinRentDays
	}
	if req.SellPrice != nil {
		book.SellPrice = req.SellPrice
	}
	if req.Images != nil {
		book.Images = req.Images
	}
	if req.PickupLocation != "" {
		book.PickupLocation = req.PickupLocation
	}

	if err := s.bookRepo.Update(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) Delete(ctx context.Context, bookID, userID uint64) error {
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return err
	}

	if book.UserID != userID {
		return errors.New("无权限删除")
	}

	return s.bookRepo.Delete(ctx, bookID)
}

func (s *BookService) List(ctx context.Context, query *model.BookListQuery) ([]model.Book, int64, error) {
	return s.bookRepo.List(ctx, query)
}

func (s *BookService) ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]model.Book, int64, error) {
	return s.bookRepo.ListByUserID(ctx, userID, page, pageSize)
}

func (s *BookService) GetHotBooks(ctx context.Context, limit int) ([]model.Book, error) {
	return s.bookRepo.GetHotBooks(ctx, limit)
}

func (s *BookService) GetLatestGiftBooks(ctx context.Context, limit int) ([]model.Book, error) {
	return s.bookRepo.GetLatestGiftBooks(ctx, limit)
}

func (s *BookService) AddFavorite(ctx context.Context, userID, bookID uint64) error {
	exists, _ := s.favoriteRepo.Exists(ctx, userID, bookID)
	if exists {
		return nil
	}

	favorite := &model.BookFavorite{
		UserID: userID,
		BookID: bookID,
	}
	return s.favoriteRepo.Create(ctx, favorite)
}

func (s *BookService) RemoveFavorite(ctx context.Context, userID, bookID uint64) error {
	return s.favoriteRepo.Delete(ctx, userID, bookID)
}

func (s *BookService) ListFavorites(ctx context.Context, userID uint64, page, pageSize int) ([]model.BookFavorite, int64, error) {
	return s.favoriteRepo.ListByUserID(ctx, userID, page, pageSize)
}

func generateOrderNo(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, uuid.New().String()[:20])
}
