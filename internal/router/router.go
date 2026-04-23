package router

import (
	"bookshare/internal/config"
	"bookshare/internal/database"
	"bookshare/internal/handler"
	"bookshare/internal/middleware"
	"bookshare/internal/repository"
	"bookshare/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *database.RedisClient, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.Cors())

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	rentalOrderRepo := repository.NewRentalOrderRepository(db)
	sellOrderRepo := repository.NewSellOrderRepository(db)
	giftRecordRepo := repository.NewGiftRecordRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	favoriteRepo := repository.NewBookFavoriteRepository(db)

	userService := service.NewUserService(userRepo, walletRepo, &cfg.JWT, &cfg.Platform)
	bookService := service.NewBookService(bookRepo, userRepo, favoriteRepo)
	rentalService := service.NewRentalService(rentalOrderRepo, bookRepo, userRepo, walletRepo, txRepo, redis, &cfg.Platform)
	sellService := service.NewSellService(sellOrderRepo, bookRepo, userRepo, walletRepo, txRepo, redis, &cfg.Platform)
	giftService := service.NewGiftService(giftRecordRepo, bookRepo, userRepo, redis)

	userHandler := handler.NewUserHandler(userService)
	bookHandler := handler.NewBookHandler(bookService)
	rentalHandler := handler.NewRentalHandler(rentalService)
	sellHandler := handler.NewSellHandler(sellService)
	giftHandler := handler.NewGiftHandler(giftService)
	guideHandler := handler.NewGuideHandler()

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.POST("/logout", middleware.JWTAuth(&cfg.JWT), userHandler.Logout)
			users.GET("/me", middleware.JWTAuth(&cfg.JWT), userHandler.GetMe)
			users.PUT("/me", middleware.JWTAuth(&cfg.JWT), userHandler.UpdateProfile)
			users.GET("/:id", userHandler.GetUser)
		}

		books := api.Group("/books")
		{
			books.POST("", middleware.JWTAuth(&cfg.JWT), bookHandler.Create)
			books.GET("", bookHandler.List)
			books.GET("/hot", bookHandler.GetHotBooks)
			books.GET("/gift/latest", bookHandler.GetLatestGiftBooks)
			books.GET("/my", middleware.JWTAuth(&cfg.JWT), bookHandler.GetMyBooks)
			books.GET("/favorites", middleware.JWTAuth(&cfg.JWT), bookHandler.GetFavorites)
			books.GET("/isbn/:isbn", bookHandler.QueryISBN)
			books.GET("/:id", bookHandler.GetByID)
			books.PUT("/:id", middleware.JWTAuth(&cfg.JWT), bookHandler.Update)
			books.DELETE("/:id", middleware.JWTAuth(&cfg.JWT), bookHandler.Delete)
			books.POST("/:id/favorite", middleware.JWTAuth(&cfg.JWT), bookHandler.AddFavorite)
			books.DELETE("/:id/favorite", middleware.JWTAuth(&cfg.JWT), bookHandler.RemoveFavorite)
		}

		rentals := api.Group("/rentals")
		rentals.Use(middleware.JWTAuth(&cfg.JWT))
		{
			rentals.POST("", rentalHandler.CreateOrder)
			rentals.GET("/my", rentalHandler.GetMyOrders)
			rentals.GET("/:id", rentalHandler.GetOrder)
			rentals.POST("/:id/pay", rentalHandler.PayOrder)
			rentals.POST("/:id/cancel", rentalHandler.CancelOrder)
			rentals.POST("/:id/confirm-pickup", rentalHandler.ConfirmPickup)
			rentals.POST("/:id/return", rentalHandler.ReturnBook)
			rentals.POST("/:id/inspect", rentalHandler.Inspect)
			rentals.POST("/:id/rate", rentalHandler.Rate)
		}

		sells := api.Group("/sells")
		sells.Use(middleware.JWTAuth(&cfg.JWT))
		{
			sells.POST("", sellHandler.CreateOrder)
			sells.GET("/my", sellHandler.GetMyOrders)
			sells.GET("/:id", sellHandler.GetOrder)
			sells.POST("/:id/pay", sellHandler.PayOrder)
			sells.POST("/:id/cancel", sellHandler.CancelOrder)
			sells.POST("/:id/ship", sellHandler.Ship)
			sells.POST("/:id/pickup", sellHandler.ConfirmPickup)
			sells.POST("/:id/confirm", sellHandler.ConfirmReceive)
			sells.POST("/:id/rate", sellHandler.Rate)
		}

		gifts := api.Group("/gifts")
		gifts.Use(middleware.JWTAuth(&cfg.JWT))
		{
			gifts.POST("", giftHandler.CreateRecord)
			gifts.GET("/my", giftHandler.GetMyRecords)
			gifts.GET("/:id", giftHandler.GetRecord)
			gifts.POST("/:id/confirm", giftHandler.Confirm)
			gifts.POST("/:id/reject", giftHandler.Reject)
			gifts.POST("/:id/cancel", giftHandler.Cancel)
			gifts.POST("/:id/deliver", giftHandler.Deliver)
			gifts.POST("/:id/rate", giftHandler.Rate)
			gifts.GET("/book/:book_id/applications", giftHandler.GetApplications)
		}

		api.GET("/guide", guideHandler.GetGuide)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
