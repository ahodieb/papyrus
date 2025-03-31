package cmd

import "os"

type Config struct {
	File       string
	EventsFile string
	Editor     string
}

const (
	DefaultFile       = "~/notes.txt"
	DefaultEventsFile = "~/events.txt"
	DefaultEditor     = "vim"
)

func LoadConfig() Config {
	return Config{
		File:       GetEnvOrDefault("NOTES", DefaultFile),
		EventsFile: GetEnvOrDefault("EVENTS", DefaultEventsFile),
		Editor:     GetEnvOrDefault("EDITOR", DefaultEditor),
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
