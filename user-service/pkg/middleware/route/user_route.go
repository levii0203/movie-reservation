package routes

import (

	"github.com/levii0203/user-service/internal/handler"
	"github.com/levii0203/user-service/internal/service"
	
	"github.com/gin-gonic/gin"

)


func UserRoutes(r *gin.Engine){

	User_Service :=  service.NewUserService()
	Handler := handler.NewUserHandler(User_Service)


	r.GET("/:id",Handler.GetUserByID())
	r.POST("/signup",Handler.SignUp())
	r.POST("/login",Handler.Login())


}