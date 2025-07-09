package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var (

	ErrFailedSigned = fmt.Errorf("failed to sign token")
)



var (

	jwt_secret = "password"
)


type Claims struct {

	user_id string
	email string
	jwt.RegisteredClaims

}


func Sign(user_id,email string) (string,error){
	claims := Claims {
		user_id: user_id,
		email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "user_service",
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	token_string, err := token.SignedString(jwt_secret)
	if err!=nil {
		return err.Error(),ErrFailedSigned
	}

	return token_string, nil
}


