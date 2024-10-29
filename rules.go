package gup

import (
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
)

type ValidationRule func(key string, value interface{}) (interface{}, error)

type comparationFunction func(val float64, err error) (interface{}, error)

type RuleList struct {
	Values *map[string]interface{}
	Rules  []ValidationRule
}

func R() *RuleList {
	return &RuleList{}
}

func (ruleList *RuleList) UseValues(values *map[string]interface{}) {
	ruleList.Values = values
}

func (ruleList *RuleList) Execute(key string, value interface{}) (interface{}, []string) {
	var parsedValue interface{}
	var errors []string

	for _, executeRule := range ruleList.Rules {
		validated, err := executeRule(key, value)

		if err != nil {
			errors = append(errors, err.Error())
		}

		parsedValue = validated
	}

	return parsedValue, errors
}

func (ruleList *RuleList) Add(rule ValidationRule) *RuleList {
	ruleList.Rules = append(ruleList.Rules, rule)
	return ruleList
}

func (ruleList *RuleList) Required() *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		validationError := errors.New("this field is required")

		switch val := value.(type) {
		case string:
			if strings.Trim(val, "") == "" {
				return value, validationError
			}
		case nil:
			return value, validationError
		}
		return value, nil
	})

	return ruleList
}

func (ruleList *RuleList) Email() *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		validationError := errors.New("must be a valid email address")

		switch val := value.(type) {
		case string:
			_, err := mail.ParseAddress(val)

			if err != nil {
				return value, validationError
			}
		default:
			return value, validationError
		}

		return value, nil
	})

	return ruleList
}

func (ruleList *RuleList) Number() *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		validationError := errors.New("must be a valid number")

		switch val := value.(type) {
		case string:
			_, err := strconv.ParseFloat(val, 64)

			if err != nil {
				return value, validationError
			}

			return value, nil
		case nil:
			return value, validationError
		}

		return value, nil
	})

	return ruleList
}

func (ruleList *RuleList) MinValue(minValue float64) *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		validationError := fmt.Errorf("minimum value must be %.2f", minValue)

		compare := func(val float64, err error) (interface{}, error) {
			if err != nil {
				return val, fmt.Errorf("invalid value informed: \"%s\"", value)
			}

			if val < minValue {
				return value, validationError
			}

			return val, nil
		}

		return compareValue(value, compare)
	})

	return ruleList
}

func (ruleList *RuleList) MaxValue(maxValue float64) *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		validationError := fmt.Errorf("maximum value must be %.2f", maxValue)

		compare := func(val float64, err error) (interface{}, error) {
			if err != nil {
				return val, fmt.Errorf("invalid value informed: \"%s\"", value)
			}

			if val > maxValue {
				return value, validationError
			}

			return val, nil
		}

		return compareValue(value, compare)
	})
	return ruleList
}

func compareValue(value interface{}, compare comparationFunction) (interface{}, error) {
	switch val := value.(type) {
	case string:
		parsedValue, err := strconv.ParseFloat(val, 64)
		return compare(float64(parsedValue), err)
	case int:
		return compare(float64(val), nil)
	case int8:
		return compare(float64(val), nil)
	case int16:
		return compare(float64(val), nil)
	case int32:
		return compare(float64(val), nil)
	case int64:
		return compare(float64(val), nil)
	case uint:
		return compare(float64(val), nil)
	case uint8:
		return compare(float64(val), nil)
	case uint16:
		return compare(float64(val), nil)
	case uint32:
		return compare(float64(val), nil)
	case uint64:
		return compare(float64(val), nil)
	case float32:
		return compare(float64(val), nil)
	case float64:
		return compare(val, nil)
	default:
		return value, fmt.Errorf("unknown type %T", val)
	}
}

func (ruleList *RuleList) Trim() *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		val, ok := value.(string)

		if ok {
			return strings.Trim(val, " "), nil
		}

		return value, nil
	})

	return ruleList
}

func (ruleList *RuleList) MinLength(minLenght int) *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		val, ok := value.(string)

		if ok && len(val) < minLenght {
			return value, fmt.Errorf("minimum length must be %d", minLenght)
		}

		return value, nil
	})

	return ruleList
}

func (ruleList *RuleList) MaxLength(maxLenght int) *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		val, ok := value.(string)
		if ok && len(val) > maxLenght {
			return value, fmt.Errorf("maximum length must be %d", maxLenght)
		}

		return value, nil
	})

	return ruleList
}

func (ruleList *RuleList) SameAs(field string) *RuleList {
	ruleList.Rules = append(ruleList.Rules, func(key string, value interface{}) (interface{}, error) {
		validationError := fmt.Errorf("must be the same as %s", field)
		values := *ruleList.Values

		if !reflect.DeepEqual(values[field], value) {
			return value, validationError
		}

		return value, nil
	})

	return ruleList
}
