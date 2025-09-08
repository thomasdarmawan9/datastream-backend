package dto

type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"password123"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"newuser"`
	Password string `json:"password" example:"newpassword"`
	Role     string `json:"role" example:"viewer"`
}
