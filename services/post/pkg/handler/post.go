package posthandler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	authservice "github.com/wafi04/chatting-app/services/auth/pkg/service"
	postservice "github.com/wafi04/chatting-app/services/post/pkg/service"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/pkg/response"
	"github.com/wafi04/chatting-app/services/shared/types"
)

type PostHandler struct {
	postclient *postservice.PostService
	auhClient  *authservice.AuthService
}

type PostWithUser struct {
	*types.Post
	User         *types.GetUserResponse `json:"user"`
	LikeCount    int32                  `json:"like_count"`
	CommentCount int32                  `json:"comment_count"`
}

func NewGateway(postService *postservice.PostService, authservice *authservice.AuthService) *PostHandler {
	return &PostHandler{
		postclient: postService,
		auhClient:  authservice,
	}
}
func (h *PostHandler) HandleCreatePost(c *gin.Context) {
	// Ensure the request is a POST request
	userID, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	// Define a struct to bind JSON data
	var reqData struct {
		Caption  string   `json:"caption"`
		Tags     []string `json:"tags"`
		Mentions []string `json:"mentions"`
		Location string   `json:"location"`
		Image    string   `json:"image"`
	}

	// Bind JSON data from the request body
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON data: %v", err)})
		return
	}
	// Prepare the gRPC request
	if reqData.Caption == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Caption is required"})
		return
	}
	imageData, err := base64.StdEncoding.DecodeString(strings.Split(reqData.Image, ",")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to decode image: %v", err)})
		return
	}

	req := &types.CreatePostRequest{
		UserId:   userID.UserId,
		Caption:  reqData.Caption,
		Location: reqData.Location,
		Tags:     reqData.Tags,
		Mentions: reqData.Mentions,
		Media: []*types.MediaUpload{
			{
				FileData: []byte(imageData),
				FileName: "uploaded_image",
				FileType: "image/jpeg",
			},
		},
	}

	resp, err := h.postclient.CreatePost(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create post: %v", err)})
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, "Created Post Successfully", resp)
}

func (h *PostHandler) HandleGetPostByUser(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	// Get posts
	postsData, err := h.postclient.GetUserPosts(c, &types.GetUserPostsRequest{
		UserId:  user.UserId,
		Page:    int32(page),
		PerPage: 10,
	})
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get posts", err.Error())
		return
	}

	// Get user info
	userInfo, err := h.auhClient.GetUser(c, &types.GetUserRequest{
		UserId: user.UserId,
	})
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get user info", err.Error())
		return
	}

	// Combine posts with user info
	postsWithUser := make([]*PostWithUser, len(postsData.Posts))
	for i, post := range postsData.Posts {

		if err != nil {
			return
		}
		postsWithUser[i] = &PostWithUser{
			Post:         post,
			User:         userInfo,
			LikeCount:    0,
			CommentCount: 0,
		}
	}

	response.SendSuccessResponse(c, http.StatusOK, "Fetch Data Successfully", gin.H{
		"posts": postsWithUser,
	})
}
func (h *PostHandler) HandleGetAllPosts(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	// Get posts
	postsData, err := h.postclient.GetAllPosts(c, &types.GetAllPostsRequest{
		Page:    int32(page),
		PerPage: 10,
	})
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get posts", err.Error())
		return
	}

	// Get user info
	userInfo, err := h.auhClient.GetUser(c, &types.GetUserRequest{
		UserId: user.UserId,
	})
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to get user info", err.Error())
		return
	}

	// Combine posts with user info
	postsWithUser := make([]*PostWithUser, len(postsData.Posts))
	for i, post := range postsData.Posts {
		postsWithUser[i] = &PostWithUser{
			Post: post,
			User: userInfo,
		}
	}

	response.SendSuccessResponse(c, http.StatusOK, "Fetch Data Successfully", gin.H{
		"posts": postsWithUser,
	})
}
