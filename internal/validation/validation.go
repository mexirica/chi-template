// Package validation provides functions and types for input validation and business rule enforcement.
package validation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Init initializes the validator instance for struct validation.
func Init() {
	validate = validator.New()
}

type ValidationErrorResponse struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value,omitempty"`
}

// BindAndValidate decodes the request body into the destination struct and validates it.
// Returns a slice of validation errors and an error if validation fails.
func BindAndValidate(r *http.Request, dst any) ([]ValidationErrorResponse, error) {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return nil, err
	}

	err := validate.Struct(dst)
	if err == nil {
		return nil, nil
	}

	var errors []ValidationErrorResponse
	for _, err := range err.(validator.ValidationErrors) {
		var val string
		if err.Param() != "" {
			val = err.Param()
		}
		errors = append(errors, ValidationErrorResponse{
			FailedField: err.StructNamespace(),
			Tag:         err.Tag(),
			Value:       val,
		})
	}

	return errors, fmt.Errorf("validation failed")
}
