package utils

type Response struct {
	StatusCode int         `json:"-"`
	ErrorMsg   string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type UpdateUserResponse struct {
}

func Success(statusCode int) *Response {
	return &Response{StatusCode: statusCode}
}

func Ok(statusCode int, data interface{}) *Response {
	return &Response{StatusCode: statusCode, Data: data}
}

func Error(statusCode int, errorMsg string) *Response {
	return &Response{StatusCode: statusCode, ErrorMsg: errorMsg}
}
