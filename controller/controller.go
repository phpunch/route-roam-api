package controller

import (
	"github.com/phpunch/route-roam-api/service"
)

type Controller interface {
	authController
}

type controller struct {
	service service.Service
}

func NewController(s service.Service) Controller {
	return &controller{
		service: s,
	}
}
