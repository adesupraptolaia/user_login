package validator

import (
	"errors"
	"fmt"
	"time"

	validator_lib "github.com/go-playground/validator/v10"
)

var validator *validator_lib.Validate

func init() {
	validator = validator_lib.New()

	validator.RegisterValidation("date", func(fl validator_lib.FieldLevel) bool {
		_, err := time.Parse("2006-01-02", fl.Field().String())
		return err == nil
	})
}

func ValidateStruct(data interface{}) error {
	err := validator.Struct(data)
	if err != nil {
		var errMsg string
		if _, ok := err.(*validator_lib.InvalidValidationError); ok {
			return errors.New("bad request")
		}

		for _, err := range err.(validator_lib.ValidationErrors) {
			errMsg += fmt.Sprintf("%s is %s, ", err.StructField(), err.Tag())
		}

		return errors.New(errMsg)
	}

	return nil
}
