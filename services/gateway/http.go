package gateway

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/wafi04/chatting-app/config"
	"github.com/wafi04/chatting-app/config/database"
	authhandler "github.com/wafi04/chatting-app/services/auth/pkg/handler"
	authrepository "github.com/wafi04/chatting-app/services/auth/pkg/repository"
	authservice "github.com/wafi04/chatting-app/services/auth/pkg/service"
	"github.com/wafi04/chatting-app/services/comments"
	"github.com/wafi04/chatting-app/services/likes"
	posthandler "github.com/wafi04/chatting-app/services/post/handler"
	cloudrepo "github.com/wafi04/chatting-app/services/post/repository/cloud"
	postrepo "github.com/wafi04/chatting-app/services/post/repository/post"
	postservice "github.com/wafi04/chatting-app/services/post/service"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpRoutes(db *database.Database, mongoClient *mongo.Client, cld *cloudinary.Cloudinary) *gin.Engine {
	r := gin.Default()
	config.SetUpCors(r)
	middleware.ResponseTime(r)
	CheckCoon(r)

	// Cloudinary setup

	// Auth dependencies
	authRepo := authrepository.NewUserRepository(db.DB)
	authService := authservice.NewAuthService(authRepo)
	authHandler := authhandler.NewGateway(authService)

	commentRepo := comments.NewCommentRepository(db.DB, authRepo)

	commetService := comments.NewCommntService(commentRepo)
	commentHandler := comments.NewCommntHandler(commetService)
	// Post dependencies
	postRepo := postrepo.NewPostRepository(db.DB, commentRepo, authRepo)
	cloudRepo := cloudrepo.NewCloudinaryService(cld)
	postService := postservice.NewPostService(cloudRepo, postRepo)
	postHandler := posthandler.NewGateway(postService, authService)

	likerepo := likes.NewLikeRepository(mongoClient)
	likeHandler := likes.NewLikeHandler(likerepo)

	// Routes
	api := r.Group("/api/v1")
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware())

	auth := api.Group("/auth")
	authhandler.RegisterRoutes(auth, authHandler)
	post := authenticated.Group("/post")
	posthandler.RegisterRoutes(post, postHandler)

	comment := authenticated.Group("/comment")
	comments.RegisterRoutes(comment, commentHandler)
	like := authenticated.Group("/likes")
	likes.RegisterRoutes(like, likeHandler)
	return r
}
