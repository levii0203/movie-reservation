package routes

import (
	"github.com/gin-gonic/gin"
)

func LobbyRoute(router *gin.Engine){
	router.POST("/seatlock")
}