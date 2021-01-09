package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phpunch/route-roam-api/controller"
	"github.com/phpunch/route-roam-api/infrastructure/db"
	"github.com/phpunch/route-roam-api/log"
	"github.com/phpunch/route-roam-api/middleware"
	"github.com/phpunch/route-roam-api/repository"
	"github.com/phpunch/route-roam-api/service"
	"github.com/spf13/viper"
	"os"
)

func getLogger() log.Logger {
	logLevel := viper.GetString("Log.Level")
	logLevel = log.NormalizeLogLevel(logLevel)

	logColor := viper.GetBool("Log.Color")
	logJSON := viper.GetBool("Log.JSON")

	logger, err := log.NewLogger(&log.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logLevel,
		ConsoleJSONFormat: logJSON,
		Color:             logColor,
	}, log.InstanceZapLogger)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return logger
}

func main() {
	////////////// GIN //////
	logger := getLogger()
	dbHandler, err := db.NewDB(logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := repository.NewRepository(dbHandler)
	s := service.NewService(r)
	c := controller.NewController(s)
	mw := middleware.New(s)

	router := gin.Default()

	router.Use(middleware.GinMiddleware("http://localhost:3000"))

	router.POST("/register", c.RegisterUser)
	router.POST("/login", c.LoginUser)

	router.Use(mw.AuthorizeToken())
	{
		router.POST("/logout", c.LogoutUser)
		router.POST("/file", c.UploadFiles)
		router.GET("/file/*filepath", c.GetFile)
		router.POST("/post", c.CreatePost)
		router.POST("/like", c.LikePost)
		router.POST("/unlike", c.UnlikePost)
		router.GET("/posts", c.GetPosts)
		router.POST("/comment", c.CommentPost)
	}

	router.Run()
}
