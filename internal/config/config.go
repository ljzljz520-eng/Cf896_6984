package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Address      string
	DatabasePath string
	SeedFixtures bool
}

func FromEnv() Config {
	c := Config{Address: ":8080", DatabasePath: "league.db", SeedFixtures: true}
	if value := os.Getenv("LEAGUE_ADDR"); value != "" {
		c.Address = value
	}
	if value := os.Getenv("LEAGUE_DB"); value != "" {
		c.DatabasePath = value
	}
	if value := os.Getenv("LEAGUE_SEED"); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			c.SeedFixtures = parsed
		}
	}
	return c
}

func (c Config) Validate() error {
	if c.Address == "" || c.DatabasePath == "" {
		return fmt.Errorf("address and database path are required")
	}
	return nil
}
