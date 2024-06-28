package base

type Response struct {
	IsSuccess bool                `json:"success"`
	Message   string              `json:"message"`
	Error     string              `json:"error,omitempty"`
	Status    uint                `json:"status"`
	Data      any                 `json:"data"`
	Meta      *PaginationResponse `json:"meta,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type PaginationResponse struct {
	Page     int64 `json:"page"`
	PerPage  int64 `json:"per_page"`
	LastPage int64 `json:"last_page"`
	Total    int64 `json:"total"`
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

func CreatePaginatedResponse(msg string, statusCode uint, d any, pageMeta PaginationResponse) Response {
	return Response{
		IsSuccess: true, Message: msg, Status: statusCode, Data: d, Meta: &pageMeta,
	}
}

func CreateAuthResponse(token string, role string) AuthResponse {
	return AuthResponse{
		Token: token, Role: role,
	}
}
