package utils

import "github.com/go-playground/validator/v10"

// ValidateLoanState is a custom validator function to check if a loan is in the required state
func ValidateLoanState(fl validator.FieldLevel) bool {
	state := fl.Field().String()
	validStates := []string{"proposed", "approved", "invested", "disbursed"}

	for _, validState := range validStates {
		if state == validState {
			return true
		}
	}

	return false
}
