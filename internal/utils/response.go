package utils

type Response struct {
	StatusCode int         `json:"-"`
	ErrorMsg   string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Status     bool        `json:"status,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type GetUserResponse struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Avatar   string `json:"avatar,omitempty"`
}

func Success(statusCode int) *Response {
	return &Response{StatusCode: statusCode, Status: true}
}

func Fail(statusCode int) *Response {
	return &Response{StatusCode: statusCode, Status: false}
}

func Ok(statusCode int, data interface{}) *Response {
	return &Response{StatusCode: statusCode, Data: data, Status: true}
}

func Error(statusCode int, errorMsg string) *Response {
	return &Response{StatusCode: statusCode, ErrorMsg: errorMsg, Status: false}
}
