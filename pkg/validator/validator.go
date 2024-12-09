package validator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Custom check to validate the provided password. Could not find an easy way to rely
// on the OpenAPI/Struct Validator and making our own validator would be the best way.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasLower = regexp.MustCompile(`[[:lower:]]`)
	var hasUpper = regexp.MustCompile(`[[:upper:]]`)
	var hasNumber = regexp.MustCompile(`[[:digit:]]`)
	var hasCharacters = regexp.MustCompile(`[[:graph:]]`)

	if !hasLower.MatchString(password) || !hasUpper.MatchString(password) || !hasNumber.MatchString(password) || !hasCharacters.MatchString(password) {
		return false
	}

	return true
}

func CheckConfig(c any, errormessage map[string]string) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("validatePassword", validatePassword)
	if err != nil {
		return err
	}

	err = validate.Struct(c)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			field := fieldError.StructField()
			tag := fieldError.Tag()
			errorKey := fmt.Sprintf("%s.%s", field, tag)

			if message, keyFound := errormessage[errorKey]; keyFound {
				return errors.New(message)
			}
		}
	}
	return nil
}
