package service

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/pkg/utils/token"
)


var (

	ErrEmptyRequest = fmt.Errorf("no email/phone")

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


	token, err := token.Sign(res.ID.Hex(), res.Email)
	if err != nil{
		return "",err
	}

	return token,nil

}


