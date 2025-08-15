package main

import (
	"os"
	"fmt"
	"github.com/levii0203/booking-service/internal/config"
	"github.com/levii0203/booking-service/pkg/middleware/routes"
	"github.com/levii0203/booking-service/pkg/utils/rabbitmq"
	rdb "github.com/levii0203/booking-service/pkg/utils/redis"

	"github.com/gin-gonic/gin"
)

func main(){
	config.LoadEnv()

	// initializing rabbitmq client
	rabbitmq.Init()
	defer rabbitmq.CLIENT.Close()

	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	routes.LobbyRoute(router)

	//Pinging the redis
	rdb.Ping()

	router.Run(fmt.Sprintf(":%s",os.Getenv("PORT")))

}