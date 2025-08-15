package handler

import (
	"github.com/levii0203/booking-service/internal/model"
	"github.com/levii0203/booking-service/pkg/utils/rabbitmq"
	"github.com/levii0203/booking-service/pkg/utils/redis"
	"context"
	"encoding/json"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request body")
)
// Seat lock handler
//publishes seat lock alert to movie service
func SeatLockHandler() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Writer.Header().Set("Content-Type", "application/json")
		var seatLock model.SeatLock
		var seatredis model.SeatRedis

		if err := c.ShouldBind(&seatLock); err != nil {
			c.JSON(400, gin.H{"error":ErrInvalidRequest})
			return
		}
		if seatLock.MovieID=="" || seatLock.UserID=="" {
			c.JSON(300,gin.H{"error":ErrInvalidRequest})
			return
		}
		seatredis.Seat = seatLock.Seat
		seatredis.UserID = seatLock.UserID

		ctx_r,cancel_r := context.WithTimeout(context.Background(),5*time.Second)
		defer cancel_r()
		ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		items, err:= rdb.Client.LRange(ctx_r, seatLock.MovieID, 0, -1).Result()
		if err!=nil {
			c.JSON(500,gin.H{"error":"seat lock failed"})
			return
		}
		for _,item := range items {
			var s model.SeatRedis
			json.Unmarshal([]byte(item),&s)
			if s.Seat==seatLock.Seat {
				c.JSON(300,gin.H{"error":"seat is already locked"})
				return
			}
		}
	
		t,err := json.Marshal(seatredis)
		if err!=nil {
			c.JSON(500,gin.H{"error":"seat lock failed"})
			return
		}
		
		rdb.Client.LPush(ctx,seatLock.MovieID,string(t),30*time.Minute)

		err = rabbitmq.CLIENT.AlertSeatLocked(seatLock)
		if err!=nil {
			c.JSON(500, gin.H{"error":"cannot create a book request"})
			return
		}
	}
} 

func SeatUnlockHandler() gin.HandlerFunc {
	return func(c *gin.Context){


	}
}