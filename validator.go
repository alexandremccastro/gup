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

func (v *validator[T]) ToJSON() []byte {
	jsonb, err := json.Marshal(v.errors)

	if err != nil {
		return []byte("{}")
	}

	return jsonb
}

func (validator *validator[T]) GetErrors() map[string][]string {
	return validator.errors
}

func (validator *validator[T]) GetValues() map[string]interface{} {
	return validator.values
}

func (validator *validator[T]) Validate(object *T) bool {
	var hasErrors bool

	for key := range validator.fields {
		validator.values[key] = reflect.ValueOf(object).Elem().FieldByName(strings.Title(key)).Interface()
	}

	for key, rules := range validator.fields {
		rules.UseValues(&validator.values)
		parsedValue, errors := rules.Execute(key, validator.values[key])
		validator.values[key] = parsedValue

		if len(errors) > 0 {
			hasErrors = true
			validator.errors[key] = errors
		}
	}

	return !hasErrors
}
