package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/middleware"
	"github.com/xyz-corp/xyz-football-api/internal/service"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register godoc
// @Summary      Registrasi admin
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Payload registrasi"
// @Success      201 {object} response.Envelope{data=dto.UserResponse}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}
	user, err := h.svc.Register(req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.Created(c, "Registrasi berhasil", user)
}

// Login godoc
// @Summary      Login admin
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Kredensial login"
// @Success      200 {object} response.Envelope{data=dto.LoginResponse}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.svc.Login(req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Login berhasil", result)
}

// Refresh godoc
// @Summary      Perpanjang access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh token"
// @Success      200 {object} response.Envelope{data=dto.LoginResponse}
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.svc.Refresh(req)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Token berhasil diperpanjang", result)
}

// Me godoc
// @Summary      Profil user login
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Envelope{data=dto.UserResponse}
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userUUID, exists := c.Get(middleware.ContextUserUUIDKey)
	if !exists {
		handleErr(c, apperror.Unauthorized("Tidak terautentikasi"))
		return
	}
	user, err := h.svc.Me(userUUID.(string))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, "Profil berhasil diambil", user)
}
