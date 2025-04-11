package errcode

import (
	"github.com/gin-gonic/gin"
)

// Response 定义统一的响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误信息
	Data    interface{} `json:"data"`    // 数据
}

// NewResponse 创建一个新的响应
func NewResponse(code int, data interface{}) *Response {
	message, ok := ErrCodeMessages[code]
	if !ok {
		message = "未知错误"
	}

	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// WithMessage 设置自定义响应信息
func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

// JSON 输出JSON响应
func (r *Response) JSON(c *gin.Context) {
	httpStatus, ok := ErrCodeToHTTPStatus[r.Code]
	if !ok {
		httpStatus = 500 // 默认返回500
	}

	c.JSON(httpStatus, gin.H{
		"code":    r.Code,
		"message": r.Message,
		"data":    r.Data,
	})
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	r := NewResponse(0, data).WithMessage("成功")
	r.JSON(c)
}

// Fail 失败响应
func Fail(c *gin.Context, code int, data interface{}) {
	r := NewResponse(code, data)
	r.JSON(c)
}

// FailWithMessage 带自定义消息的失败响应
func FailWithMessage(c *gin.Context, code int, message string, data interface{}) {
	r := NewResponse(code, data).WithMessage(message)
	r.JSON(c)
}
