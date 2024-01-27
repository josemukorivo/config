package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	env "github.com/joho/godotenv"
)

// ErrInvalidConfig is returned when the config is not a pointer to struct.
var ErrInvalidConfig = errors.New("config: invalid config must be a pointer to struct")
var ErrRequiredField = errors.New("config: required field missing value")

// Parse parses the config, the config must be a pointer to struct and the struct can contain nested structs.
func Parse(prefix string, cfg any, envFiles ...string) error {
	// Load the .env file if it exists.
	env.Load(envFiles...)
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
				return &FieldError{
					fieldName:  fieldName,
					fieldType:  f.Kind().String(),
					fieldValue: value,
					fieldErr:   err,
				}
			}
		}
	}
	return nil
}

// MustParse parses the config and panics if an error occurs.
func MustParse(prefix string, cfg any) {
	if err := Parse(prefix, cfg); err != nil {
		panic(err)
	}
}
