package gup

import (
	"reflect"
	"testing"
)

func TestValidationWithInvalidData(t *testing.T) {
	type UserRegistration struct {
		Name               string
		Email              string
		Password           string
		PasswordValidation string
	}

	myRegistration := UserRegistration{
		Name:               "A",
		Email:              "alexandre",
		Password:           "easyp",
		PasswordValidation: "hardpassword",
	}

	registerValidator := NewValidator[UserRegistration](map[string]*RuleList{
		"name":               R().Required().Trim().MinLength(2).MaxLength(100),
		"email":              R().Required().Email(),
		"password":           R().Required().MinLength(8),
		"passwordValidation": R().Required().SameAs("password"),
	})

	if !registerValidator.Validate(&myRegistration) {
		errors := registerValidator.GetErrors()

		if !reflect.DeepEqual(errors["name"][0], "minimum length must be 2") {
			t.Error("invalid validation message for name")
		}

		if !reflect.DeepEqual(errors["email"][0], "must be a valid email address") {
			t.Error("invalid validation message for email")
		}

		if !reflect.DeepEqual(errors["password"][0], "minimum length must be 8") {
			t.Error("invalid validation message for password")
		}

		if !reflect.DeepEqual(errors["passwordValidation"][0], "must be the same as password") {
			t.Error("invalid validation message for passwordValidation")
		}

	} else {
		t.Error("validation is expected to fail")
	}
}

func TestValidationWithValidData(t *testing.T) {
	type UserRegistration struct {
		Name               string
		Email              string
		Password           string
		PasswordValidation string
	}

	myRegistration := UserRegistration{
		Name:               "Alexandre",
		Email:              "alexandre@gmail.com",
		Password:           "easypassword",
		PasswordValidation: "easypassword",
	}

	registerValidator := NewValidator[UserRegistration](map[string]*RuleList{
		"name":               R().Required().Trim().MinLength(2).MaxLength(100),
		"email":              R().Required().Email(),
		"password":           R().Required().MinLength(8),
		"passwordValidation": R().Required().SameAs("password"),
	})

	if registerValidator.Validate(&myRegistration) == false {
		t.Error("validation is expected to pass")
	}
}
