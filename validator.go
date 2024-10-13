package gup

import (
	"encoding/json"
	"reflect"
	"strings"
)

type validator[T interface{}] struct {
	fields map[string]*RuleList
	errors map[string][]string
	values map[string]interface{}
}

func NewValidator[T interface{}](fields map[string]*RuleList) *validator[T] {
	return &validator[T]{fields: fields, errors: map[string][]string{}, values: map[string]interface{}{}}
}

func (v *validator[T]) GetErrors() []byte {
	jsonb, err := json.Marshal(v.errors)

	if err != nil {
		return []byte("{}")
	}

	return jsonb
}

func (v *validator[T]) GetValues() map[string]interface{} {
	return v.values
}

func (v *validator[T]) Validate(object *T) bool {
	var hasErrors bool

	for key, rules := range v.fields {

		value := reflect.ValueOf(object).Elem().FieldByName(strings.Title(key))

		parsedValue, errors := rules.Execute(key, value.String())

		if len(errors) > 0 {
			hasErrors = true
			v.errors[key] = errors
		} else {
			v.values[key] = parsedValue
			value.Set(reflect.ValueOf(parsedValue))
		}
	}

	return !hasErrors
}
