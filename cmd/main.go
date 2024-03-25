package main

import (
	"log"

	"github.com/geedotrar/mygram/internal/handler"
	"github.com/geedotrar/mygram/internal/infrastructure"
	"github.com/geedotrar/mygram/internal/repository"
	"github.com/geedotrar/mygram/internal/router"
	"github.com/geedotrar/mygram/internal/service"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {

	server()
}

func server() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	g := gin.Default()
	g.Use(gin.Recovery())

	usersGroup := g.Group("/users")

	gorm := infrastructure.NewGormPostgres()
	userRepo := repository.NewUserQuery(gorm)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)
	userRouter := router.NewUserRouter(usersGroup, userHdl)
	userRouter.Mount()
	photosGroup := g.Group("/photos")
	photoRepo := repository.NewPhotoQuery(gorm)
	photoSvc := service.NewPhotoService(photoRepo, userRepo)
	photoHdl := handler.NewPhotoHandler(photoSvc)
	photoRouter := router.NewPhotoRouter(photosGroup, photoHdl)
	photoRouter.Mount()
	commentsGroup := g.Group("/comments")
	commentRepo := repository.NewCommentQuery(gorm)
	commentSvc := service.NewCommentService(commentRepo, userRepo, photoRepo)
	commentHdl := handler.NewCommentHandler(commentSvc)
	commentRouter := router.NewCommentRouter(commentsGroup, commentHdl)
	commentRouter.Mount()
	socialMediasGroup := g.Group("/socialmedias")
	socialMediaRepo := repository.NewSocialMediaQuery(gorm)
	socialMediaSvc := service.NewSocialMediaService(socialMediaRepo, userRepo)
	socialMediaHdl := handler.NewSocialMediaHandler(socialMediaSvc)
	socialMediaRouter := router.NewSocialMediaRouter(socialMediasGroup, socialMediaHdl)
	socialMediaRouter.Mount()

	g.Run(":3000")
}
