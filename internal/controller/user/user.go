package user_controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/internal/usecase"
	"github.com/adesupraptolaia/user_login/internal/utils"
	"github.com/adesupraptolaia/user_login/pkg/jwt"
	"github.com/adesupraptolaia/user_login/pkg/validator"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uc usecase.UserUC
}

func NewUserHandler(uc usecase.UserUC) UserHandler {
	return UserHandler{
		uc: uc,
	}
}

// Login godoc
// @Summary Login
// @Description Login using username and password
// @Tags Public
// @Accept  json
// @Produce  json
// @Param payload body entity.UserRequest true "payload"
// @Success 200 {object} TokenSuccessResp
// @Response 400 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /login [post]
func (h UserHandler) Login(ctx echo.Context) error {
	req := entity.UserRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	user, err := h.uc.GetUserByUsername(req.Username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	// check password
	if !utils.IsPasswordMatch(user.Password, req.Password) {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("wrong username or password"))
	}

	accessToken, err := jwt.CreateAccessToken(user.Ksuid, user.Role)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	refreshToken, err := jwt.CreateRefreshToken(user.Ksuid, user.Role)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusCreated, SuccessTokenResponse(user.Ksuid, accessToken, refreshToken))
}

// RefreshAccessToken godoc
// @Summary Refresh Token
// @Description Refresh AccessToken
// @Tags Public
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {refresh_token}"
// @Success 200 {object} TokenSuccessResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /refresh [get]
func (h UserHandler) RefreshToken(ctx echo.Context) error {
	refreshToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	user, err := h.uc.GetUserByKsuid(claims.UserKsuid)
	if err != nil || user.Role != claims.Role {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	accessToken, err := jwt.CreateAccessToken(user.Ksuid, user.Role)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusOK, SuccessTokenResponse(user.Ksuid, accessToken, refreshToken))
}

// CreateUser godoc
// @Summary Create new User
// @Description Only admin can create new user
// @Tags Private
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param payload body entity.UserRequest true "Request Payload"
// @Success 201 {object} UserSuccessResp
// @Response 400 {object} ErrorResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /user/create [post]
func (h UserHandler) CreateUser(ctx echo.Context) error {
	userProfile := entity.UserRequest{}
	err := ctx.Bind(&userProfile)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	accessToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	err = jwt.IsAdmin(accessToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	err = validator.ValidateStruct(userProfile)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
	}

	newUser, err := h.uc.CreateUser(userProfile)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	newUser.Password = ""

	return ctx.JSON(http.StatusCreated, SuccessResponse(newUser))
}

// DeleteUser godoc
// @Summary Delete User
// @Description Only admin can delete user
// @Tags Private
// @Accept  json
// @Produce  json
// @Param ksuid path string true "Ksuid of User"
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} UserSuccessResp
// @Response 401 {object} ErrorResp
// @Response 500 {object} ErrorResp
// @Router /user/{user_ksuid} [delete]
func (h UserHandler) DeleteUser(ctx echo.Context) error {
	ksuid := ctx.Param("ksuid")

	accessToken, err := getBearerToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	err = jwt.IsAdmin(accessToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorize"))
	}

	deletedUser, err := h.uc.DeleteUser(ksuid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
	}

	deletedUser.Password = ""

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
