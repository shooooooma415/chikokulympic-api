package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	loadedFiles map[string]bool
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		loadedFiles: make(map[string]bool),
	}
}

func LoadEnvFile(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("env file not found: %s", absPath)
	}

	err = godotenv.Load(absPath)
	if err != nil {
		return fmt.Errorf("failed to load env file %s: %w", absPath, err)
	}

	log.Printf("Loaded environment variables from %s", absPath)
	return nil
}

func LoadEnvFileOrDefault(filename string) {
	err := LoadEnvFile(filename)
	if err != nil {
		log.Printf("WARN: %v", err)
	}
}

func LoadEnvFiles(filenames ...string) error {
	for _, filename := range filenames {
		if err := LoadEnvFile(filename); err != nil {
			return err
		}
	}
	return nil
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Required environment variable not set: %s", key)
	}
	return value
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
