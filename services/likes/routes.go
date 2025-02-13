package likes

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *LikeHandler) {
	r.POST("/comment/:id", handler.HandleChangeLikeComment)
	r.POST("/post/:id", handler.HandleChangeLikePost)
	r.GET("/post/:id/isliked", handler.HandleGetPostLikedUser)
	r.GET("/comment/:id/isliked", handler.HandleGetCommentLikedUser)
	r.GET("/user/comments", handler.HandleGetUserCommentLikes)
	r.GET("/user/posts", handler.HandleGetUserPostLikes)
	r.GET("/comment/:id/count", handler.HandleGetCommentLikesCount)
	r.GET("/post/:id/count", handler.HandleGetPostLikesCount)
}
