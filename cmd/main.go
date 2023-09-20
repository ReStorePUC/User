package main

import (
	"github.com/gin-gonic/gin"
	"github.com/restore/user/config"
	"github.com/restore/user/controller"
	"github.com/restore/user/handler"
	"github.com/restore/user/repository"
	"github.com/restore/user/service"
)

func main() {
	config.Init()
	kgCfg := config.NewKongConfig()
	dbCfg := config.NewDBConfig()

	db, err := repository.Init(dbCfg)
	if err != nil {
		panic(err)
	}

	uRepo := repository.NewUser(db)
	kong := service.NewKong(kgCfg)
	uController := controller.NewUser(uRepo, kong)
	uHandler := handler.NewUser(uController)

	router := gin.Default()
	router.POST("/profile", uHandler.Register)
	router.POST("/login", uHandler.Login)

	router.POST("/private/store", uHandler.RegisterStore)
	router.GET("/private/store/:id", uHandler.GetStore)
	router.GET("/private/store/search/:name", uHandler.SearchStore)
	router.GET("/private/profile/:id", uHandler.GetProfile)
	router.PUT("/private/profile/:id", uHandler.UpdateProfile)

	router.Run(":8080")
}
