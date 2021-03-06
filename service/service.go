package service

import (
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"github.com/phpunch/route-roam-api/model"
	"github.com/phpunch/route-roam-api/repository"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service interface {
	fileService
	postService
	authService
	RegisterUser(username string, password string) (int64, error)
	LoginUser(username string, password string) (int64, error)
}

type service struct {
	repository repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) RegisterUser(
	username string, password string,
) (int64, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	hashedPwd := hashAndSalt([]byte(password))
	user := &model.User{
		Username:    username,
		Password: hashedPwd,
	}

	return s.repository.AddUser(user)
}

func (s *service) LoginUser(
	username string, password string,
) (int64, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	user, err := s.repository.GetUser(username)
	if user == nil {
		return 0, fmt.Errorf("user not found")
	}
	if err != nil {
		return 0, err
	}
	if !comparePasswords(user.Password, password) {
		return 0, fmt.Errorf("Invalid password")
	}
	return int64(user.ID), nil
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Log.Debugf("%v", err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		log.Log.Debugf("%v", err)
		return false
	}
	return true
}
