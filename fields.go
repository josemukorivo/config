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

// parseField parses a string value into a field.
func parseField(value string, field reflect.Value) error {
	switch field.Kind() {
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
