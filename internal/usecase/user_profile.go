package usecase

import (
	"fmt"
	"time"

	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/internal/repo"
)

type UserProfileUC interface {
	GetUserProfile(string) (*entity.UserProfile, error)
	CreateUserProfile(entity.CreateUserRequest) (*entity.UserProfile, error)
	UpdateUserProfile(string, entity.UserProfile) (*entity.UserProfile, error)
	DeleteUserProfile(string) (*entity.UserProfile, error)
}

type userProfile struct {
	repo repo.UserProfilesRepo
	auth repo.AuthRepo
}

func NewUserProfile(userProfileRepo repo.UserProfilesRepo, authRepo repo.AuthRepo) UserProfileUC {
	return &userProfile{
		repo: userProfileRepo,
		auth: authRepo,
	}
}

func (uc *userProfile) GetUserProfile(userksuid string) (*entity.UserProfile, error) {
	userProfile, err := uc.repo.GetUserProfile(userksuid)
	if err != nil {
		return nil, fmt.Errorf("user with ksuid %s not found", userksuid)
	}

	userProfile.DateOfBirth = convertDatetime(userProfile.DateOfBirth)

	return userProfile, nil
}
func (uc *userProfile) CreateUserProfile(userProfileReq entity.CreateUserRequest) (*entity.UserProfile, error) {
	userRequest := entity.User{
		Username: userProfileReq.Username,
		Password: userProfileReq.Password,
	}

	user, err := uc.auth.CreateUser(userRequest)
	if err != nil {
		return nil, fmt.Errorf("error when create user to auth service %s", err.Error())
	}

	data := entity.UserProfile{
		UserKsuid:   user.Ksuid,
		Name:        userProfileReq.Name,
		DateOfBirth: userProfileReq.DateOfBirth,
		Address:     userProfileReq.Address,
	}

	userProfile, err := uc.repo.CreateUserProfile(data)

	if err != nil {
		deleteUserInAuthService(user.Ksuid, uc.auth)

		return nil, fmt.Errorf("failed when create user_profiles")
	}

	return userProfile, nil
}

func (uc *userProfile) UpdateUserProfile(userKsuid string, userProfile entity.UserProfile) (*entity.UserProfile, error) {
	_, err := uc.repo.GetUserProfile(userKsuid)
	if err != nil {
		return nil, fmt.Errorf("user_profiles with ksuid %s not found", userKsuid)
	}

	userProfile.UserKsuid = userKsuid

	userProfileResp, err := uc.repo.UpdateUserProfile(userProfile)
	if err != nil {
		return nil, fmt.Errorf("failed when update user_profiles with ksuid %s", userKsuid)
	}

	return userProfileResp, nil
}

func (uc *userProfile) DeleteUserProfile(userKsuid string) (*entity.UserProfile, error) {
	_, err := uc.repo.GetUserProfile(userKsuid)
	if err != nil {
		return nil, fmt.Errorf("user_profiles with ksuid %s not found", userKsuid)
	}

	if _, err := uc.auth.DeleteUser(userKsuid); err != nil {
		return nil, fmt.Errorf("failed when delete user to auth_service with ksuid %s", userKsuid)
	}

	deletedUser, err := uc.repo.DeleteUserProfile(userKsuid)
	if err != nil {
		return nil, fmt.Errorf("failed when delete user_profiles with ksuid %s", userKsuid)
	}

	deletedUser.DateOfBirth = convertDatetime(deletedUser.DateOfBirth)

	return deletedUser, nil
}

func convertDatetime(dt string) string {
	t, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return dt
	}
	return t.Format("2006-01-02")
}

func deleteUserInAuthService(ksuid string, auth repo.AuthRepo) {
	for maxRetry := 5; maxRetry > 0; maxRetry-- {
		if _, err := auth.DeleteUser(ksuid); err == nil {
			return
		}

	}
}
