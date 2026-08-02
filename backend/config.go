package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config maps country ID → its rules. Loaded from the resource file.
type Config map[string]CountryConfig

func loadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	if len(c) == 0 {
		return nil, fmt.Errorf("config %s has no countries", path)
	}
	return c, nil
}
