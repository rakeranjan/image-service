package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// func RegisterValidator() {
// 	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
// 		v.RegisterValidation("userName", ValidateUserName)
// 		v.RegisterValidation("phoneNmber", ValidatePhoneNumber)
// 	}
// }

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	regex := `^\d{10}$`
	phoneNmber := fl.Field().String()
	return regexp.MustCompile(regex).MatchString(phoneNmber)
}

func ValidateUserName(fl validator.FieldLevel) bool {
	regex := `^[a-zA-Z][a-zA-Z0-9._-]{2,15}$`
	username := fl.Field().String()
	return regexp.MustCompile(regex).MatchString(username)
}
