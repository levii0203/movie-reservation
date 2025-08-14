package service

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/pkg/utils/token"
	"user-service/pkg/utils/redis"
)


var (

	ErrEmptyRequest = fmt.Errorf("no email/phone")
	ErrInternalServer = fmt.Errorf("internal server error")

)


type UserService struct {
	User_repo repository.UserRepositoryInterface
}


func NewUserService() *UserService {
	return &UserService{
		User_repo: repository.NewUserRepository(),
	}
}

func (s *UserService) RegisterUser(user *model.User) error {
	if user.Email == "" && user.Phone == "" {
		return ErrEmptyRequest
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.User_repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(email,password string) (string,error) {
	if email == "" && password == "" {
		return "",ErrEmptyRequest
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := s.User_repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "",err
	}


	token, err := token.Sign(*res)
	if err != nil{
		return "",err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rdb.Client.Set(ctx,token,1,30*time.Second).Err(); err!=nil {
		return "", ErrInternalServer
	}

	return token,nil

}

func (s *UserService) ChangePassword(email, password string) error {
	return nil
}
