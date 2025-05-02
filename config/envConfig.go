package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// EnvConfig は環境変数の設定を管理する構造体です
type EnvConfig struct {
	// 環境変数のキャッシュ（必要に応じて拡張可能）
	loadedFiles map[string]bool
}

// NewEnvConfig は新しいEnvConfigインスタンスを作成します
func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		loadedFiles: make(map[string]bool),
	}
}

// LoadEnvFile は指定された環境変数ファイルを読み込みます
// ファイルが存在しない場合はエラーを返します
func LoadEnvFile(filename string) error {
	// ファイルの絶対パスを取得
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// ファイルの存在確認
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("env file not found: %s", absPath)
	}

	// 環境変数ファイルの読み込み
	err = godotenv.Load(absPath)
	if err != nil {
		return fmt.Errorf("failed to load env file %s: %w", absPath, err)
	}

	log.Printf("Loaded environment variables from %s", absPath)
	return nil
}

// LoadEnvFileOrDefault は指定された環境変数ファイルを読み込みます
// ファイルが存在しない場合はエラーを返さずに続行します
func LoadEnvFileOrDefault(filename string) {
	err := LoadEnvFile(filename)
	if err != nil {
		log.Printf("WARN: %v", err)
	}
}

// LoadEnvFiles は複数の環境変数ファイルを読み込みます
// 引数に渡された順番で読み込むため、後から読み込まれた値が優先されます
func LoadEnvFiles(filenames ...string) error {
	for _, filename := range filenames {
		if err := LoadEnvFile(filename); err != nil {
			return err
		}
	}
	return nil
}

// GetRequiredEnv は指定された環境変数を取得します
// 環境変数が設定されていない場合はパニックを発生させます
func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Required environment variable not set: %s", key)
	}
	return value
}

// GetEnvWithDefault は指定された環境変数を取得します
// 環境変数が設定されていない場合はデフォルト値を返します
func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
