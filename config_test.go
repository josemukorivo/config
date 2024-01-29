package config

import (
	"os"
	"testing"
	"time"
)

type Config struct {
	Host string
	Port int    `default:"8080" env:"app_port"`
	User string `env:"config_user" default:"joseph" required:"true"`
}

func TestParse(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_HOST", "localhost")
	os.Setenv("APP_PORT", "8080")

	var cfg Config

	err := Parse("app", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Host != "localhost" {
		t.Fatalf("expected host to be localhost, got %s", cfg.Host)
	}

	if cfg.Port != 8080 {
		t.Fatalf("expected port to be 8080, got %d", cfg.Port)
	}
}

func TestParseInvalidConfig(t *testing.T) {
	tests := []struct {
		description string
		input       any
	}{
		{
			description: "string type config",
			input:       new(string),
		},
		{
			description: "map type config",
			input:       make(map[string]string),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			err := Parse("app", tc.input)
			if err != ErrInvalidConfig {
				t.Fatalf("expected ErrInvalidConfig, got %v", err)
			}
		})
	}
}

func TestParseValueConfig(t *testing.T) {
	var cfg Config

	// cfg is not a pointer type.
	err := Parse("app", cfg)
	if err != ErrInvalidConfig {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestFieldIntError(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_PORT", "not_a_number")
	os.Setenv("APP_HOST", "localhost")

	var cfg Config
	if err := Parse("app", &cfg); err == nil {
		t.Fatal("expected error, got nil")
	} else {
		if v, ok := err.(*FieldError); !ok {
			t.Fatalf("expected FieldError, got %v", v)
		}
	}
}

func TestAlternateEnvName(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_CONFIG_USER", "root")

	var cfg Config
	if err := Parse("app", &cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.User != "root" {
		t.Fatalf("expected user to be root, got %s", cfg.User)
	}
}

func TestAlternateWithoutPrefix(t *testing.T) {
	os.Clearenv()
	os.Setenv("MY_PORT", "9000")

	spec := struct {
		Port int `env:"my_port"`
	}{}

	if err := Parse("app", &spec); err != nil {
		t.Fatal(err)
	}

	if spec.Port != 9000 {
		t.Fatalf("expected port to be 9000, got %d", spec.Port)
	}
}

func TestDefault(t *testing.T) {
	os.Clearenv()
	spec := struct {
		User string `default:"joseph"`
	}{}
	if err := Parse("app", &spec); err != nil {
		t.Fatal(err)
	}

	if spec.User != "joseph" {
		t.Fatalf("expected user to be joseph, got %s", spec.User)
	}
}

func TestRequired(t *testing.T) {

	spec := struct {
		Host string `required:"true"`
	}{}

	if err := Parse("app", &spec); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseDuration(t *testing.T) {

	spec := struct {
		Timeout time.Duration `default:"2s"`
	}{}

	if err := Parse("app", &spec); err != nil {
		t.Fatal(err)
	}
}
func TestNestedStruct(t *testing.T) {

	spec := struct {
		Web struct {
			Host string
		}
		DB struct {
			Port int
		}
	}{}

	os.Clearenv()
	os.Setenv("APP_WEB_HOST", "localhost")
	os.Setenv("APP_DB_PORT", "5432")

	if err := Parse("app", &spec); err != nil {
		t.Fatal(err)
	}

	if spec.Web.Host != "localhost" {
		t.Fatalf("expected web host to be localhost, got %s", spec.Web.Host)
	}

	if spec.DB.Port != 5432 {
		t.Fatalf("expected db port to be 5432, got %d", spec.DB.Port)
	}

}

func TestMustParse(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_HOST", "localhost")
	os.Setenv("APP_PORT", "8080")

	var cfg Config

	MustParse("app", &cfg)

	if cfg.Host != "localhost" {
		t.Fatalf("expected host to be localhost, got %s", cfg.Host)
	}

	if cfg.Port != 8080 {
		t.Fatalf("expected port to be 8080, got %d", cfg.Port)
	}

	defer func() {
		if err := recover(); err != nil {
			return
		}
		t.Fatal("expected panic, got nil")
	}()

	m := make(map[string]string)
	MustParse("app", m)

}
