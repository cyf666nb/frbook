package model

type RegisterRequest struct {
	Account   string `json:"account" binding:"required,min=2,max=50"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
	BirthDate string `json:"birth_date" binding:"required"`
	Code      string `json:"code" binding:"required,len=6"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
	User     *User  `json:"user"`
}

type UpdateUserRequest struct {
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	School       string   `json:"school"`
	Campus       string   `json:"campus"`
	Grade        string   `json:"grade"`
	InterestTags []string `json:"interest_tags"`
}

type SendCodeRequest struct {
	Target string `json:"target" binding:"required"`
	Type   int8   `json:"type" binding:"required,oneof=1 2 3"`
}

type CreateBookRequest struct {
	ISBN           string   `json:"isbn"`
	Title          string   `json:"title" binding:"required"`
	Author         string   `json:"author"`
	Publisher      string   `json:"publisher"`
	CoverImage     string   `json:"cover_image"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Mode           int8     `json:"mode" binding:"required,oneof=1 2 3"`
	DailyRent      *float64 `json:"daily_rent"`
	WeeklyRent     *float64 `json:"weekly_rent"`
	Deposit        *float64 `json:"deposit"`
	MinRentDays    int      `json:"min_rent_days"`
	SellPrice      *float64 `json:"sell_price"`
	Images         []string `json:"images" binding:"required,min=1"`
	PickupLocation string   `json:"pickup_location"`
}

type UpdateBookRequest struct {
	Title          string   `json:"title"`
	Author         string   `json:"author"`
	Publisher      string   `json:"publisher"`
	CoverImage     string   `json:"cover_image"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	DailyRent      *float64 `json:"daily_rent"`
	WeeklyRent     *float64 `json:"weekly_rent"`
	Deposit        *float64 `json:"deposit"`
	MinRentDays    int      `json:"min_rent_days"`
	SellPrice      *float64 `json:"sell_price"`
	Images         []string `json:"images"`
	PickupLocation string   `json:"pickup_location"`
}

type BookListQuery struct {
	Mode     int8   `form:"mode"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
	School   string `form:"school"`
	Sort     string `form:"sort" binding:"omitempty,oneof=created_at view_count price"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type CreateRentalRequest struct {
	BookID   uint64 `json:"book_id" binding:"required"`
	RentDays int    `json:"rent_days" binding:"required,min=1"`
}

type PayRentalRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type InspectRentalRequest struct {
	Passed bool   `json:"passed"`
	Remark string `json:"remark"`
}

type RentalListQuery struct {
	Role     int8   `form:"role" binding:"omitempty,oneof=1 2"`
	Status   *int8  `form:"status"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type CreateSellRequest struct {
	BookID uint64 `json:"book_id" binding:"required"`
}

type PaySellRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type ShipSellRequest struct {
	DeliveryType    int8   `json:"delivery_type" binding:"required,oneof=1 2"`
	DeliveryCompany string `json:"delivery_company"`
	DeliveryNo      string `json:"delivery_no"`
}

type RefundSellRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type SellListQuery struct {
	Role     int8  `form:"role" binding:"omitempty,oneof=1 2"`
	Status   *int8 `form:"status"`
	Page     int   `form:"page" binding:"omitempty,min=1"`
	PageSize int   `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type CreateGiftRequest struct {
	BookID  uint64 `json:"book_id" binding:"required"`
	Message string `json:"message"`
}

type GiftListQuery struct {
	Role     int8  `form:"role" binding:"omitempty,oneof=1 2"`
	Status   *int8 `form:"status"`
	Page     int   `form:"page" binding:"omitempty,min=1"`
	PageSize int   `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type RateRequest struct {
	Rating  int8   `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type ISBNBookInfo struct {
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Cover     string `json:"cover"`
	Summary   string `json:"summary"`
}

type BookDetail struct {
	*Book
	IsFavorited bool   `json:"is_favorited"`
	OwnerName   string `json:"owner_name"`
	OwnerSchool string `json:"owner_school"`
}

type RentalOrderDetail struct {
	*RentalOrder
	Book       *Book `json:"book"`
	Owner      *User `json:"owner"`
	Renter     *User `json:"renter"`
	IsOverdue  bool  `json:"is_overdue"`
	OverdueDays int  `json:"overdue_days"`
}

type SellOrderDetail struct {
	*SellOrder
	Book   *Book `json:"book"`
	Seller *User `json:"seller"`
	Buyer  *User `json:"buyer"`
}

type GiftRecordDetail struct {
	*GiftRecord
	Book     *Book `json:"book"`
	Giver    *User `json:"giver"`
	Receiver *User `json:"receiver"`
}

type UserStats struct {
	TotalPublish     int `json:"total_publish"`
	TotalRental      int `json:"total_rental"`
	TotalSell        int `json:"total_sell"`
	TotalGift        int `json:"total_gift"`
	TotalTransaction int `json:"total_transaction"`
}

type PlatformStats struct {
	TotalUsers      int64   `json:"total_users"`
	TotalBooks      int64   `json:"total_books"`
	TotalRental     int64   `json:"total_rental"`
	TotalSell       int64   `json:"total_sell"`
	TotalGift       int64   `json:"total_gift"`
	TotalTransaction float64 `json:"total_transaction"`
}

type SchoolRank struct {
	School     string `json:"school"`
	BookCount  int64  `json:"book_count"`
	TradeCount int64  `json:"trade_count"`
}

type BookRank struct {
	BookID    uint64 `json:"book_id"`
	Title     string `json:"title"`
	ViewCount int    `json:"view_count"`
}
