package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse map[string]string

var validate = validator.New()

func ValidateRequest(s interface{}) ErrorResponse {
	errors := ErrorResponse{}

	// Check for validation errors
	if err := validate.Struct(s); err != nil {
		// Cast error to ValidationErrors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				// Create a simplified error map: "field": "tag_failed"
				// e.g., "email": "required", "password": "min"
				field := e.Field() // Name of the field
				tag := e.Tag()     // The rule that failed (e.g. required, min, email)
				param := e.Param() // The param (e.g. 7 for min=7)

				msg := fmt.Sprintf("failed on rule '%s'", tag)
				if param != "" {
					msg = fmt.Sprintf("failed on rule '%s' (param: %s)", tag, param)
				}

				errors[field] = msg
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}
