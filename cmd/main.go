package main

import (
	pb "github.com/ReStorePUC/protobucket/generated"
	"github.com/gin-gonic/gin"
	"github.com/restore/user/config"
	"github.com/restore/user/controller"
	"github.com/restore/user/handler"
	"github.com/restore/user/repository"
	"github.com/restore/user/service"
	"google.golang.org/grpc"
	"log"
	"net"
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
	fHandler := handler.NewFile()

	// GRPC
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterUserServer(s, handler.NewUserServer(uController))
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// HTTP
	router := gin.Default()
	router.POST("/profile", uHandler.Register)
	router.POST("/login", uHandler.Login)

	router.POST("/file", fHandler.UploadFile)
	router.GET("/file/:file", fHandler.GetFile)
	router.DELETE("/file/:file", fHandler.DeleteFile)

	router.POST("/private/store", uHandler.RegisterStore)
	router.GET("/private/store/:id", uHandler.GetStore)
	router.GET("/private/store/search/:name", uHandler.SearchStore)
	router.GET("/private/profile/:id", uHandler.GetProfile)
	router.PUT("/private/profile/:id", uHandler.UpdateProfile)

	router.Run(":8080")
}
