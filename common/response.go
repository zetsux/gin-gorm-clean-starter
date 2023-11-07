package common

type Response struct {
	IsSuccess bool   `json:"success"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Status    uint   `json:"status"`
	Data      any    `json:"data"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type EmptyObj struct {
}

func CreateFailResponse(msg string, err string, statusCode uint) Response {
	return Response{
		IsSuccess: false, Message: msg, Error: err, Status: statusCode, Data: nil,
	}
}

func CreateSuccessResponse(msg string, statusCode uint, d any) Response {
	return Response{
		IsSuccess: true, Message: msg, Status: statusCode, Data: d,
	}
}

func CreateAuthResponse(token string, role string) AuthResponse {
	return AuthResponse{
		Token: token, Role: role,
	}
}
