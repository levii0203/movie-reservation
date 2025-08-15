package routes

import (
	"github.com/levii0203/booking-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func LobbyRoute(router *gin.Engine){
	router.POST("/seatlock", handler.SeatLockHandler())
}