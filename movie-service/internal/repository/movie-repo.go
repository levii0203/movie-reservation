package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/levii0203/movie-service/internal/db"
	"github.com/levii0203/movie-service/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoID = fmt.Errorf("no movie id provided")
	ErrMovieNotFound = fmt.Errorf("movie not found")
	ErrInvalidMovie = fmt.Errorf("invalid movie details")
	ErrMissingParameterValue = fmt.Errorf("parameter/value missing")
	ErrMovieAlreadyExists = fmt.Errorf("movie with similar details already exists")
	ErrClientDisconnected = fmt.Errorf("client disconnected")
)


type MovieRepository interface {
	// should be given movie's mongodb_id as the second parameter
	// if id is nil, returns ErrNoID
	GetMovieByID(ctx context.Context, id primitive.ObjectID) (*model.Movie, error)
	// inserting movie , counts if movie already exists before insert operation
	// checks if movie is a valid model
	InsertMovie(ctx context.Context, movie *model.Movie) (string,error)
	// updating a single parameter in a query
	// parameter & its value shouldn't be emoty , returns ErrMissingParameterValue
	UpdateMovieSingleParameter(ctx context.Context, id primitive.ObjectID, parameter string, value any ) (*model.Movie,error)
	// should be given movie's mongodb_id as the second parameter
	// if id is nil, returns ErrNoID
	// returns ErrMovieNotFound, if movie doesn;t exist
	 DeleteMovieByID(ctx context.Context, id primitive.ObjectID) error
	

	 FetchAllByCity(ctx context.Context, city string) ([]model.Movie, error)


	 FetchAllByCitySortedByViewCount(ctx context.Context, city string) ([]model.Movie, error)

	 
	 FetchAllByNameAndCity(ctx context.Context, name string, city string) ([]model.Movie,error)


	 UpdateFilledMovie(ctx context.Context, movieID primitive.ObjectID, seat string) error
	
}

type movieRepository struct {
	movieCollection *mongo.Collection
}


func NewMovieRepository() MovieRepository{
	return &movieRepository{
		movieCollection: db.Client.Database("test").Collection("movies"),
	}
}


func (repo *movieRepository) GetMovieByID(ctx context.Context, id primitive.ObjectID) (*model.Movie, error) {
	if id.IsZero() {
		return nil, ErrNoID
	}
	fmt.Println(id)

	var movie model.Movie
	movie_id := id.Hex()
	res := repo.movieCollection.FindOne(ctx,bson.M{"_id":movie_id}).Decode(&movie)
	if res != nil {
		if res == mongo.ErrNoDocuments {
			return nil,ErrMovieNotFound
		}

		return nil , res
	}

	return &movie,nil

}


func (repo *movieRepository) InsertMovie(ctx context.Context, movie *model.Movie) (string,error) {
	if movie.Title == "" || movie.ReleaseDate == "" {
		return "",ErrInvalidMovie
	}

	ctx_cnt,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	count, err := repo.movieCollection.CountDocuments(ctx_cnt,bson.M{"title":movie.Title,"city":movie.City,"cinema":movie.Cinema})
	if count>0 {
		return "", ErrMovieAlreadyExists
	}
	if err!=nil {
		return "", err
	}
	movie.ID = primitive.NilObjectID
	res,err:= repo.movieCollection.InsertOne(ctx, movie)
	if err!=nil {
		return "",err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(),nil
}

// for updating any movie parameter 
func (repo *movieRepository) UpdateMovieSingleParameter(ctx context.Context, id primitive.ObjectID, parameter string, value any ) (*model.Movie,error) {
	if value == "" || parameter == "" {
		return nil,ErrMissingParameterValue
	}

	res := repo.movieCollection.FindOneAndUpdate(ctx, bson.M{"_id":id},bson.M{parameter:value});
	if res.Err() != nil {
		if errors.Is(res.Err(),mongo.ErrNilDocument){
			return nil,ErrMovieNotFound
		} else if errors.Is(res.Err(),mongo.ErrClientDisconnected) {
			return nil,ErrClientDisconnected
		}
	}

	var movie model.Movie

	err := res.Decode(&movie)
	if err != nil {
		return nil , err
	}

	return &movie, err
}

func (repo *movieRepository) DeleteMovieByID(ctx context.Context, id primitive.ObjectID) error {
	if id.IsZero() {
		return ErrNoID
	}

	res , err := repo.movieCollection.DeleteOne(ctx, bson.M{"_id":id})
	if err != nil {
		if errors.Is(err,mongo.ErrNilDocument){
			return ErrMovieNotFound
		} else if errors.Is(err,mongo.ErrClientDisconnected) {
			return ErrClientDisconnected
		}
	}

	if res.DeletedCount == 0 {
		return ErrMovieNotFound
	}

	return nil
}

func (repo *movieRepository) FetchAllByCity(ctx context.Context, city string) ([]model.Movie, error) {
	if len(city)==0 {
		return nil, fmt.Errorf("empty parameters")
	}

	cur , err := repo.movieCollection.Find(ctx, bson.M{"city":city})
	if err!=nil{
		if errors.Is(err,mongo.ErrNilCursor){
			return nil, mongo.ErrNilCursor
		} else if errors.Is(err,mongo.ErrClientDisconnected) {
			return nil, ErrClientDisconnected
		}
	}
	defer cur.Close(context.TODO())
	
	var movies []model.Movie

	ctx_cur, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	cur.All(ctx_cur,&movies)

	return movies, nil
}

func (repo *movieRepository) FetchAllByCitySortedByViewCount(ctx context.Context, city string) ([]model.Movie, error) {
	if len(city)==0 {
		return nil, fmt.Errorf("empty parameters")
	}

	// sorting by view counts in descending , to get trending movies in a city
	sort_options := options.Find()
	sort_options.SetSort(bson.M{"view_count":-1})

	cur,err := repo.movieCollection.Find(ctx, bson.M{"city":city}, sort_options)
	if err!=nil {
		if errors.Is(err,mongo.ErrNilCursor){
			return nil, mongo.ErrNilCursor
		} else if errors.Is(err,mongo.ErrClientDisconnected) {
			return nil, ErrClientDisconnected
		}
	}
	defer cur.Close(context.TODO())

	var movies []model.Movie

	ctx_cur, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	cur.All(ctx_cur,&movies)

	return movies,nil
}

func (repo *movieRepository) FetchAllByNameAndCity(ctx context.Context, name string, city string) ([]model.Movie,error) {
	if len(city)==0 || len(name)==0 {
		return nil, fmt.Errorf("empty parameters")
	}

	curr, err := repo.movieCollection.Find(ctx,bson.M{"city":city, "name":name})
	if err!=nil {
		if errors.Is(err,mongo.ErrNilCursor){
			return nil, mongo.ErrNilCursor
		} else if errors.Is(err,mongo.ErrClientDisconnected) {
			return nil, ErrClientDisconnected
		}
	}
	defer curr.Close(context.TODO())

	var movies []model.Movie

	ctx_cur, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	curr.All(ctx_cur,&movies)

	return movies,nil
}

func (repo *movieRepository) UpdateFilledMovie(ctx context.Context, movieID primitive.ObjectID, seat string) error {
	if movieID.IsZero() || seat=="" {
		return ErrNoID
	}

	id := movieID.Hex()
	update := bson.M{"$addToSet":bson.M{"filled_seats":seat}}
	err := repo.movieCollection.FindOneAndUpdate(ctx,bson.M{"_id":id}, update )
	if err.Err()!=nil {
		if errors.Is(err.Err(),mongo.ErrNilDocument) {
			return ErrMovieNotFound
		} else if errors.Is(err.Err(),mongo.ErrClientDisconnected){
			return ErrClientDisconnected
		}
		return fmt.Errorf("internal server error")
	}

	return nil
}