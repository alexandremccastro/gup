package gup

import (
	"encoding/json"
	"reflect"
	"strings"
)

type validation[T interface{}] struct {
	fields map[string]*RuleList
	errors map[string][]string
	values map[string]interface{}
}

func NewValidation[T interface{}](fields map[string]*RuleList) *validation[T] {
	return &validation[T]{fields: fields, errors: map[string][]string{}, values: map[string]interface{}{}}
}

func (v *validation[T]) ToJSON() []byte {
	jsonb, err := json.Marshal(v.errors)

	if err != nil {
		return []byte("{}")
	}

	return jsonb
}

func (validation *validation[T]) GetErrors() map[string][]string {
	return validation.errors
}

func (validation *validation[T]) GetValues() map[string]interface{} {
	return validation.values
}

func (validation *validation[T]) Validate(object *T) bool {
	var hasErrors bool

	for key := range validation.fields {
		validation.values[key] = reflect.ValueOf(object).Elem().FieldByName(strings.Title(key)).Interface()
	}

	for key, rules := range validation.fields {
		rules.UseValues(&validation.values)
		parsedValue, errors := rules.Execute(key, validation.values[key])
		validation.values[key] = parsedValue

		if len(errors) > 0 {
			hasErrors = true
			validation.errors[key] = errors
		}
	}

	return !hasErrors
}
