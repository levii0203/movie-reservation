package token

import (
	"fmt"
	"time"
	"user-service/internal/model"

	"github.com/golang-jwt/jwt/v5"
)


var (
	ErrFailedSigned = fmt.Errorf("failed to sign token")
	ErrInvalidToken = fmt.Errorf("invalid token")
	ErrInvalidClaims = fmt.Errorf("invalid claims in token")
)

const (

	jwt_secret string = "password"
	expiry_time time.Duration = 30 * time.Second
)

type Claims struct {

	User model.User

	jwt.RegisteredClaims
}

func Sign(user model.User) (string,error){
	claims := Claims {
		User : user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "user_service",
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry_time)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	token_string, err := token.SignedString([]byte(jwt_secret))
	if err!=nil {
		return err.Error(),ErrFailedSigned
	}

	return token_string, nil
}

func Validate(token_string string) (*Claims,error) {
	claims := &Claims{}
    token, err := jwt.ParseWithClaims(token_string, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwt_secret), nil
	})
    if err != nil {
        return nil, err
    }

	if _,ok := token.Claims.(*Claims); !ok {
		return nil , ErrInvalidClaims
	}

    if !token.Valid {
        return nil, ErrInvalidToken
    }
    return claims, nil
}

