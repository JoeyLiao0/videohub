package user

type UserInfo struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Avatar   string `json:"avatar,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type GetUserResponse struct {
	User UserInfo `json:"user" binding:"required"`
}
