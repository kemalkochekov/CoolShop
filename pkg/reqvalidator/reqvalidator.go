package reqvalidator

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validate            *validator.Validate
	phoneRegex          = `^(\+[1-9]\d{9,14}|\d{10,15})$`
	emailRegex          = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	compiledPhoneRegexp *regexp.Regexp
	compiledEmailRegexp *regexp.Regexp
)

func init() {
	validate = validator.New()

	compiledPhoneRegexp = regexp.MustCompile(phoneRegex)
	compiledEmailRegexp = regexp.MustCompile(emailRegex)

	err := validate.RegisterValidation("phone", validatePhoneNumber)
	if err != nil {
		log.Printf("[reqvalidator][init] Unable to put validator for phonuNumber %v", err)
	}

	err = validate.RegisterValidation("email", validateEmail)
	if err != nil {
		log.Printf("[reqvalidator][init] Unable to put validator for email %v", err)
	}
}

// ValidateRequest is
func ValidateRequest(request interface{}) error {
	return validate.Struct(request)
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()

	return compiledPhoneRegexp.MatchString(phoneNumber)
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	return compiledEmailRegexp.MatchString(email)
}
