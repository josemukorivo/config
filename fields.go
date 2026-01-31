package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FieldError is returned when a field cannot be parsed.
type FieldError struct {
	fieldName  string
	fieldType  string
	fieldValue string
	fieldErr   error
}

// Error returns the error message for the FieldError. It includes the field name, the field value,
// the field type and the error returned by the parser.
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

// Field represents a field in a struct.
type Field struct {
	Name     string
	Field    reflect.Value
	Key      string // The key used to look up the value in the environment.
	EnvKey   string // The environment variable name used when overriding the default key.
	Tags     reflect.StructTag
	Required bool
	Default  string
}

// extractFields extracts the fields from the struct and returns a slice of Fields.
func extractFields(prefix string, cfg any) ([]Field, error) {
	if reflect.TypeOf(cfg).Kind() != reflect.Ptr {
		return nil, ErrInvalidConfig
	}
	v := reflect.ValueOf(cfg).Elem()
	if v.Kind() != reflect.Struct {
		return nil, ErrInvalidConfig
	}
	t := v.Type()

	fields := make([]Field, 0, v.NumField())
	for i := range v.NumField() {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		if f.Kind() == reflect.Struct {
			newPrefix := fmt.Sprintf("%s_%s", prefix, t.Field(i).Name)
			err := Parse(newPrefix, f.Addr().Interface())
			if err != nil {
				return nil, err
			}
			continue
		}

		envKey := strings.ToUpper(t.Field(i).Tag.Get("env"))
		key := t.Field(i).Name

		if envKey != "" {
			key = envKey
		}
		if prefix != "" {
			key = strings.ToUpper(fmt.Sprintf("%s_%s", prefix, key))
		}

		key = strings.ToUpper(key)
		required := isTrue(t.Field(i).Tag.Get("required"))
		def := t.Field(i).Tag.Get("default")

		field := Field{
			Name:     t.Field(i).Name,
			Field:    f,
			Tags:     t.Field(i).Tag,
			Key:      key,
			Required: required,
			Default:  def,
			EnvKey:   envKey,
		}

		fields = append(fields, field)
	}
	return fields, nil
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
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			return fmt.Errorf("map keys must be strings")
		}
		mapValue := strings.TrimSpace(value)
		if len(mapValue) >= 2 && mapValue[0] == '"' && mapValue[len(mapValue)-1] == '"' {
			unquoted, err := strconv.Unquote(mapValue)
			if err == nil {
				mapValue = unquoted
			} else {
				mapValue = mapValue[1 : len(mapValue)-1]
			}
		} else if len(mapValue) >= 2 && mapValue[0] == '\'' && mapValue[len(mapValue)-1] == '\'' {
			mapValue = mapValue[1 : len(mapValue)-1]
		}
		mapPtr := reflect.New(t)
		if err := json.Unmarshal([]byte(mapValue), mapPtr.Interface()); err != nil {
			return err
		}
		field.Set(mapPtr.Elem())
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

func isTrue(value string) bool {
	b, _ := strconv.ParseBool(value)
	return b
}
