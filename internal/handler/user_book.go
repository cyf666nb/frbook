package handler

import (
	"strconv"

	"bookshare/internal/middleware"
	"bookshare/internal/model"
	"bookshare/internal/service"
	"bookshare/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, gin.H{"id": user.ID})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodePasswordError, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *UserHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeUserNotFound)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeUserNotFound)
		return
	}

	publicUser := gin.H{
		"id":           user.ID,
		"nickname":     user.Nickname,
		"avatar":       user.Avatar,
		"school":       user.School,
		"campus":       user.Campus,
		"credit_score": user.CreditScore,
	}

	response.Success(c, publicUser)
}

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	book, err := h.bookService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, book)
}

func (h *BookHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	userID := middleware.GetUserID(c)

	book, err := h.bookService.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, response.CodeBookNotFound)
		return
	}

	response.Success(c, book)
}

func (h *BookHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	book, err := h.bookService.Update(c.Request.Context(), id, userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, book)
}

func (h *BookHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.bookService.Delete(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *BookHandler) List(c *gin.Context) {
	var query model.BookListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err)
		return
	}

	books, total, err := h.bookService.List(c.Request.Context(), &query)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.PageSuccess(c, books, total, query.Page, query.PageSize)
}

func (h *BookHandler) GetMyBooks(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	books, total, err := h.bookService.ListByUserID(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.PageSuccess(c, books, total, page, pageSize)
}

func (h *BookHandler) GetHotBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	books, err := h.bookService.GetHotBooks(c.Request.Context(), limit)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, books)
}

func (h *BookHandler) GetLatestGiftBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	books, err := h.bookService.GetLatestGiftBooks(c.Request.Context(), limit)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, books)
}

func (h *BookHandler) AddFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.bookService.AddFavorite(c.Request.Context(), userID, id); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *BookHandler) RemoveFavorite(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.bookService.RemoveFavorite(c.Request.Context(), userID, id); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *BookHandler) GetFavorites(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	favorites, total, err := h.bookService.ListFavorites(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.PageSuccess(c, favorites, total, page, pageSize)
}

func (h *BookHandler) QueryISBN(c *gin.Context) {
	isbn := c.Param("isbn")

	bookInfo := &model.ISBNBookInfo{
		ISBN:      isbn,
		Title:     "示例图书",
		Author:    "示例作者",
		Publisher: "示例出版社",
		Cover:     "https://example.com/cover.jpg",
		Summary:   "这是一本示例图书",
	}

	response.Success(c, bookInfo)
}
