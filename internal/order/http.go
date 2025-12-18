package main

import (
	"fmt"

	"github.com/SInITRS/gorder/common"
	client "github.com/SInITRS/gorder/common/client/order"
	"github.com/SInITRS/gorder/order/app"
	"github.com/SInITRS/gorder/order/app/command"
	"github.com/SInITRS/gorder/order/app/query"
	"github.com/SInITRS/gorder/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		err  error
		resp struct {
			CustomerID  string `json:"customer_id"`
			OrderID     string `json:"order_id"`
			RedirectURL string `json:"redirect_url"`
		}
	)

	defer func() {
		H.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp.CustomerID = req.CustomerId
	resp.OrderID = r.OrderID
	resp.RedirectURL = fmt.Sprintf("http://192.168.2.32:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID)

}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp struct {
			Order *client.Order `json:"order"`
		}
	)
	defer func() {
		H.Response(c, err, &resp)
	}()
	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		return
	}
	resp.Order = convertor.NewOrderConvertor().EntityToClient(o)
}
