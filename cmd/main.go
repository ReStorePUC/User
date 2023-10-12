package main

import (
	pb "github.com/ReStorePUC/protobucket/user"
	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowFiles:       true,
	}))

	router.POST("/profile", uHandler.Register)
	router.POST("/login", uHandler.Login)

	router.POST("/file", fHandler.UploadFile)
	router.Static("/view-file/", "./uploads")
	router.DELETE("/file/:file", fHandler.DeleteFile)
	router.GET("/store/search/:name", uHandler.SearchStore)
	router.GET("/store/admin/search", uHandler.SearchAdminStore)

	router.POST("/private/store", uHandler.RegisterStore)
	router.GET("/store/:id", uHandler.GetStore)
	router.GET("/private/profile/:id", uHandler.GetProfile)
	router.PUT("/private/profile/:id", uHandler.UpdateProfile)

	router.GET("/private/self/store", uHandler.GetSelfStore)
	router.GET("/private/self/profile", uHandler.GetSelfProfile)

	router.Run(":8080")
}
