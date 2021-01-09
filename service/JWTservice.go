package service

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type authService interface {
	CreateToken(userid int64) (string, error)
}

func (s *service) CreateToken(userid int64) (string, error) {
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userid
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
