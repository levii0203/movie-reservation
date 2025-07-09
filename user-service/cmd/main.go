package main

import (
	"fmt"
	"user-service/internal/config"
	"user-service/pkg/middleware/route"

	"github.com/gin-gonic/gin"
)


const PORT string = "8000"


func main(){

	config.LoadEnv()
	router := gin.New()


	router.Use(gin.Logger())
	router.Use(gin.Recovery())



	routes.UserRoutes(router)

	

	router.Run(fmt.Sprintf(":"+PORT))
}