package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// TODO
type CustomValidator interface {
	ValidateUseCase(payload interface{}) *UseCaseValError
}

type customValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() CustomValidator {
	validate := validator.New()
	validate.RegisterValidation("timeformat", timeFormatValidation)
	return &customValidator{Validator: validate}
}

type ValidationError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func timeFormatValidation(fl validator.FieldLevel) bool {
	layout := time.RFC3339 // "2006-01-02T15:04:05Z07:00"
	value := fl.Field().String()
	_, err := time.Parse(layout, value)
	return err == nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", strings.ToUpper(err.Field()))
	case "timeformat":
		return fmt.Sprintf("'%s' must be a valid time format (example: %s)", err.Field(), time.RFC3339)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

type UseCaseValError struct {
	ValidationErros []ValidationError
	ErrorType       string
}

func (cv *customValidator) ValidateUseCase(payload interface{}) *UseCaseValError {
	var validationErrors []ValidationError

	err := cv.Validator.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Rule:    err.Tag(),
				Message: getErrorMessage(err),
			})
		}

		return &UseCaseValError{
			ValidationErros: validationErrors,
			ErrorType:       "validation error",
		}
	}

	return nil
}

func (e *UseCaseValError) Error() string {
	return (fmt.Sprintf(e.ErrorType))
}

func (e *UseCaseValError) GetValidationErrors() []ValidationError {
	return e.ValidationErros
}
