package common

import (
	"encoding/json"
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
	r := response{
		Errno:   0,
		Msg:     "success",
		Data:    data,
		TraceID: tracing.TraceID(c.Request.Context()),
	}
	resp, _ := json.Marshal(r)
	c.Set("response", string(resp))
	c.JSON(http.StatusOK, r)
}
func (base *BaseResponse) error(c *gin.Context, err error) {
	r := response{
		Errno:   2,
		Msg:     err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	}
	resp, _ := json.Marshal(r)
	c.Set("response", string(resp))
	c.JSON(http.StatusOK, r)
}
