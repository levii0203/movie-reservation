package routes

import (
	"movie-service/internal/handler"
	"movie-service/internal/service"

	"github.com/gin-gonic/gin"
)

func MovieRoute(router *gin.Engine) {
	movie_service := service.NewMovieService()
	Handler := handler.NewMovieHandler(movie_service)

	router.GET("/:id",Handler.GetMovieThruID())
	router.POST("/add",Handler.RegisterMovie())
	router.POST("/:id",Handler.DeleteByID())
}