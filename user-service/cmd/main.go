package main

import (
	"fmt"
	"os"
	"user-service/internal/config"
	"user-service/pkg/middleware/cors"
	"user-service/pkg/middleware/route"
	"user-service/pkg/middleware/rate-limiter"
	"user-service/pkg/utils/redis"

	"github.com/gin-gonic/gin"
)

const ( 
	PORT string = "8000"
)


func main(){

	config.LoadEnv()
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.CORSMiddleware())

	rdb.Ping()
	defer rdb.Client.Close()
	fmt.Println("Redis connection successful")

	routes.UserRoutes(router)

	router.GET("/test",limiter.ApiRateLimiter(),func(c *gin.Context){
		user , ok := c.Get("user")
		if !ok {
			c.JSON(400,gin.H{"error":"no user found"})
			return
		}
		c.JSON(200,gin.H{"user":user})
	})

	router.Run(fmt.Sprintf(":"+os.Getenv("PORT")))
}