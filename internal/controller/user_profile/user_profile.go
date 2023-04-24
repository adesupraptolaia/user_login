package user_profile_controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/internal/usecase"
	"github.com/adesupraptolaia/user_login/pkg/jwt"
	"github.com/adesupraptolaia/user_login/pkg/validator"
	"github.com/labstack/echo/v4"
)

type userProfileHandler struct {
	uc usecase.UserProfileUC
}

func NewUserProfileHandler(uc usecase.UserProfileUC) userProfileHandler {
	return userProfileHandler{
		uc: uc,
	}
}

// GetUser godoc
// @Summary Get a user by Userksuid
// @Description Get a user by Userksuid
// @Tags users
// @Accept  json
// @Produce  json
// @Param user_ksuid path string true "Ksuid of User"
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} SuccessResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /user/{user_ksuid} [get]
func (h userProfileHandler) GetUser(ctx echo.Context) error {
	userKsuid := ctx.Param("user_ksuid")

	accessToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	// Admin and User can access this API
	if err = jwt.ValidateAccessToken(accessToken, userKsuid); err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	newUser, err := h.uc.GetUserProfile(userKsuid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusOK, SuccessResponse(newUser))
}

// CreateUser godoc
// @Summary Create New User
// @Description Only admin can create new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param payload body entity.CreateUserRequest true "Request Payload"
// @Success 201 {object} SuccessResp
// @Response 400 {object} ErrorResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /user/create [post]
func (h userProfileHandler) CreateUser(ctx echo.Context) error {
	userProfile := entity.CreateUserRequest{}
	if err := ctx.Bind(&userProfile); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	accessToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	if err = jwt.IsAdmin(accessToken); err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	err = validator.ValidateStruct(userProfile)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	newUser, err := h.uc.CreateUserProfile(userProfile)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusCreated, SuccessResponse(newUser))
}

// UpdateUser godoc
// @Summary Update User Profile
// @Description Only admin can update user profile
// @Tags users
// @Accept  json
// @Produce  json
// @Param user_ksuid path string true "Ksuid of User"
// @Param Authorization header string true "Bearer {token}"
// @Param payload body entity.UserProfile true "Request Payload"
// @Success 200 {object} SuccessResp
// @Response 400 {object} ErrorResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /user/{user_ksuid}/update [post]
func (h userProfileHandler) UpdateUser(ctx echo.Context) error {
	userKsuid := ctx.Param("user_ksuid")

	userProfile := entity.UserProfile{}
	if err := ctx.Bind(&userProfile); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	accessToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	if err = jwt.IsAdmin(accessToken); err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	err = validator.ValidateStruct(userProfile)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	updatedUser, err := h.uc.UpdateUserProfile(userKsuid, userProfile)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusOK, SuccessResponse(updatedUser))
}

// DeleteUser godoc
// @Summary Delete User
// @Description Only admin can delete user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user_ksuid path string true "Ksuid of User"
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} SuccessResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /user/{user_ksuid} [delete]
func (h userProfileHandler) DeleteUser(ctx echo.Context) error {
	userKsuid := ctx.Param("user_ksuid")

	accessToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	if err = jwt.IsAdmin(accessToken); err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	deletedUser, err := h.uc.DeleteUserProfile(userKsuid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusOK, SuccessResponse(deletedUser))
}

func getBearerToken(auth string) (string, error) {
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("missing bearer token")
	}

	token := strings.Replace(auth, "Bearer ", "", -1)

	if token == "" {
		return "", fmt.Errorf("missing bearer token")
	}

	return token, nil
}
