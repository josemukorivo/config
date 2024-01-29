package config

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// FieldError is returned when a field cannot be parsed.
type FieldError struct {
	fieldName  string
	fieldType  string
	fieldValue string
	fieldErr   error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("config: error assigning to field %s: converting '%s' to type %s. details: %s",
		e.fieldName, e.fieldValue, e.fieldType, e.fieldErr,
	)
}

// Setter is the interface that wraps the Decode method. A type that implements the Setter interface
// can set it's value from a string value passed to the Set method.
type Setter interface {
	Set(value string) error
}

// parseField parses a string value into a field.
func parseField(value string, field reflect.Value) error {
	t := field.Type()

	// If the field implements the Setter interface, use it to set it's value.
	// Otherwise, use the default parser. This allows for custom types to be used.
	if setter := extractSetter(field); setter != nil {
		return setter.Set(value)
	}

	switch t.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && field.Type().PkgPath() == "time" && field.Type().Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, field.Type().Bits())
		}
		if err != nil {
			return err
		}
		field.SetInt(val)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	}
	return nil
}

// extractInterface extracts the interface from a field. It checks if the field implements the interface
// and if not, it checks if the field's address implements the interface. If the interface is found,
// the ok parameter is set to true. Otherwise, it is set to false.
func extractInterface(field reflect.Value, fn func(any, *bool)) {
	var ok bool
	if field.CanInterface() {
		fn(field.Interface(), &ok)
	}
	if !ok && field.CanAddr() {
		fn(field.Addr().Interface(), &ok)
	}
}

// extractSetter returns a Setter if the field implements the Setter interface.
// Otherwise, it returns nil.
func extractSetter(field reflect.Value) Setter {
	var s Setter
	// Check if the field implements the Setter interface.
	extractInterface(field, func(v any, ok *bool) {
		s, *ok = v.(Setter)
	})
	return s
}
