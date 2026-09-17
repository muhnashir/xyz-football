package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
)

var instance = validator.New()

// Details converts a go-playground validator error into apperror.Detail list
// with human readable (Indonesian) messages for the response envelope.
func Details(err error) []apperror.Detail {
	var details []apperror.Detail
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range verrs {
			details = append(details, apperror.Detail{
				Field:   strings.ToLower(fe.Field()),
				Message: message(fe),
			})
		}
		return details
	}
	details = append(details, apperror.Detail{Field: "", Message: err.Error()})
	return details
}

func message(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "min":
		return fmt.Sprintf("%s minimal %s", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s", field, fe.Param())
	case "email":
		return fmt.Sprintf("%s harus berupa email yang valid", field)
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", field, fe.Param())
	case "uuid":
		return fmt.Sprintf("%s harus berupa uuid yang valid", field)
	case "gtfield", "nefield":
		return fmt.Sprintf("%s tidak valid", field)
	default:
		return fmt.Sprintf("%s tidak valid (%s)", field, fe.Tag())
	}
}

func Struct(s interface{}) error {
	return instance.Struct(s)
}
