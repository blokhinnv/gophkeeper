// Package validation provides utility functions for validating data structures.
package validation

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// Validate is a validator instance that can be used to validate data structures.
var Validate = validator.New()

// validateExpirationDate is a custom validation function that checks if a string
// has the format "MM/YY".
func validateExpirationDate(fl validator.FieldLevel) bool {
	// Convert the field value to a string
	str := fl.Field().String()

	// Define the pattern to match
	pattern := `(?P<Month>\d\d)/(?P<Year>\d\d)`

	// Use regex to match the pattern
	re := regexp.MustCompile(pattern)
	submatch := re.FindStringSubmatch(str)
	if len(submatch) != 3 {
		return false
	}
	// Check month
	monthInt, err := strconv.Atoi(submatch[1])
	if err != nil {
		return false
	}
	if monthInt >= 1 && monthInt <= 12 {
		return true
	}
	return false
}

// init registers the validateExpirationDate function as a custom validator with
// the Validate instance.
func init() {
	Validate.RegisterValidation("exp_date", validateExpirationDate)
}
