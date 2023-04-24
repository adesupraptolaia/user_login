package repo

import (
	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserProfilesRepo interface {
	GetUserProfile(string) (*entity.UserProfile, error)
	CreateUserProfile(entity.UserProfile) (*entity.UserProfile, error)
	UpdateUserProfile(entity.UserProfile) (*entity.UserProfile, error)
	DeleteUserProfile(string) (*entity.UserProfile, error)
}

type userProfileRepo struct {
	db *gorm.DB
}

func NewUserProfile(db *gorm.DB) UserProfilesRepo {
	return &userProfileRepo{
		db: db.Table("user_profiles").Debug(),
	}
}

func (repo *userProfileRepo) GetUserProfile(userKsuid string) (*entity.UserProfile, error) {
	result := entity.UserProfile{UserKsuid: userKsuid}

	err := repo.db.First(&result).Error
	if err != nil {
		log.Errorf("error when GetUserProfile, err: %s", err.Error())
		return nil, err
	}

	return &result, err
}

func (repo *userProfileRepo) CreateUserProfile(userProfile entity.UserProfile) (*entity.UserProfile, error) {
	err := repo.db.Create(&userProfile).Error
	if err != nil {
		log.Errorf("error when CreateUserProfile, err: %s", err.Error())
		return nil, err
	}

	return &userProfile, err
}

func (repo *userProfileRepo) UpdateUserProfile(userProfile entity.UserProfile) (*entity.UserProfile, error) {
	err := repo.db.Save(&userProfile).Error
	if err != nil {
		log.Errorf("error when UpdateUserProfile+, err: %s", err.Error())
		return nil, err
	}

	return &userProfile, nil
}

func (repo *userProfileRepo) DeleteUserProfile(userKsuid string) (*entity.UserProfile, error) {
	userProfile, err := repo.GetUserProfile(userKsuid)
	if err != nil {
		return nil, err
	}

	err = repo.db.Delete(&entity.UserProfile{UserKsuid: userKsuid}).Error
	if err != nil {
		log.Errorf("error when DeleteUserProfile, err: %s", err.Error())
		return nil, err
	}

	return userProfile, nil
}
