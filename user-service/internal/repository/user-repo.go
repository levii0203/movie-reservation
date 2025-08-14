package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user-service/internal/db"
	"user-service/internal/model"


	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var  ( 

	ErrEmptyRequest = fmt.Errorf("no email/phone")
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrDatabaseError = fmt.Errorf("database error")
	ErrInvalidPassword = fmt.Errorf("password is not valid")
	ErrBookingNotFound = fmt.Errorf("user has no this booking")
	ErrInvalidUserID = fmt.Errorf("invalid user ID")

)


type UserRepositoryInterface interface { 
	
	CreateUser(ctx context.Context, user *model.User) error

	FindUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)

	FindUserByEmail(ctx context.Context, email string) (*model.User , error)

	AddBookingByEmail(ctx context.Context, email string, booking_id primitive.ObjectID) error

	DeleteBookingByEmail(ctx context.Context, email string, booking_id primitive.ObjectID) error

	UpdatePhone(ctx context.Context, email string , phone string ) error 

	ChangePassword(ctx context.Context, email string, password string) error

}

type UserRepository struct {
	collection *mongo.Collection
}


func NewUserRepository() UserRepositoryInterface {
	return &UserRepository{
		collection: db.Client.Database("test").Collection("users"),
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *model.User)  error {
	if user.Email == "" || user.Phone == "" {
		return ErrEmptyRequest
	}

	validate := validator.New()
	if err:=validate.Struct(user); err!=nil {
		return err
	}

	ctx_cnt , cancel_cnt := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel_cnt()

	n , _ := repo.collection.CountDocuments(ctx_cnt, bson.M{"email":user.Email})
	if n>0 {
		return ErrUserAlreadyExists
	}

	_,err := repo.collection.InsertOne(ctx,user)
	if err!=nil {
		return err
	}


	return nil

}

func (repo *UserRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	if id.IsZero(){
		return nil, ErrInvalidUserID
	}

	var user model.User


	err := repo.collection.FindOne(ctx, bson.M{"_id":id}).Decode(&user)
	if err!=nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil , ErrUserNotFound
		}
		return nil , err
	}

	return &user, nil

}

func (repo *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := repo.collection.FindOne(ctx, bson.M{"email":email}).Decode(&user)
	if err!=nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return nil , ErrUserNotFound
		}
		return nil , err
	}

	return &user, nil
}

func (repo *UserRepository) AddBookingByEmail(ctx context.Context, email string, booking_id primitive.ObjectID) error {

	res, err:= repo.collection.UpdateOne(ctx, bson.M{"email":email}, bson.M{"$push":  bson.M{"bookings":booking_id},})
	if err!=nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return ErrUserNotFound
		} 
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return ErrDatabaseError
		}
		return err
	}

	if res.MatchedCount==0 {
		return ErrUserNotFound
	}
	
	return nil
}


func (repo *UserRepository) DeleteBookingByEmail(ctx context.Context, email string, booking_id primitive.ObjectID) error {
	if booking_id.IsZero() {
		return fmt.Errorf("booking ID cannot be zero")
	}

	res, err:= repo.collection.UpdateOne(ctx, bson.M{"email":email}, bson.M{"$pull":  bson.M{"bookings":booking_id},})
	if err!=nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return ErrUserNotFound
		} 
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return ErrDatabaseError
		}
		return err
	}

	if res.MatchedCount==0 {
		return ErrUserNotFound
	}
	
	if res.ModifiedCount == 0 {
		return ErrBookingNotFound
	}

	return nil
}






func (repo *UserRepository) UpdatePhone(ctx context.Context, email string , phone string ) error {

	res, err := repo.collection.UpdateOne(ctx, bson.M{"email":email}, bson.M{"$set": bson.M{"phone":phone},})
	if err!=nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return ErrUserNotFound
		} 
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return ErrDatabaseError
		}
		return err
	}

	if res.MatchedCount==0 {
		return ErrUserNotFound
	}
	
	return nil
}



func (repo *UserRepository) ChangePassword(ctx context.Context, email string, password string) error {

	if len(password)<6 {
		return ErrInvalidPassword
	}

	res,err := repo.collection.UpdateOne(ctx, bson.M{"email":email}, bson.M{"$set": bson.M{"password":password},})
		if err!=nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return ErrUserNotFound
		} 
		if errors.Is(err, mongo.ErrClientDisconnected) {
			return ErrDatabaseError
		}
		return err
	}

	if res.MatchedCount==0 {
		return ErrUserNotFound
	}
	
	return nil

}


