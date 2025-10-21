package app

import (
	"github.com/SInITRS/gorder/order/app/command"
	"github.com/SInITRS/gorder/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
