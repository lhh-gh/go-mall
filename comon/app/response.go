package app

import (
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/comon/logger"
)

// response 统一返回结构体
type response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`       //omitempty: 忽略空值
	Pagination *pagination `json:"pagination,omitempty"` //
}

// NewResponse 创建新的响应实例
// 参数说明：
//   - ctx: 上下文对象
//
// 返回值：
//   - *response: 返回新创建的响应实例
//
// 说明：
func NewResponse(ctx *gin.Context) *response {
	return &response{
		ctx: ctx,
	}
}

// SetPagination 设置Response的分页信息
func (r *response) SetPagination(pagination *pagination) *response {
	r.Pagination = pagination
	return r
}

func (r *response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		val, _ := r.ctx.Get("traceid")
		requestId = val.(string)
	}
	r.RequestId = requestId
	r.Data = data

	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

func (r *response) SuccessOk() {
	r.Success("")
}

func (r *response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		val, _ := r.ctx.Get("traceid")
		requestId = val.(string)
	}
	r.RequestId = requestId
	// 兜底记一条响应错误, 项目自定义的AppError中有错误链条, 方便出错后排查问题
	logger.New(r.ctx).Error("api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
