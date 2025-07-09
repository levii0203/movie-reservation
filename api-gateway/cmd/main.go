package main

import (
	"fmt"
	"log"
	"apigateway/pkg/middleware"

	"github.com/gin-gonic/gin"
)


const PORT string = "3000"



func main(){


	gin.ForceConsoleColor()


	gin.SetMode(gin.DebugMode)


	router := gin.New()



	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Version_one())



	v1 := router.Group("/v1")


	v1.GET("/")

	v1.GET("/users")

	v1.GET("/book")


	log.Fatal(router.Run(fmt.Sprintf(":"+PORT)))

}