package user_controller

import (
	"github.com/adesupraptolaia/user_login/internal/entity"
)

// swagger:model
type UserSuccessResp struct {
	// success
	Status string       `json:"status"`
	Data   *entity.User `json:"data"`
}

// swagger:model
type TokenData struct {
	Ksuid        string `json:"ksuid,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// swagger:model
type TokenSuccessResp struct {
	// success
	Status string    `json:"status"`
	Data   TokenData `json:"data"`
}

// swagger:model
type ErrorResp struct {
	// error
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

func SuccessTokenResponse(ksuid, accessToken, refreshToken string) TokenSuccessResp {
	return TokenSuccessResp{
		Status: "success",
		Data: TokenData{
			Ksuid:        ksuid,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func SuccessResponse(data *entity.User) UserSuccessResp {
	return UserSuccessResp{
		Status: "success",
		Data:   data,
	}
}

func ErrorResponse(error_message string) ErrorResp {
	return ErrorResp{
		Status:       "error",
		ErrorMessage: error_message,
	}
}
