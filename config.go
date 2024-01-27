package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidConfig is returned when the config is not a pointer to struct.
var ErrInvalidConfig = errors.New("config: invalid config must be a pointer to struct")
var ErrRequiredField = errors.New("config: required field missing value")

// FieldError is returned when a field cannot be parsed.
type FieldError struct {
	Key   string
	Name  string
	Type  string
	Value string
	Err   error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("config: field %s of type %s has invalid value %s", e.Name, e.Type, e.Value)
}

func Parse(prefix string, cfg any) error {
	if reflect.TypeOf(cfg).Kind() != reflect.Ptr {
		return ErrInvalidConfig
	}
	v := reflect.ValueOf(cfg).Elem()
	if v.Kind() != reflect.Struct {
		return ErrInvalidConfig
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			newPrefix := fmt.Sprintf("%s_%s", prefix, t.Field(i).Name)
			err := Parse(newPrefix, f.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}
		if f.CanSet() {
			var fieldName string

			customVariable := t.Field(i).Tag.Get("env")
			if customVariable != "" {
				fieldName = customVariable
			} else {
				fieldName = t.Field(i).Name
			}
			key := strings.ToUpper(fmt.Sprintf("%s_%s", prefix, fieldName))
			value := os.Getenv(key)
			// If you can't find the value, try to find the value without the prefix.
			if value == "" && customVariable != "" {
				key := strings.ToUpper(fieldName)
				value = os.Getenv(key)
			}

			def := t.Field(i).Tag.Get("default")
			if value == "" && def != "" {
				value = def
			}

			req := t.Field(i).Tag.Get("required")

			if value == "" {
				if req == "true" {
					return ErrRequiredField
				}
				continue
			}

			err := parseField(value, f)
			if err != nil {
				return &FieldError{Key: key, Name: fieldName, Type: f.Kind().String(), Value: value, Err: err}
			}

		}
	}
	return nil
}

func MustParse(prefix string, cfg any) {
	if err := Parse(prefix, cfg); err != nil {
		panic(err)
	}
}

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
