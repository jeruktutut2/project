package helper

import (
	"reflect"
	"regexp"
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
	isPasswordLoop:
		for _, value := range password {
			if unicode.IsUpper(value) && unicode.IsLetter(value) && !isUpper {
				isUpper = true
			}

			if unicode.IsLower(value) && unicode.IsLetter(value) && !isLower {
				isLower = true
			}

			if isUpper && isLower {
				break isPasswordLoop
			}
		}

		if !isSpesialCharacter || !isUpper || !isLower {
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

func GetValidatorError(validatorError error, structRequest interface{}) (result map[string]interface{}) {
	validationErrors := validatorError.(validator.ValidationErrors)
	val := reflect.ValueOf(structRequest)
	result = make(map[string]interface{})
validationErrorLoop:
	for _, fieldError := range validationErrors {
		structField, ok := val.Type().FieldByName(fieldError.Field())
		// fmt.Println("structField:", structField, ok)
		if !ok {
			// result = nil
			result["property"] = "couldn't find property: " + fieldError.Field()
			return
		}
		structJsonTag := structField.Tag.Get("json")
		if fieldError.Tag() == "usernamevalidator" {
			result[structJsonTag] = "please use only uppercase and lowercase letter and number and min 5 and max 8 alphanumeric"
			continue validationErrorLoop
		} else if fieldError.Tag() == "passwordvalidator" {
			result[structJsonTag] = "please use only uppercase, lowercase, number and must have 1 uppercase. lowercase, number, @, _, -, min 8 and max 20"
			continue validationErrorLoop
		} else if fieldError.Tag() == "telephonevalidator" {
			result[structJsonTag] = "please use only number and + "
			continue validationErrorLoop
		} else if fieldError.Tag() == "email" {
			result[structJsonTag] = "please input a correct email format "
			continue validationErrorLoop
		} else {
			result[structJsonTag] = "is " + fieldError.Tag()
			continue validationErrorLoop
		}

		// result[structJsonTag] = fieldError.Tag()
	}
	return
}
