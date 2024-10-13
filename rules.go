package gup

import (
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
)

type ValidationRule func(key, value string) (string, error)

type RuleList struct {
	Rules []ValidationRule
}

func R() *RuleList {
	return &RuleList{}
}

func (l *RuleList) Execute(key, value string) (string, []string) {
	var parsedValue string
	var errors []string

	for _, executeRule := range l.Rules {
		validated, err := executeRule(key, value)

		if err != nil {
			errors = append(errors, err.Error())
		}

		parsedValue = validated
	}

	return parsedValue, errors
}

func (l *RuleList) Add(rule ValidationRule) *RuleList {
	l.Rules = append(l.Rules, rule)
	return l
}

func (l *RuleList) Required() *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		if strings.Trim(value, "") == "" {
			return value, errors.New("This field is required.")
		}

		return value, nil
	})

	return l
}

func (l *RuleList) Email() *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		_, err := mail.ParseAddress(value)

		if err != nil {
			return value, errors.New("Must be a valid email address.")
		}

		return value, nil
	})

	return l
}

func (l *RuleList) Number() *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		_, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return value, errors.New("Must be a valid number.")
		}

		return value, nil
	})

	return l
}

func (l *RuleList) MinValue(minValue float64) *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		val, err := strconv.ParseFloat(value, 64)

		if err != nil || val < minValue {
			return value, errors.New(fmt.Sprintf("Min value must be %.2f.", minValue))
		}

		return value, nil
	})

	return l
}

func (l *RuleList) MaxValue(maxValue float64) *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		val, err := strconv.ParseFloat(value, 64)

		if err != nil || val > maxValue {
			return value, errors.New(fmt.Sprintf("Max value must be %.2f.", maxValue))
		}

		return value, nil
	})
	return l
}

func (l *RuleList) Trim() *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		return strings.Trim(value, " "), nil
	})

	return l
}

func (l *RuleList) MinLenght(minLenght int) *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		if len(value) < minLenght {
			return value, errors.New(fmt.Sprintf("Min lenght must be %d.", minLenght))
		}

		return value, nil
	})

	return l
}

func (l *RuleList) MaxLenght(maxLenght int) *RuleList {
	l.Rules = append(l.Rules, func(key, value string) (string, error) {
		if len(value) > maxLenght {
			return value, errors.New(fmt.Sprintf("Max lenght must be %d.", maxLenght))
		}

		return value, nil
	})

	return l
}
