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
	Code     string `json:"code" binding:"required"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// Code     string `json:"code" binding:"required"`
}

// type UploadAvatarRequest struct {
// 	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
// }

type UpdatePasswordRequest struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type SendEmailVerificationRequest struct {
	Email string `json:"email" binding:"required" validate:"email"`
}
