package response

import (
	"github.com/LucienVen/golang-auth-service/internal/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// PageResponse 分页响应结构
type PageResponse struct {
	Total int64       `json:"total"` // 总数
	Items interface{} `json:"items"` // 数据列表
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.ErrCodeSuccess,
		Message: errors.GetMessage(errors.ErrCodeSuccess),
		Data:    data,
	})
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, total int64, items interface{}) {
	Success(c, PageResponse{
		Total: total,
		Items: items,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	if message == "" {
		message = errors.GetMessage(code)
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = errors.GetMessage(errors.ErrCodeParamInvalid)
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:    errors.ErrCodeParamInvalid,
		Message: message,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    errors.ErrCodeUnauthorized,
		Message: errors.GetMessage(errors.ErrCodeUnauthorized),
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    errors.ErrCodeForbidden,
		Message: errors.GetMessage(errors.ErrCodeForbidden),
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code:    errors.ErrCodeNotFound,
		Message: errors.GetMessage(errors.ErrCodeNotFound),
	})
}

// ServerError 服务器错误
func ServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    errors.ErrCodeSystemError,
		Message: err.Error(),
	})
}
