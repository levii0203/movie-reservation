package handler

import (
	"fmt"
	"github.com/levii0203/movie-service/internal/model"
	"github.com/levii0203/movie-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidRequestBody = fmt.Errorf("invalid request body")
	ErrInternalServerError = fmt.Errorf("internal server error")
	ErrNoIdGiven = fmt.Errorf("no id was given")
)

type MovieHandler struct {
	movie_service service.MovieService
}

func NewMovieHandler(s service.MovieService) *MovieHandler {
	return &MovieHandler{
		movie_service: s,
	}
}

func (h *MovieHandler) RegisterMovie() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Writer.Header().Set("Content-Type","application/json")
		if c.Request.Header.Get("Content-Type")!="application/json" {
			c.JSON(300, gin.H{"error": ErrInvalidRequestBody.Error()})
			return
		}

		var movie model.Movie

		if err:=c.ShouldBindJSON(&movie); err!=nil {
			c.JSON(300,gin.H{"error":err.Error()})
			return
		}

		validate := validator.New()

		err := validate.Struct(movie)
		if err != nil {
			c.JSON(300, gin.H{"error": ErrInvalidRequestBody.Error()})
			return
		}

		id,err := h.movie_service.PutMovie(&movie)
		if err!=nil {
			c.JSON(300,gin.H{"error":err.Error()})
			return
		}

		c.JSON(200,gin.H{"id":id,"error":nil})
	}
}

func (h *MovieHandler) GetMovieThruID() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Header().Set("Content-Type","application/json")
		id := c.Param("id")
		if id == "" {
			id = c.Query("id")
		}
		if id=="" {
			c.JSON(300,gin.H{"error":ErrNoIdGiven.Error()})
			return
		}

		doc,err := h.movie_service.GetMovie(id)
		if err!=nil {
			c.JSON(300,gin.H{"error":err.Error()})
			return
		}

		c.JSON(200,gin.H{"movie":doc,"error":nil})
	}
}

func (h *MovieHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Writer.Header().Set("Content-Type","application/json")
		id := c.Param("id")
		if id == "" {
			id = c.Query("id")
		}
		if id=="" {
			c.JSON(300,gin.H{"error":ErrNoIdGiven.Error()})
			return
		}

		res := h.movie_service.DeleteMovie(id)
		if res!=nil {
			c.JSON(300,gin.H{"error":res.Error()})
			return
		}

		c.JSON(200,gin.H{"ok":true})
	}
}