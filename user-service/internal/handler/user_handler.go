package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"user-service/internal/service"
	"user-service/internal/model"


	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/go-playground/validator/v10"
)


var (

	ErrNoId = fmt.Errorf("no user id provided")
	ErrInvalidUserID = fmt.Errorf("invalid user id provided")
	ErrInvalidUser = fmt.Errorf("invalid user data provided")

)


type UserHandler interface {
	 GetUserByID() gin.HandlerFunc
	 SignUp() gin.HandlerFunc
	 Login() gin.HandlerFunc
}


type userHandler struct {
	user_service *service.UserService
}


func NewUserHandler(user_service *service.UserService) UserHandler {
	return &userHandler{
		user_service: user_service,	
	}
}


func (h *userHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("id")
		if user_id == "" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrNoId.Error(),"user":nil})
			return
		}

		ctx,cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		ID, err:= primitive.ObjectIDFromHex(user_id)
		if err !=nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserID.Error(),"user":nil})
			return
		}


		user , err := h.user_service.User_repo.FindUserByID(ctx, ID)
		if err != nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(),"user":nil})
		}

		c.JSON(http.StatusOK, gin.H{"error":nil, "user":user})
	}
}




func (h *userHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context){
		var user model.User

		if err:= c.ShouldBindJSON(&user); err!=nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := h.user_service.RegisterUser(&user)
		if err!=nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error":nil})
	}
}


func (h *userHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")

		type user_details struct {
			Email		string	`bson:"email" json:"email" validate:"required,email"`
			Password   	string  `bson:"password" json:"-" validate:"required,min=8,alphanum"`
		}

		var user user_details
		if err := c.ShouldBindJSON(&user); err != nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {	
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUser.Error()})
			return	
		}


		token, err:= h.user_service.LoginUser(user.Email, user.Password);
		if err!=nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": nil, "token": token})

	
	}
}