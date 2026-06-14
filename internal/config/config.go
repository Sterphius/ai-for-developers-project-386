package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ListenAddr          string
	DatabaseURL         string
	OwnerID             string
	WindowDays          int
	WorkingDayStartHour int
	WorkingDayEndHour   int
	SlotStepMinutes     int
	Location            *time.Location
}

func Load() Config {
	location := time.UTC
	if value := os.Getenv("TIMEZONE"); value != "" {
		if loaded, err := time.LoadLocation(value); err == nil {
			location = loaded
		}
	}

	return Config{
		ListenAddr:          getEnv("LISTEN_ADDR", ":8080"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/calendar?sslmode=disable"),
		OwnerID:             getEnv("OWNER_ID", "owner-1"),
		WindowDays:          getEnvInt("WINDOW_DAYS", 14),
		WorkingDayStartHour: getEnvInt("WORKING_DAY_START_HOUR", 9),
		WorkingDayEndHour:   getEnvInt("WORKING_DAY_END_HOUR", 18),
		SlotStepMinutes:     getEnvInt("SLOT_STEP_MINUTES", 30),
		Location:            location,
	}
}

func getEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}

func getEnvInt(name string, fallback int) int {
	if value := os.Getenv(name); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}

	return fallback
}
