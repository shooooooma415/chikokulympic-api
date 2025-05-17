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
	cache       map[string]string
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		loadedFiles: make(map[string]bool),
		cache:       make(map[string]string),
	}
}

func (e *EnvConfig) LoadEnvFile(filename string) error {
	if e.loadedFiles[filename] {
		return nil
	}

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

	e.loadedFiles[filename] = true
	e.ClearCache() 
	log.Printf("Loaded environment variables from %s", absPath)
	return nil
}

func (e *EnvConfig) LoadEnvFileOrDefault(filename string) {
	err := e.LoadEnvFile(filename)
	if err != nil {
		log.Printf("WARN: %v", err)
	}
}

func (e *EnvConfig) LoadEnvFiles(filenames ...string) error {
	for _, filename := range filenames {
		if err := e.LoadEnvFile(filename); err != nil {
			return err
		}
	}
	return nil
}

func (e *EnvConfig) TryLoadEnvFile(filename string) (bool, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %w", err)
	}

	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return false, nil
	}

	err = e.LoadEnvFile(filename)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (e *EnvConfig) LoadFromFileOrEnv(filename string) {
	loaded, err := e.TryLoadEnvFile(filename)
	if err != nil {
		log.Printf("WARN: Failed to load env file %s: %v", filename, err)
	}

	if loaded {
		log.Printf("Using environment variables from file: %s", filename)
	} else {
		log.Printf("Environment file %s not found, using system environment variables", filename)
	}
}

func (e *EnvConfig) Get(key string) string {
	if value, exists := e.cache[key]; exists {
		return value
	}

	value := os.Getenv(key)
	e.cache[key] = value
	return value
}

func (e *EnvConfig) GetRequired(key string) string {
	value := e.Get(key)
	if value == "" {
		log.Panicf("Required environment variable not set: %s", key)
	}
	return value
}

func (e *EnvConfig) GetWithDefault(key, defaultValue string) string {
	value := e.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (e *EnvConfig) Set(key, value string) {
	os.Setenv(key, value)
	e.cache[key] = value
}

func (e *EnvConfig) ClearCache() {
	e.cache = make(map[string]string)
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

func LoadFromFileOrEnv(filename string) {
	ec := NewEnvConfig()
	ec.LoadFromFileOrEnv(filename)
}

func TryLoadEnvFile(filename string) (bool, error) {
	ec := NewEnvConfig()
	return ec.TryLoadEnvFile(filename)
}
