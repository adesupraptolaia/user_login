package entity

// swagger:model
type User struct {
	Ksuid    string `json:"ksuid"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Role     string `json:"role"`
}

// swagger:model
type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

const (
	ADMIN string = "admin"
	USER  string = "user"
)
