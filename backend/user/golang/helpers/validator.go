package helpers

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func UsernameValidator(validate *validator.Validate) {
	validate.RegisterValidation("usernamevalidator", func(fl validator.FieldLevel) bool {
		usernameRegex := `^[a-zA-Z\d]{5,12}$`
		return regexp.MustCompile(usernameRegex).MatchString(fl.Field().String())
	})
}

func PasswordValidator(validate *validator.Validate) {
	validate.RegisterValidation("passwordvalidator", func(fl validator.FieldLevel) bool {
		passwordRegex := `^[a-zA-Z\d@_-]{8,20}$`
		password := fl.Field().String()
		ok := regexp.MustCompile(passwordRegex).MatchString(password)
		if !ok {
			return false
		}

		isSpesialCharacter := strings.ContainsAny(password, "@ | _ | -")

		isUpper := false
		isLower := false
		isNumber := false
	isPasswordLoop:
		for _, value := range password {
			if unicode.IsUpper(value) && unicode.IsLetter(value) && !isUpper {
				isUpper = true
			} else if unicode.IsLower(value) && unicode.IsLetter(value) && !isLower {
				isLower = true
			} else if _, err := strconv.Atoi(string(value)); err == nil {
				isNumber = true
			}

			if isUpper && isLower && isNumber {
				break isPasswordLoop
			}
		}

		if !isSpesialCharacter || !isUpper || !isLower || !isNumber {
			return false
		}
		return true
	})
}

func TelephoneValidator(validate *validator.Validate) {
	validate.RegisterValidation("telephonevalidator", func(fl validator.FieldLevel) bool {
		regexString := `^[\d+]{14}$`
		return regexp.MustCompile(regexString).MatchString(fl.Field().String())
	})
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetValidatorError(validatorError error, structRequest interface{}) (errorMessages []ErrorMessage) {
	validationErrors := validatorError.(validator.ValidationErrors)
	val := reflect.ValueOf(structRequest)
	for _, fieldError := range validationErrors {
		var errorMessage ErrorMessage
		structField, ok := val.Type().FieldByName(fieldError.Field())
		if !ok {
			errorMessage.Field = "property"
			errorMessage.Message = "couldn't find property: " + fieldError.Field()
			errorMessages = append(errorMessages, errorMessage)
			return
		}
		errorMessage.Field = structField.Tag.Get("json")
		if fieldError.Tag() == "usernamevalidator" {
			errorMessage.Message = "please use only uppercase and lowercase letter and number and min 5 and max 8 alphanumeric"
		} else if fieldError.Tag() == "passwordvalidator" {
			errorMessage.Message = "please use only uppercase, lowercase, number and must have 1 uppercase. lowercase, number, @, _, -, min 8 and max 20"
		} else if fieldError.Tag() == "telephonevalidator" {
			errorMessage.Message = "please use only number and + "
		} else if fieldError.Tag() == "email" {
			errorMessage.Message = "please input a correct email format "
		} else if fieldError.Tag() == "gte" {
			errorMessage.Message = "please input greater than equal to " + fieldError.Param()
		} else {
			errorMessage.Message = "is " + fieldError.Tag()
		}
		errorMessages = append(errorMessages, errorMessage)
	}
	return
}
