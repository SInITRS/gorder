package main

import (
	"github.com/SInITRS/gorder/common/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	serviceName := viper.GetString("order.service-name")
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {

	})

}
