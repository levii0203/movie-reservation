package main

import (
	"os"
	"fmt"
	"booking-service/internal/config"
	"booking-service/pkg/middleware/routes"
	"booking-service/pkg/utils/rabbitmq"
	rdb "booking-service/pkg/utils/redis"

	"github.com/gin-gonic/gin"
)

func main(){
	config.LoadEnv()

	defer rabbitmq.Booking.Close()
	defer rabbitmq.AmqpClient.Close()

	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	routes.LobbyRoute(router)

	rdb.Ping()

	router.Run(fmt.Sprintf(":%s",os.Getenv("PORT")))

}