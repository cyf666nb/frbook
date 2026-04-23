package handler

import (
	"strconv"

	"bookshare/internal/middleware"
	"bookshare/internal/model"
	"bookshare/internal/service"
	"bookshare/pkg/response"

	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	rentalService *service.RentalService
}

func NewRentalHandler(rentalService *service.RentalService) *RentalHandler {
	return &RentalHandler{rentalService: rentalService}
}

func (h *RentalHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.CreateRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	order, err := h.rentalService.CreateOrder(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, order)
}

func (h *RentalHandler) PayOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.PayRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.rentalService.PayOrder(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RentalHandler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.rentalService.CancelOrder(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RentalHandler) ConfirmPickup(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.rentalService.ConfirmPickup(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RentalHandler) ReturnBook(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.rentalService.ReturnBook(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RentalHandler) Inspect(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.InspectRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.rentalService.Inspect(c.Request.Context(), id, userID, &req); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *RentalHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	order, err := h.rentalService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeOrderNotFound)
		return
	}

	response.Success(c, order)
}

func (h *RentalHandler) GetMyOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var query model.RentalListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err)
		return
	}

	orders, total, err := h.rentalService.ListByUserID(c.Request.Context(), userID, &query)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.PageSuccess(c, orders, total, query.Page, query.PageSize)
}

func (h *RentalHandler) Rate(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.RateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.rentalService.Rate(c.Request.Context(), id, userID, &req); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

type SellHandler struct {
	sellService *service.SellService
}

func NewSellHandler(sellService *service.SellService) *SellHandler {
	return &SellHandler{sellService: sellService}
}

func (h *SellHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.CreateSellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	order, err := h.sellService.CreateOrder(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, order)
}

func (h *SellHandler) PayOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.PaySellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.sellService.PayOrder(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *SellHandler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.sellService.CancelOrder(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *SellHandler) Ship(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.ShipSellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.sellService.Ship(c.Request.Context(), id, userID, &req); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *SellHandler) ConfirmPickup(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.sellService.ConfirmPickup(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *SellHandler) ConfirmReceive(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.sellService.ConfirmReceive(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *SellHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	order, err := h.sellService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeOrderNotFound)
		return
	}

	response.Success(c, order)
}

func (h *SellHandler) GetMyOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var query model.SellListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err)
		return
	}

	orders, total, err := h.sellService.ListByUserID(c.Request.Context(), userID, &query)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.PageSuccess(c, orders, total, query.Page, query.PageSize)
}

func (h *SellHandler) Rate(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.RateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.sellService.Rate(c.Request.Context(), id, userID, &req); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

type GiftHandler struct {
	giftService *service.GiftService
}

func NewGiftHandler(giftService *service.GiftService) *GiftHandler {
	return &GiftHandler{giftService: giftService}
}

func (h *GiftHandler) CreateRecord(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.CreateGiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	record, err := h.giftService.CreateRecord(c.Request.Context(), userID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, record)
}

func (h *GiftHandler) Confirm(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.giftService.Confirm(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GiftHandler) Reject(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.giftService.Reject(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GiftHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.giftService.Cancel(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GiftHandler) Deliver(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	if err := h.giftService.Deliver(c.Request.Context(), id, userID); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *GiftHandler) GetRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	record, err := h.giftService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeOrderNotFound)
		return
	}

	response.Success(c, record)
}

func (h *GiftHandler) GetMyRecords(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var query model.GiftListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err)
		return
	}

	records, total, err := h.giftService.ListByUserID(c.Request.Context(), userID, &query)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.PageSuccess(c, records, total, query.Page, query.PageSize)
}

func (h *GiftHandler) GetApplications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	records, err := h.giftService.GetApplicationsByBookID(c.Request.Context(), bookID, userID)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, records)
}

func (h *GiftHandler) Rate(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError)
		return
	}

	var req model.RateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err)
		return
	}

	if err := h.giftService.Rate(c.Request.Context(), id, userID, &req); err != nil {
		response.ErrorWithMessage(c, response.CodeError, err.Error())
		return
	}

	response.Success(c, nil)
}
