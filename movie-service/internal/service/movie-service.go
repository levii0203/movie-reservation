package service

import (
	"context"
	"movie-service/internal/model"
	"movie-service/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieService interface {
	PutMovie(movie *model.Movie) (string,error)

	GetMovie(id string) (*model.Movie, error)

	DeleteMovie(id string) (error)
}

type movieService struct {
	Movie_repo  repository.MovieRepository
}

func NewMovieService() MovieService {
	return &movieService{
		Movie_repo: repository.NewMovieRepository(),
	}
}

func (s *movieService) PutMovie(movie *model.Movie) (string,error) {
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	res, err := s.Movie_repo.InsertMovie(ctx, movie)
	if err!=nil {
		return "",err
	}
	return res,nil
}

func (s *movieService) GetMovie(id string) (*model.Movie, error){
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	movie_id , _ := primitive.ObjectIDFromHex(id)
	doc, err := s.Movie_repo.GetMovieByID(ctx,movie_id)
	if err!=nil {
		return nil, err
	}
	return doc, nil
}

func (s *movieService) DeleteMovie(id string) (error){
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	movie_id, _ := primitive.ObjectIDFromHex(id)
	err := s.Movie_repo.DeleteMovieByID(ctx,movie_id)
	if(err!=nil){
		return err
	}

	return nil
}
