# config: Efficient Configuration and Environment Variable Management for Golang

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/josemukorivo/config)](https://goreportcard.com/report/github.com/josemukorivo/config)
[![GoDoc](https://godoc.org/github.com/josemukorivo/config?status.svg)](https://godoc.org/github.com/josemukorivo/config)

Config is a lightweight and flexible Golang package designed to simplify configuration management in your applications. It seamlessly handles configuration files and environment variables, providing a unified interface for accessing and managing your application's settings.

## Features

- **Environment Variables**: Automatically bind environment variables to configuration fields, simplifying deployment and dynamic configuration.
- **Nested Configuration**: Support for nested and structured configuration data for better organization and readability.
- **Default Values**: Define default values for configuration parameters, ensuring your application runs smoothly even when specific settings are not provided.
- **Validation**: Validate configuration values against predefined rules, catching errors early in the application lifecycle.

## Installation

```bash
go get -u github.com/josemukorivo/config

