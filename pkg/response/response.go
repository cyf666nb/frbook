package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

const (
	CodeSuccess = 0
	CodeError   = 1
)

const (
	CodeParamError      = 1001
	CodeUnauthorized    = 1002
	CodeForbidden       = 1003
	CodeNotFound        = 1004
	CodeInternalError   = 1005
	CodeTooManyRequests = 1006
)

const (
	CodeUserExists       = 2001
	CodeUserNotFound     = 2002
	CodePasswordError    = 2003
	CodeVerifyCodeError  = 2004
	CodeAgeNotQualified  = 2005
	CodeTokenInvalid     = 2006
	CodeTokenExpired     = 2007
)

const (
	CodeBookNotFound   = 3001
	CodeBookOffline    = 3002
	CodeBookLocked     = 3003
	CodeBookNotOwner   = 3004
	CodeBookAlreadyTraded = 3005
)

const (
	CodeOrderNotFound     = 4001
	CodeOrderStatusError  = 4002
	CodeBalanceNotEnough  = 4003
	CodeLockFailed        = 4004
	CodeOrderNotOwner     = 4005
	CodeOrderTimeout      = 4006
)

const (
	CodePaymentFailed = 5001
	CodeRefundFailed  = 5002
)

var codeMessages = map[int]string{
	CodeSuccess:           "成功",
	CodeError:             "失败",
	CodeParamError:        "参数错误",
	CodeUnauthorized:      "未授权",
	CodeForbidden:         "禁止访问",
	CodeNotFound:          "资源不存在",
	CodeInternalError:     "服务器内部错误",
	CodeTooManyRequests:   "请求过于频繁",
	CodeUserExists:        "用户已存在",
	CodeUserNotFound:      "用户不存在",
	CodePasswordError:     "密码错误",
	CodeVerifyCodeError:   "验证码错误",
	CodeAgeNotQualified:   "年龄不符合要求",
	CodeTokenInvalid:      "Token无效",
	CodeTokenExpired:      "Token已过期",
	CodeBookNotFound:      "图书不存在",
	CodeBookOffline:       "图书已下架",
	CodeBookLocked:        "图书已被锁定",
	CodeBookNotOwner:      "非图书所有者",
	CodeBookAlreadyTraded: "图书已交易",
	CodeOrderNotFound:     "订单不存在",
	CodeOrderStatusError:  "订单状态错误",
	CodeBalanceNotEnough:  "余额不足",
	CodeLockFailed:        "锁定失败",
	CodeOrderNotOwner:     "非订单所有者",
	CodeOrderTimeout:      "订单已超时",
	CodePaymentFailed:     "支付失败",
	CodeRefundFailed:      "退款失败",
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: codeMessages[CodeSuccess],
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int) {
	message, ok := codeMessages[code]
	if !ok {
		message = "未知错误"
	}
	
	httpStatus := http.StatusBadRequest
	switch code {
	case CodeUnauthorized:
		httpStatus = http.StatusUnauthorized
	case CodeForbidden:
		httpStatus = http.StatusForbidden
	case CodeNotFound:
		httpStatus = http.StatusNotFound
	case CodeTooManyRequests:
		httpStatus = http.StatusTooManyRequests
	case CodeInternalError:
		httpStatus = http.StatusInternalServerError
	}
	
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func ErrorWithMessage(c *gin.Context, code int, message string) {
	httpStatus := http.StatusBadRequest
	if code == CodeUnauthorized {
		httpStatus = http.StatusUnauthorized
	} else if code == CodeForbidden {
		httpStatus = http.StatusForbidden
	} else if code == CodeNotFound {
		httpStatus = http.StatusNotFound
	} else if code == CodeInternalError {
		httpStatus = http.StatusInternalServerError
	}
	
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func ParamError(c *gin.Context, err error) {
	message := "参数错误"
	if err != nil {
		message = err.Error()
	}
	ErrorWithMessage(c, CodeParamError, message)
}

func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
