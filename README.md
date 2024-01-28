# Efficient Configuration and Environment Variable Management for Golang

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/josemukorivo/config)](https://goreportcard.com/report/github.com/josemukorivo/config)
[![GoDoc](https://godoc.org/github.com/josemukorivo/config?status.svg)](https://godoc.org/github.com/josemukorivo/config)

Config is a lightweight and flexible Golang package designed to simplify configuration management in your applications. It seamlessly handles configuration files and environment variables, providing a unified interface for accessing and managing your application's settings.

## Features

- Simple and easy to use
- Supports environment variables
- Supports configuration files
- Supports nested configuration
- Supports default values
- Supports validation
- Supports custom configuration sources

## Installation

```bash
go get -u github.com/josemukorivo/config
```

## Usage

#### Configuration File

`main.go`
```go
package main

import (
	"fmt"
	"log"

	"github.com/josemukorivo/config"
)

type Config struct {
	Host string
	Port int
}

func main() {
	var cfg Config

	// Load configuration from environment variables file
	if err := config.Parse("app", &cfg, ".env.local"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
}
```

`.env.local`
```bash
APP_HOST=localhost
APP_PORT=8080
```

### Default Values

```go
package main

import (
	"fmt"
	"log"

	"github.com/josemukorivo/config"
)

type Config struct {
	Host string `default:"localhost"`
	Port int    `default:"8080"`
}

func main() {
	var cfg Config

	// Looks for a file named .env by default
	if err := config.Parse("app", &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
}
```

### Nested Configuration

```go

package main

import (
	"fmt"
	"log"

	"github.com/josemukorivo/config"
)

type Config struct {
	Host string
	Port int
	DB   struct {
		Host     string
		Port     int
		Username string
		Password string
	}
}

func main() {
	var cfg Config

	if err := config.Parse("app", &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
	fmt.Println(cfg.DB.Host)
	fmt.Println(cfg.DB.Port)
	fmt.Println(cfg.DB.Username)
	fmt.Println(cfg.DB.Password)
}
```

### Validation

```go
package main

import (
	"fmt"
	"log"

	"github.com/josemukorivo/config"
)

type Config struct {
	Host string `required:"true"`
	Port int    `required:"true"`
}


func main() {
	var cfg Config

	if err := config.Parse("app", &cfg); err != nil {
		log.Fatal(err) // Missing required configuration parameters if Host or Port are not provided in the environment
	}

	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
