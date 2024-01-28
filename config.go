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
var (
	ErrInvalidConfig = errors.New("config: invalid config must be a pointer to struct")
)

// Parse parses the config, the config must be a pointer to struct and the struct can contain nested structs.
// The prefix is used to prefix the environment variables. For example, if the prefix is "app" and the struct
// contains a field named "Host", the environment variable will be "APP_HOST". If the struct contains a nested
// struct, the prefix will be the original prefix plus the nested struct name. For example, if the prefix is "app"
// and the nested struct is named "DB", the environment variable will be "APP_DB_HOST". Parse take an optional
// list of .env files to load. If the .env file exists, it will be loaded before parsing the config. By default,
// Parse will look for a .env file and parse it.
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
					return errors.New("config: required field missing value")
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
// See Parse for more information. MustParse is a wrapper around Parse.
func MustParse(prefix string, cfg any, envFiles ...string) {
	if err := Parse(prefix, cfg, envFiles...); err != nil {
		panic(err)
	}
}
