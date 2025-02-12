package gateway

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/wafi04/chatting-app/config"
	"github.com/wafi04/chatting-app/config/database"
	"github.com/wafi04/chatting-app/config/env"
	authhandler "github.com/wafi04/chatting-app/services/auth/pkg/handler"
	authrepository "github.com/wafi04/chatting-app/services/auth/pkg/repository"
	authservice "github.com/wafi04/chatting-app/services/auth/pkg/service"
	posthandler "github.com/wafi04/chatting-app/services/post/pkg/handler"
	cloudrepo "github.com/wafi04/chatting-app/services/post/pkg/repository/cloud"
	postrepo "github.com/wafi04/chatting-app/services/post/pkg/repository/post"
	postservice "github.com/wafi04/chatting-app/services/post/pkg/service"
	"github.com/wafi04/chatting-app/services/shared/middleware"
)

func SetUpRoutes(db *database.Database) *gin.Engine {
	r := gin.Default()
	config.SetUpCors(r)
	middleware.ResponseTime(r)
	CheckCoon(r)

	cld, err := cloudinary.NewFromParams(
		env.LoadEnv("CLOUDINARY_CLOUD_NAME"),
		env.LoadEnv("CLOUDINARY_API_KEY"),
		env.LoadEnv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		panic(err)

	}

	authrepo := authrepository.NewUserRepository(db.DB)
	authservice := authservice.NewAuthService(authrepo)
	authsrv := authhandler.NewGateway(authservice)

	postrepo := postrepo.NewPostRepository(db.DB)
	cloudrepo := cloudrepo.NewCloudinaryService(cld)
	postservice := postservice.NewPostService(cloudrepo, postrepo)
	postsrv := posthandler.NewGateway(postservice, authservice)

	api := r.Group("/api/v1")

	authenticated := r.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	auth := api.Group("/auth")
	authhandler.RegisterRoutes(auth, authsrv)

	post := api.Group("/post")
	posthandler.RegisterRoutes(post, postsrv)

	return r
}
