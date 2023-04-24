package repo

import (
	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UsersRepo interface {
	GetUserByKsuid(string) (*entity.User, error)
	GetUserByUsername(string) (*entity.User, error)
	CreateUser(entity.User) (*entity.User, error)
	UpdateUser(string, entity.User) (*entity.User, error)
	DeleteUser(string) (*entity.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) UsersRepo {
	return &userRepo{
		db: db.Table("users").Debug(),
	}
}

func (repo *userRepo) GetUserByKsuid(ksuid string) (*entity.User, error) {
	result := entity.User{}

	err := repo.db.
		Where("ksuid = ?", ksuid).
		First(&result).Error
	if err != nil {
		log.Errorf("error when GetUserByKsuid, err: %s", err.Error())
		return nil, err
	}

	return &result, err
}

func (repo *userRepo) GetUserByUsername(username string) (*entity.User, error) {
	result := entity.User{}

	err := repo.db.
		Where("username = ?", username).
		First(&result).Error
	if err != nil {
		log.Errorf("error when GetUserByUsername, err: %s", err.Error())
		return nil, err
	}

	return &result, err
}

func (repo *userRepo) CreateUser(user entity.User) (*entity.User, error) {
	err := repo.db.Create(&user).Error
	if err != nil {
		log.Errorf("error when CreateUser, err: %s", err.Error())
		return nil, err
	}

	return &user, err
}

func (repo *userRepo) UpdateUser(ksuid string, user entity.User) (*entity.User, error) {
	err := repo.db.
		Where("ksuid = ?", ksuid).
		Save(&user).Error
	if err != nil {
		log.Errorf("error when UpdateUser, err: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) DeleteUser(ksuid string) (*entity.User, error) {
	user := &entity.User{}

	err := repo.db.
		Where("ksuid = ? AND role = ?", ksuid, entity.USER).
		Delete(user).Error
	if err != nil {
		log.Errorf("error when DeleteUser, err: %s", err.Error())
		return nil, err
	}

	return user, nil
}
