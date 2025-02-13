package likes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/pkg/response"
)

type LikeHandler struct {
	likerepo Repository
}

func NewLikeHandler(likerepo Repository) *LikeHandler {
	return &LikeHandler{
		likerepo: likerepo,
	}
}

// HandleChangeLikeComment handles toggling likes on comments
func (lh *LikeHandler) HandleChangeLikeComment(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Comment ID is required")
		return
	}

	data, err := lh.likerepo.ChangeLikeComment(c, user.UserId, commentID)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to change like", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Like status changed successfully", data)
}

// HandleChangeLikePost handles toggling likes on posts
func (lh *LikeHandler) HandleChangeLikePost(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID := c.Param("id")
	if postID == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Post ID is required")
		return
	}

	err = lh.likerepo.ChangeLikePost(c, user.UserId, postID)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to change like", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Like status changed successfully", nil)
}

// HandleGetUserCommentLikes gets all comments liked by a user
func (lh *LikeHandler) HandleGetUserCommentLikes(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	likes, err := lh.likerepo.GetUserCommentLikes(c, user.UserId)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get user comment likes", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "User comment likes retrieved successfully", likes)
}

// HandleGetUserPostLikes gets all posts liked by a user
func (lh *LikeHandler) HandleGetUserPostLikes(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	likes, err := lh.likerepo.GetUserPostLikes(c, user.UserId)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get user post likes", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "User post likes retrieved successfully", likes)
}

// HandleGetCommentLikesCount gets the total number of likes for a comment
func (lh *LikeHandler) HandleGetCommentLikesCount(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Comment ID is required")
		return
	}

	count, err := lh.likerepo.GetCommentLikesCount(c, commentID)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get comment likes count", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Comment likes count retrieved successfully", gin.H{"count": count})
}

// HandleGetPostLikesCount gets the total number of likes for a post
func (lh *LikeHandler) HandleGetPostLikesCount(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Post ID is required")
		return
	}

	count, err := lh.likerepo.GetPostLikesCount(c, postID)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get post likes count", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Post likes count retrieved successfully", gin.H{"count": count})
}
func (lh *LikeHandler) HandleGetPostLikedUser(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	postID := c.Param("id")
	if postID == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Post ID is required")
		return
	}

	count, err := lh.likerepo.GetUserLiked(c, "post_id", postID, user.UserId)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get post likes count", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Post likes count retrieved successfully", count)
}
func (lh *LikeHandler) HandleGetCommentLikedUser(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	commentID := c.Param("id")
	if commentID == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Post ID is required")
		return
	}

	count, err := lh.likerepo.GetUserLiked(c, "comment_id", commentID, user.UserId)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusInternalServerError, "Failed to get post likes count", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Post likes count retrieved successfully", count)
}
