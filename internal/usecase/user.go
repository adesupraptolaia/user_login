package usecase

import (
	"fmt"

	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/internal/repo"
	"github.com/adesupraptolaia/user_login/internal/utils"
	"github.com/segmentio/ksuid"
)

type UserUC interface {
	GetUserByUsername(string) (*entity.User, error)
	GetUserByKsuid(string) (*entity.User, error)
	CreateUser(entity.UserRequest) (*entity.User, error)
	UpdateUser(string, entity.User) (*entity.User, error)
	DeleteUser(string) (*entity.User, error)
}

type user struct {
	repo repo.UsersRepo
}

func NewUser(repo repo.UsersRepo) UserUC {
	return &user{
		repo: repo,
	}
}

func (uc *user) GetUserByUsername(username string) (*entity.User, error) {
	user, err := uc.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user with username %s not found", username)
	}

	return user, nil
}

func (uc *user) GetUserByKsuid(ksuid string) (*entity.User, error) {
	user, err := uc.repo.GetUserByKsuid(ksuid)
	if err != nil {
		return nil, fmt.Errorf("user with ksuid %s not found", ksuid)
	}

	return user, nil
}

func (uc *user) CreateUser(userReq entity.UserRequest) (*entity.User, error) {
	user, _ := uc.repo.GetUserByUsername(userReq.Username)
	if user != nil {
		return nil, fmt.Errorf("user with username %s already exist", user.Username)
	}

	data := entity.User{
		Ksuid:    ksuid.New().String(),
		Username: userReq.Username,
		Password: utils.HashPassword(userReq.Password),
		Role:     entity.USER,
	}

	user, err := uc.repo.CreateUser(data)
	if err != nil {
		return nil, fmt.Errorf("failed when create user_profiles")
	}

	return user, nil
}

func (uc *user) UpdateUser(ksuid string, user entity.User) (*entity.User, error) {
	if _, err := uc.repo.GetUserByKsuid(ksuid); err != nil {
		return nil, fmt.Errorf("user with ksuid %s not exist", ksuid)
	}

	user.Ksuid = ksuid

	userResp, err := uc.repo.UpdateUser(ksuid, user)
	if err != nil {
		return nil, fmt.Errorf("failed when update user_profiles with ksuid %s", ksuid)
	}

	return userResp, nil
}

func (uc *user) DeleteUser(ksuid string) (*entity.User, error) {
	user, err := uc.repo.GetUserByKsuid(ksuid)
	if err != nil {
		return nil, fmt.Errorf("user with ksuid %s not exist", ksuid)
	}

	_, err = uc.repo.DeleteUser(ksuid)
	if err != nil {
		return nil, fmt.Errorf("failed when delete user_profiles with ksuid %s", ksuid)
	}

	return user, nil
}
