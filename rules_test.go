package gup

import (
	"strings"
	"testing"
)

func TestTrimString(t *testing.T) {
	rule := R().Trim()

	result, errors := rule.Execute("name", "  Alexandre          ")

	if len(result.(string)) > 9 {
		t.Errorf("Trim: expected string length was 9 but got %d", len(result.(string)))
	}

	if len(errors) > 0 {
		t.Errorf("Trim: expect 0 errors in rule execution but got %d errors", len(errors))
	}
}

func TestRequiredFieldIsSet(t *testing.T) {
	rule := R().Required()

	_, errors := rule.Execute("name", "Alexandre")

	if len(errors) == 1 {
		t.Errorf("expect 0 errors in rule execution but got %d", len(errors))
	}
}

func TestRequiredFieldIsUnset(t *testing.T) {
	rule := R().Required()

	_, errors := rule.Execute("name", nil)

	if len(errors) == 0 {
		t.Errorf("expect 1 errors in rule execution but got %d", len(errors))
	}
}

func TestEmailIsValid(t *testing.T) {
	rule := R().Email()

	_, errors := rule.Execute("email", "alexandrecastro@gmail.com")

	if len(errors) == 1 {
		t.Errorf("Email: expected 0 errors in rule execution but got %d", len(errors))
	}
}

func TestEmailIsInvalid(t *testing.T) {
	rule := R().Email()

	_, errors := rule.Execute("email", 22)

	if len(errors) == 0 {
		t.Errorf("expected 1 error in rule execution but got %d", len(errors))
	}
}

func TestNumberIsValid(t *testing.T) {
	rule := R().Number()

	_, errors := rule.Execute("price", 200)

	if len(errors) == 1 {
		t.Errorf("expected 0 errors in rule execution but got %d", len(errors))
	}
}

func TestNumberIsInvalid(t *testing.T) {
	rule := R().Number()

	_, errors := rule.Execute("price", "122")

	if len(errors) == 1 {
		t.Errorf("expected 1 error in rule execution but got %d", len(errors))
	}
}

func TestMinLengthIsInvalid(t *testing.T) {
	rule := R().MinLength(10)

	_, errors := rule.Execute("name", "Tes")

	if len(errors) == 0 {
		t.Errorf("expected 0 errors in rule execution but got %d", len(errors))
	}
}

func TestMinLengthIsValid(t *testing.T) {
	rule := R().MinLength(10)

	_, errors := rule.Execute("name", "Testing name")

	if len(errors) == 1 {
		t.Errorf("expected 1 error in rule execution but got %d", len(errors))
	}
}

func TestMinValueIsValid(t *testing.T) {
	rule := R().MinValue(10)

	_, errors := rule.Execute("name", "11")

	if len(errors) == 1 {
		t.Errorf("expected 0 errors in rule execution but got %d", len(errors))
	}
}

func TestMinValueIsInvalid(t *testing.T) {
	rule := R().MinValue(10)

	_, errors := rule.Execute("name", 9)

	if strings.Compare(errors[0], "minimum value must be 10.00") != 0 {
		t.Errorf("invalid error message: %s", errors[0])
	}

	if len(errors) == 0 {
		t.Errorf("expected 1 error in rule execution but got %d", len(errors))
	}
}

func TestMaxValueIsValid(t *testing.T) {
	rule := R().MaxValue(10)

	_, errors := rule.Execute("name", "9")

	if len(errors) == 1 {
		t.Errorf("expected 0 errors in rule execution but got %d", len(errors))
	}
}

func TestMaxValueIsInvalid(t *testing.T) {
	rule := R().MaxValue(10)

	_, errors := rule.Execute("name", 11)

	if strings.Compare(errors[0], "maximum value must be 10.00") != 0 {
		t.Errorf("invalid error message: %s", errors[0])
	}

	if len(errors) == 0 {
		t.Errorf("expected 1 error in rule execution but got %d", len(errors))
	}
}

func TestSameAsIsInvalid(t *testing.T) {
	rule := R().SameAs("password")

	fields := map[string]interface{}{
		"password": "1234",
	}

	rule.UseValues(&fields)

	_, errors := rule.Execute("password_validation", "1234")

	t.Log(errors)

	if len(errors) == 1 {
		t.Errorf("expected 0 errors in rule execution but got %d", len(errors))
	}
}
