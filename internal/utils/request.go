package utils

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   uint8  `json:"status"`
	Role     uint8  `json:"role"`
}