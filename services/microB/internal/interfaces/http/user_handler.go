package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/dto"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/auth"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/usecase"
)

type UserHandler struct {
	uc         usecase.UserUsecase
	jwtManager *auth.JWTManager
}

func NewUserHandler(e *echo.Echo, uc usecase.UserUsecase, jwt *auth.JWTManager) {
	handler := &UserHandler{uc: uc, jwtManager: jwt}

	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.uc.Register(req.Username, req.Password, req.Role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "registered"})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) Login(c echo.Context) error {

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	user, err := h.uc.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	token, err := h.jwtManager.Generate(user.Username, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
