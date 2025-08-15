package main

import (
	"os"
	"fmt"
	"github.com/levii0203/movie-service/internal/config"
	"github.com/levii0203/movie-service/pkg/middleware/cors"
	"github.com/levii0203/movie-service/pkg/middleware/routes"
	"github.com/levii0203/movie-service/pkg/utils/rabbitmq"

	"github.com/gin-gonic/gin"
)


func main(){
	config.LoadEnv()

	rabbitmq.Init()
	defer rabbitmq.CLIENT.Close()

	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.CORSMiddleware())

	routes.MovieRoute(router)

	router.Run(fmt.Sprintf(":%s",os.Getenv("PORT")))
}