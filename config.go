package config

import (
	"errors"
	"fmt"
	"os"

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
	fields, err := extractFields(prefix, cfg)
	if err != nil {
		return err
	}

	for _, field := range fields {
		value, ok := os.LookupEnv(field.EnvKey)
		if !ok {
			value, ok = os.LookupEnv(field.Key)
		}

		def := field.Default
		if def != "" && !ok {
			value = def
		}

		if !ok && field.Required && def == "" {
			key := field.Key
			if field.EnvKey != "" {
				key = field.EnvKey
			}
			return fmt.Errorf("config: required key %s missing value", key)
		}
		err := parseField(value, field.Field)
		if err != nil {
			return &FieldError{
				fieldName:  field.Name,
				fieldType:  field.Field.Type().String(),
				fieldValue: value,
				fieldErr:   err,
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
