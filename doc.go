/*
Package config provides a simple way to parse environment variables into a struct.
It supports nested structs and custom tags.

Example:

Here is an example of how to use this package:

	package main

	import (
		"fmt"
		"os"
		"github.com/josemukorivo/config"
		)

	type Config struct {
		Host string
		Port int    `config:"default=8080,required=true,env=app_port"`
		User string `env:"config_user" default:"joseph" required:"true"`
		}

	func main() {
		os.Clearenv()
		os.Setenv("APP_HOST", "localhost")
		os.Setenv("APP_PORT", "8080")

		var cfg Config

		err := config.Parse("app", &cfg)
		if err != nil {
			panic(err)
		}

		fmt.Println(cfg.Host)
		fmt.Println(cfg.Port)
		fmt.Println(cfg.User)
		}

*/

package config
