package helpers

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func GetValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "min":
		if err.Kind() == reflect.String {
			return err.Field() + " must be at least " + err.Param() + " characters"
		}
		return err.Field() + " must be at least " + err.Param()
	case "gt":
		return err.Field() + " must be greater than " + err.Param()
	default:
		return err.Field() + " is invalid"
	}
}
