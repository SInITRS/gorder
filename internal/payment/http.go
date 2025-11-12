package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
	// Define fields and methods for the payment handler
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	// Define your HTTP routes here
	logrus.Info("hi")

	c.POST("/api/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	// Handle the webhook logic here
}
