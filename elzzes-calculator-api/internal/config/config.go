// Package config loads runtime configuration from the environment.
package config

import (
	"fmt"
	"strconv"
)

type Config struct {
	Port int
}

func (c Config) Addr() string {
	return ":" + strconv.Itoa(c.Port)
}

func Load(getenv func(string) string) (Config, error) {
	raw := getenv("PORT")
	if raw == "" {
		return Config{Port: 8080}, nil
	}

	port, err := strconv.Atoi(raw)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("invalid PORT %q: must be an integer between 1 and 65535", raw)
	}
	return Config{Port: port}, nil
}
