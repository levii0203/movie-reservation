package main

import (
	"fmt"
	"log"
	"movie-service/internal/config"
	"movie-service/pkg/middleware/cors"
	"movie-service/pkg/middleware/routes"
	"movie-service/pkg/utils/rabbitmq"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)


func main(){
	config.LoadEnv()

	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.CORSMiddleware())

	routes.MovieRoute(router)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(){
		defer wg.Done()
		rabbitmq.KeepConsumingSeatRequests()

	}()

	router.Run(fmt.Sprintf(":%s",os.Getenv("PORT")))
	
	wg.Wait()
}