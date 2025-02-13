package comments

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/pkg/response"
	"github.com/wafi04/chatting-app/services/shared/types"
)

type CommentHandler struct {
	srv *CommentService
}

func NewCommntHandler(srv *CommentService) *CommentHandler {
	return &CommentHandler{
		srv: srv,
	}
}

func (h *CommentHandler) HandleCreateComment(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)

	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		PostID   string `json:"post_id"`
		Content  string `json:"content"`
		ParentID string `json:"parent_comment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed input", err.Error())
		return
	}

	data, err := h.srv.CreateComment(c, &types.CreateComment{
		UserID:   user.UserId,
		PostID:   req.PostID,
		Content:  req.Content,
		ParentID: &req.ParentID,
	})

	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create comment", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, "Create Comment Successfullt", data)
}

func (h *CommentHandler) HandleGetComments(c *gin.Context) {
	postID := c.Param("id")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	parentID := c.Request.URL.Query().Get("parent_id")
	includeChildren := c.Request.URL.Query().Get("include_children") == "true"

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	req := &types.ListCommentsRequest{
		Page:            int32(page),
		Limit:           int32(limit),
		IncludeChildren: includeChildren,
		PostID:          postID,
	}

	if parentID != "" {
		req.ParentID = &parentID
	}

	data, err := h.srv.GetComments(c, req)

	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to Get comment", err.Error())
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Get Comment Successfullt", data)
}

func (h *CommentHandler) HandleDeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Alternative ways to get parameters:
	// Query parameter (?id=123): c.Query("id")
	// Optional query parameter with default: c.DefaultQuery("id", "default_value")
	// Check if query exists: id, exists := c.GetQuery("id")

	// Validate ID
	if id == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Category ID is required")
		return
	}
	updateReq := &types.DeleteComment{
		CommentID:      id,
		DeleteChildren: true,
	}

	category, err := h.srv.DeleteComment(c, updateReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			response.SendErrorResponse(c, http.StatusNotFound, "Category Not Found")
		case strings.Contains(err.Error(), "invalid"):
			response.SendErrorResponse(c, http.StatusBadRequest, "Invalid Request")
		default:
			response.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")

		}
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Delete Category Succesfully", category)
}
