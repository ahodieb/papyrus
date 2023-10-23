package cmd

import "os"

type Config struct {
	File   string
	Editor string
}

const (
	DefaultFile   = "~/notes.txt"
	DefaultEditor = "vim"
)

func LoadConfig() Config {
	return Config{
		File:   GetEnvOrDefault("NOTES", DefaultFile),
		Editor: GetEnvOrDefault("EDITOR", DefaultEditor),
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
