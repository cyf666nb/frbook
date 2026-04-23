package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Account      string         `gorm:"uniqueIndex;size:50;not null" json:"account"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	BirthDate    time.Time      `gorm:"type:date;not null" json:"birth_date"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	School       string         `gorm:"size:100" json:"school"`
	Campus       string         `gorm:"size:100" json:"campus"`
	Grade        string         `gorm:"size:20" json:"grade"`
	InterestTags []string       `gorm:"type:json;serializer:json" json:"interest_tags"`
	CreditScore  int            `gorm:"default:100" json:"credit_score"`
	Status       int8           `gorm:"default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsUnderAge(maxAge int) bool {
	age := time.Now().Year() - u.BirthDate.Year()
	if time.Now().YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age < maxAge
}

type Book struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64         `gorm:"index;not null" json:"user_id"`
	ISBN           string         `gorm:"size:20;index" json:"isbn"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	Author         string         `gorm:"size:255" json:"author"`
	Publisher      string         `gorm:"size:100" json:"publisher"`
	CoverImage     string         `gorm:"size:255" json:"cover_image"`
	Description    string         `gorm:"type:text" json:"description"`
	Category       string         `gorm:"size:50" json:"category"`
	Mode           int8           `gorm:"not null" json:"mode"`
	Status         int8           `gorm:"default:1;index" json:"status"`
	DailyRent      *float64       `gorm:"type:decimal(10,2)" json:"daily_rent,omitempty"`
	WeeklyRent     *float64       `gorm:"type:decimal(10,2)" json:"weekly_rent,omitempty"`
	Deposit        *float64       `gorm:"type:decimal(10,2)" json:"deposit,omitempty"`
	MinRentDays    int            `gorm:"default:1" json:"min_rent_days,omitempty"`
	SellPrice      *float64       `gorm:"type:decimal(10,2)" json:"sell_price,omitempty"`
	Images         []string       `gorm:"type:json;serializer:json;not null" json:"images"`
	PickupLocation string         `gorm:"size:255" json:"pickup_location"`
	ViewCount      int            `gorm:"default:0" json:"view_count"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Book) TableName() string {
	return "books"
}

const (
	BookModeRental = 1
	BookModeSell   = 2
	BookModeGift   = 3
)

const (
	BookStatusOffline   = 0
	BookStatusOnline    = 1
	BookStatusTraded    = 2
)

type RentalOrder struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo          string         `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	BookID           uint64         `gorm:"index;not null" json:"book_id"`
	BookTitle        string         `gorm:"size:255;not null" json:"book_title"`
	OwnerID          uint64         `gorm:"index;not null" json:"owner_id"`
	RenterID         uint64         `gorm:"index;not null" json:"renter_id"`
	DailyRent        float64        `gorm:"type:decimal(10,2);not null" json:"daily_rent"`
	Deposit          float64        `gorm:"type:decimal(10,2);not null" json:"deposit"`
	RentDays         int            `gorm:"not null" json:"rent_days"`
	TotalRent        float64        `gorm:"type:decimal(10,2);not null" json:"total_rent"`
	TotalAmount      float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	StartTime        *time.Time     `json:"start_time,omitempty"`
	EndTime          *time.Time     `json:"end_time,omitempty"`
	ActualReturnTime *time.Time     `json:"actual_return_time,omitempty"`
	Status           int8           `gorm:"default:0;index" json:"status"`
	OverdueFee       float64        `gorm:"type:decimal(10,2);default:0" json:"overdue_fee"`
	RefundAmount     *float64       `gorm:"type:decimal(10,2)" json:"refund_amount,omitempty"`
	SettledAmount    *float64       `gorm:"type:decimal(10,2)" json:"settled_amount,omitempty"`
	OwnerRating      *int8          `json:"owner_rating,omitempty"`
	RenterRating     *int8          `json:"renter_rating,omitempty"`
	OwnerComment     string         `gorm:"size:500" json:"owner_comment,omitempty"`
	RenterComment    string         `gorm:"size:500" json:"renter_comment,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Book             *Book          `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Owner            *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Renter           *User          `gorm:"foreignKey:RenterID" json:"renter,omitempty"`
}

func (RentalOrder) TableName() string {
	return "rental_orders"
}

const (
	RentalStatusPendingPayment   = 0
	RentalStatusPaidWaitingPick  = 1
	RentalStatusRenting          = 2
	RentalStatusWaitingInspection = 3
	RentalStatusCompleted        = 4
	RentalStatusCancelled        = 5
	RentalStatusOverdue          = 6
)

type SellOrder struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo        string         `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	BookID         uint64         `gorm:"index;not null" json:"book_id"`
	BookTitle      string         `gorm:"size:255;not null" json:"book_title"`
	SellerID       uint64         `gorm:"index;not null" json:"seller_id"`
	BuyerID        uint64         `gorm:"index;not null" json:"buyer_id"`
	Price          float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	PlatformFee    float64        `gorm:"type:decimal(10,2);default:0" json:"platform_fee"`
	SettleAmount   *float64       `gorm:"type:decimal(10,2)" json:"settle_amount,omitempty"`
	Status         int8           `gorm:"default:0;index" json:"status"`
	DeliveryType   *int8          `json:"delivery_type,omitempty"`
	DeliveryNo     string         `gorm:"size:50" json:"delivery_no,omitempty"`
	DeliveryCompany string        `gorm:"size:50" json:"delivery_company,omitempty"`
	DeliveryTime   *time.Time     `json:"delivery_time,omitempty"`
	ReceiveTime    *time.Time     `json:"receive_time,omitempty"`
	AutoSettleTime *time.Time     `gorm:"index" json:"auto_settle_time,omitempty"`
	SellerRating   *int8          `json:"seller_rating,omitempty"`
	BuyerRating    *int8          `json:"buyer_rating,omitempty"`
	SellerComment  string         `gorm:"size:500" json:"seller_comment,omitempty"`
	BuyerComment   string         `gorm:"size:500" json:"buyer_comment,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Book           *Book          `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Seller         *User          `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Buyer          *User          `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
}

func (SellOrder) TableName() string {
	return "sell_orders"
}

const (
	SellStatusPendingPayment  = 0
	SellStatusPaidWaitingShip = 1
	SellStatusShipped         = 2
	SellStatusCompleted       = 3
	SellStatusCancelled       = 4
)

const (
	DeliveryTypeExpress = 1
	DeliveryTypePickup  = 2
)

type GiftRecord struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	RecordNo      string         `gorm:"uniqueIndex;size:32;not null" json:"record_no"`
	BookID        uint64         `gorm:"index;not null" json:"book_id"`
	BookTitle     string         `gorm:"size:255;not null" json:"book_title"`
	GiverID       uint64         `gorm:"index;not null" json:"giver_id"`
	ReceiverID    uint64         `gorm:"index;not null" json:"receiver_id"`
	Status        int8           `gorm:"default:0;index" json:"status"`
	DeliveryType  *int8          `json:"delivery_type,omitempty"`
	DeliveryTime  *time.Time     `json:"delivery_time,omitempty"`
	GiverRating   *int8          `json:"giver_rating,omitempty"`
	ReceiverRating *int8         `json:"receiver_rating,omitempty"`
	GiverComment  string         `gorm:"size:500" json:"giver_comment,omitempty"`
	ReceiverComment string       `gorm:"size:500" json:"receiver_comment,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Book          *Book          `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Giver         *User          `gorm:"foreignKey:GiverID" json:"giver,omitempty"`
	Receiver      *User          `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
}

func (GiftRecord) TableName() string {
	return "gift_records"
}

const (
	GiftStatusPending        = 0
	GiftStatusConfirmed      = 1
	GiftStatusCompleted      = 2
	GiftStatusCancelled      = 3
)

type Transaction struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionNo string    `gorm:"uniqueIndex;size:32;not null" json:"transaction_no"`
	UserID        uint64    `gorm:"index;not null" json:"user_id"`
	RelatedUserID *uint64   `json:"related_user_id,omitempty"`
	Type          int8      `gorm:"not null;index" json:"type"`
	OrderType     int8      `gorm:"not null" json:"order_type"`
	OrderID       uint64    `gorm:"index;not null" json:"order_id"`
	OrderNo       string    `gorm:"size:32;not null" json:"order_no"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	BalanceBefore *float64  `gorm:"type:decimal(10,2)" json:"balance_before,omitempty"`
	BalanceAfter  *float64  `gorm:"type:decimal(10,2)" json:"balance_after,omitempty"`
	Status        int8      `gorm:"default:1" json:"status"`
	Remark        string    `gorm:"size:255" json:"remark,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

const (
	TransactionTypeRentalDeposit     = 1
	TransactionTypeRentalRent        = 2
	TransactionTypeRentalDepositRefund = 3
	TransactionTypeSellPayment       = 4
	TransactionTypeSellSettlement    = 5
	TransactionTypePlatformFee       = 6
)

const (
	TransactionTypeOrderRental = 1
	TransactionTypeOrderSell   = 2
)

type Wallet struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance       float64   `gorm:"type:decimal(10,2);default:0" json:"balance"`
	FrozenBalance float64   `gorm:"type:decimal(10,2);default:0" json:"frozen_balance"`
	TotalIncome   float64   `gorm:"type:decimal(10,2);default:0" json:"total_income"`
	TotalExpense  float64   `gorm:"type:decimal(10,2);default:0" json:"total_expense"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallets"
}

type CreditRecord struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint64    `gorm:"index;not null" json:"user_id"`
	ScoreChange     int       `gorm:"not null" json:"score_change"`
	ScoreBefore     int       `gorm:"not null" json:"score_before"`
	ScoreAfter      int       `gorm:"not null" json:"score_after"`
	Reason          string    `gorm:"size:255;not null" json:"reason"`
	RelatedOrderType *int8    `json:"related_order_type,omitempty"`
	RelatedOrderID  *uint64   `json:"related_order_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

func (CreditRecord) TableName() string {
	return "credit_records"
}

type BookFavorite struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:uk_user_book;not null" json:"user_id"`
	BookID    uint64    `gorm:"uniqueIndex:uk_user_book;index;not null" json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (BookFavorite) TableName() string {
	return "book_favorites"
}

type VerificationCode struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Target    string    `gorm:"index:idx_target_type;size:50;not null" json:"target"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	Type      int8      `gorm:"index:idx_target_type;not null" json:"type"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	Used      int8      `gorm:"default:0" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}

const (
	VerifyTypeRegister     = 1
	VerifyTypeLogin        = 2
	VerifyTypeResetPassword = 3
)
