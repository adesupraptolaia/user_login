package entity

// swagger:model
type UserProfile struct {
	UserKsuid   string `json:"user_ksuid,omitempty" gorm:"primaryKey"`
	Name        string `json:"name" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required,date=2006-01-02"`
	Address     string `json:"address" validate:"required"`
}

// swagger:model
type CreateUserRequest struct {
	UserProfile
	Username string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}
