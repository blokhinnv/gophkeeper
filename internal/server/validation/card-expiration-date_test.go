package validation

import (
	"testing"

	// such an import was taken from the validator source...
	. "github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
)

func TestValidateExpirationDate(t *testing.T) {
	type Arg struct {
		Date string `validate:"exp_date"`
	}
	validate := validator.New()
	err := validate.RegisterValidation("exp_date", func(fl validator.FieldLevel) bool {
		return validateExpirationDate(fl)
	})
	Equal(t, err, nil)

	tests := []struct {
		name string
		arg  Arg
		want bool
	}{
		{
			name: "valid date",
			arg:  Arg{"06/12"},
			want: true,
		},
		{
			name: "invalid date",
			arg:  Arg{"13/25"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = validate.Struct(tt.arg)
			if tt.want {
				Equal(t, err, nil)
			} else {
				NotEqual(t, err, nil)
			}
		})
	}
}
