package helper

import (
	"encoding/json"
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

type Result struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetValidatorError(validatorError error, structRequest interface{}) (results []Result) {
	validationErrors := validatorError.(validator.ValidationErrors)
	val := reflect.ValueOf(structRequest)
	for _, fieldError := range validationErrors {
		var r Result
		structField, ok := val.Type().FieldByName(fieldError.Field())
		if !ok {
			r.Field = "property"
			r.Message = "couldn't find property: " + fieldError.Field()
			results = append(results, r)
			return
		}
		r.Field = structField.Tag.Get("json")
		if fieldError.Tag() == "usernamevalidator" {
			r.Message = "please use only uppercase and lowercase letter and number and min 5 and max 8 alphanumeric"
		} else if fieldError.Tag() == "passwordvalidator" {
			r.Message = "please use only uppercase, lowercase, number and must have 1 uppercase. lowercase, number, @, _, -, min 8 and max 20"
		} else if fieldError.Tag() == "telephonevalidator" {
			r.Message = "please use only number and + "
		} else if fieldError.Tag() == "email" {
			r.Message = "please input a correct email format "
		} else {
			r.Message = "is " + fieldError.Tag()
		}
		results = append(results, r)
	}
	return
}

func ToResultsMessageResponse(requestId string, field string, message string) (response string, err error) {
	var results []Result
	var result Result
	result.Field = field
	result.Message = message
	results = append(results, result)
	var resultsByte []byte
	resultsByte, err = json.Marshal(results)
	if err != nil {
		return
	}
	response = string(resultsByte)
	return
}
