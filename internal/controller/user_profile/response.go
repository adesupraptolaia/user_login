package user_profile_controller

import (
	"github.com/adesupraptolaia/user_login/internal/entity"
)

// swagger:model
type SuccessResp struct {
	// success
	Status string             `json:"status"`
	Data   entity.UserProfile `json:"data"`
}

// swagger:model
type ErrorResp struct {
	// error
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

func SuccessResponse(data *entity.UserProfile) SuccessResp {
	return SuccessResp{
		Status: "success",
		Data:   *data,
	}
}

func ErrorResponse(error_message string) ErrorResp {
	return ErrorResp{
		Status:       "error",
		ErrorMessage: error_message,
	}
}
