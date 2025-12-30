package common

import (
	"net/http"

	"github.com/SInITRS/gorder/common/tracing"
	"github.com/gin-gonic/gin"
)

type BaseResponse struct{}

type response struct {
	Errno   int    `json:"errno"`
	Msg     string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(c *gin.Context, err error, data any) {
	if err != nil {
		base.error(c, err)
	} else {
		base.success(c, data)
	}
}

func (base *BaseResponse) success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{
		Errno:   0,
		Msg:     "success",
		Data:    data,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}
func (base *BaseResponse) error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, response{
		Errno:   2,
		Msg:     err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}
