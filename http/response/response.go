package response

import (
	"chat/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 数据结构体
type Response struct {
	// Code 业务状态码
	Code int `json:"code"`

	// Message 提示信息
	Message    string      `json:"message"`

	// Data 数据，用interface{}的目的是可以用任意数据
	Data       interface{} `json:"data"`

	// Meta 源数据,存储如请求ID,分页等信息
	Meta       Meta        `json:"meta"`

	// Errors 错误提示，如 xx字段不能为空等
	Errors     []ErrorItem `json:"errors"`
}

// Meta 元数据
type Meta struct {
	RequestId      string                 `json:"request_id"`
	// 还可以集成分页信息等
}

// ErrorItem 错误项
type ErrorItem struct {
	Key   string `json:"key"`
	Value string `json:"error"`
}

// New return response instance
func New() *Response {
	return &Response{
		Code: 200,
		Message:    "",
		Data:       nil,
		Meta: Meta{
			RequestId: utils.UUID(),
		},
		Errors: []ErrorItem{},
	}
}

// Wrapper include context
type Wrapper struct {
	ctx *gin.Context
}

// WrapContext
func WrapContext(ctx *gin.Context) *Wrapper {
	return &Wrapper{ctx:ctx}
}

// Json 输出json,支持自定义response结构体
func (wrapper *Wrapper) Json(response *Response) {
	wrapper.ctx.JSON(http.StatusOK, response)
}

// Success 成功的输出
func (wrapper *Wrapper) Success( data interface{}) {
	response := New()
	response.Data = data
	wrapper.Json(response)
}

// Error 错误输出
func (wrapper *Wrapper) Error( code int, message string) {
	response := New()
	response.Code = code
	response.Message = message
	wrapper.Json(response)
}