package authhandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	authservice "github.com/wafi04/chatting-app/services/auth/pkg/service"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/pkg/response"
	"github.com/wafi04/chatting-app/services/shared/types"
)

type AuthHandler struct {
	authservice *authservice.AuthService
	logger      logger.Logger
}

func NewGateway(authservice *authservice.AuthService) *AuthHandler {
	return &AuthHandler{
		authservice: authservice,
	}
}

func (h *AuthHandler) HandleCreateUser(c *gin.Context) {
	log.Printf("Received create user request: %s %s", c.Request.Method, c.Request.URL.Path)

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create user : %v", err.Error())
		return
	}
	log.Printf("Decoded request: %+v", &req)
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	regis := &types.CreateUserRequest{
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		IpAddress:  clientIP,
		DeviceInfo: userAgent,
	}
	resp, err := h.authservice.CreateUser(c, regis)
	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	log.Printf("Received response from auth service: %+v", resp)

	response.SendSuccessResponse(c, http.StatusCreated, "User Created Successfully", resp)

}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {

	var req Login
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendErrorResponse(c, http.StatusBadRequest, "Failed to login")
		return
	}

	if req.Name == "" || req.Password == "" {
		response.SendErrorResponse(c, http.StatusBadRequest, "Credentials failed")
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	resp, err := h.authservice.Login(c.Request.Context(), &types.LoginRequest{
		Name:       req.Name,
		Password:   req.Password,
		DeviceInfo: userAgent,
		IpAddress:  clientIP,
	})

	if err != nil {
		response.SendErrorResponseWithDetails(c, http.StatusBadRequest, "user not found : %s", err.Error())
		return
	}
	c.Header("Access-Control-Allow-Credentials", "true")
	middleware.SetRefreshTokenCookie(c, resp.RefreshToken)
	middleware.SetAccressTokenCookie(c, resp.AccessToken)
	middleware.SetSessionCookie(c, resp.SessionInfo.SessionId)
	response.SendSuccessResponse(c, http.StatusOK, "Login user successfully", resp)
}

func (h *AuthHandler) HandleGetProfile(c *gin.Context) {
	h.logger.Log((logger.InfoLevel), "called profile")
	sessionID, err := c.Cookie("auth_session")

	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	users, err := h.authservice.GetSession(c, &types.GetSessionRequest{
		SessionId: sessionID,
	})

	if err != nil {
		response.Error(http.StatusBadRequest, "Failed to get profile")
		return
	}
	response.SendSuccessResponse(c, http.StatusOK, "Profile received successfully", users)
}

func (h *AuthHandler) HandleLogout(c *gin.Context) {
	user, err := middleware.GetUserFromGinContext(c)
	if err != nil {
		response.Error(http.StatusUnauthorized, "Unauthorized")
		return
	}
	token, err := c.Cookie("access-token")
	if err != nil {
		response.Error(http.StatusUnauthorized, "Unauthorized")
		return
	}

	logout, err := h.authservice.Logout(c, &types.LogoutRequest{
		AccessToken: token,
		UserId:      user.UserId,
	})

	if err != nil {
		response.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, "Logout Successfully", logout)
}
