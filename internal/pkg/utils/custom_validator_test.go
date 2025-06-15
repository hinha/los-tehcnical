package utils

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
)

func TestValidateLoanState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		want  bool
	}{
		{"valid proposed state", "proposed", true},
		{"valid approved state", "approved", true},
		{"valid invested state", "invested", true},
		{"valid disbursed state", "disbursed", true},
		{"invalid state", "invalid", false},
		{"empty state", "", false},
	}

	validate := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.RegisterValidation("loan_state", ValidateLoanState)
			if err != nil {
				t.Fatalf("failed to register validation: %v", err)
			}

			got := ValidateLoanState(getMockFieldLevel(tt.state))
			if got != tt.want {
				t.Errorf("ValidateLoanState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getMockFieldLevel(value string) validator.FieldLevel {
	return mockFieldLevel{value: value}
}

type mockFieldLevel struct {
	value string
}

func (m mockFieldLevel) Top() reflect.Value {
	return reflect.Value{}
}

func (m mockFieldLevel) Parent() reflect.Value {
	return reflect.Value{}
}

func (m mockFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(m.value)
}

func (m mockFieldLevel) FieldName() string {
	return ""
}

func (m mockFieldLevel) StructFieldName() string {
	return ""
}

func (m mockFieldLevel) Param() string {
	return ""
}

func (m mockFieldLevel) GetTag() string {
	return ""
}

func (m mockFieldLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func (m mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}
