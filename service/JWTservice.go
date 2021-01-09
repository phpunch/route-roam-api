package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/phpunch/route-roam-api/middleware"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"strconv"
	"time"
)

type authService interface {
	CreateToken(userid int64) (*TokenDetails, error)
	CreateAuth(userid int64, td *TokenDetails) error
	ExtractTokenMetadata(r *http.Request) (*AccessDetails, error)
	FetchAuth(authD *AccessDetails) (int64, error)
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
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
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
	td.AccessToken, err = at.SignedString([]byte("secret_here"))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
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
	token, err := middleware.VerifyToken(r)
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
