package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/twinj/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authService interface {
	CreateToken(userid int64) (*TokenDetails, error)
	CreateAuth(userid int64, td *TokenDetails) error
	ExtractTokenMetadata(r *http.Request) (*AccessDetails, error)
	FetchAuth(authD *AccessDetails) (int64, error)
	DeleteAuth(uuid string) (int64, error)
	VerifyToken(tokenString string, tokenType string) (*jwt.Token, error)
	RefreshToken(token *jwt.Token) (map[string]string, error)
}

// Inspired by this blog
// https://codeburst.io/using-jwt-for-authentication-in-a-golang-application-e0357d579ce2

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}
type AccessDetails struct {
	AccessUUID string
	UserID     int64
}

func (s *service) CreateToken(userid int64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 60).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(viper.GetString("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(viper.GetString("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (s *service) CreateAuth(userid int64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := s.repository.SaveToken(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now))
	if errAccess != nil {
		return errAccess
	}
	errRefresh := s.repository.SaveToken(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (s *service) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	tokenString := ExtractToken(r)
	token, err := s.VerifyToken(tokenString, "ACCESS_SECRET")
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

func (s *service) FetchAuth(authD *AccessDetails) (int64, error) {
	userID, err := s.repository.FetchToken(authD.AccessUUID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *service) DeleteAuth(uuid string) (int64, error) {
	deleted, err := s.repository.DeleteToken(uuid)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (s *service) VerifyToken(tokenString string, tokenType string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := viper.GetString(tokenType)
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token error: %v", err)
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("failed to convert *jwt.Claims type")
	}
	return token, nil
}

func (s *service) RefreshToken(token *jwt.Token) (map[string]string, error) {
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return nil, fmt.Errorf("unauthorized")
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unauthorized: %v", err)
		}
		//Delete the previous Refresh Token
		deleted, delErr := s.DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return nil, fmt.Errorf("unauthorized: %v", delErr)
		}
		//Create new pairs of refresh and access tokens
		td, createErr := s.CreateToken(userID)
		if createErr != nil {
			return nil, fmt.Errorf("unauthorized: %v", createErr)
		}
		//save the tokens metadata to redis
		saveErr := s.CreateAuth(userID, td)
		if saveErr != nil {
			return nil, fmt.Errorf("unauthorized: %v", saveErr)
		}
		tokens := map[string]string{
			"access_token":  td.AccessToken,
			"refresh_token": td.RefreshToken,
		}
		return tokens, nil
	} else {
		return nil, fmt.Errorf("refresh expired")

	}
}
