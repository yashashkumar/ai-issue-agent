package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost           string
	ServerPort           string
	BaseURL              string
	DatabasePath         string
	MaxConcurrentAgents  int
	SpawnerChannelBuffer int
	AgentTimeoutMinutes  int
	DefaultWorkBaseDir   string
	GeminiCliPath        string
	GitHubDefaultBranch  string
	LogLevel             string
	LogFormat            string
}

func Load() *Config {
	_ = godotenv.Load() // Ignore error if .env doesn't exist, rely on system env

	return &Config{
		ServerHost:           getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		BaseURL:              getEnv("BASE_URL", "http://localhost:8080"),
		DatabasePath:         getEnv("DATABASE_PATH", "./data/agentforge.db"),
		MaxConcurrentAgents:  getEnvAsInt("MAX_CONCURRENT_AGENTS", 5),
		SpawnerChannelBuffer: getEnvAsInt("SPAWNER_CHANNEL_BUFFER", 100),
		AgentTimeoutMinutes:  getEnvAsInt("AGENT_TIMEOUT_MINUTES", 30),
		DefaultWorkBaseDir:   getEnv("DEFAULT_WORK_BASE_DIR", "./workspaces"),
		GeminiCliPath:        getEnv("GEMINI_CLI_PATH", "gemini"),
		GitHubDefaultBranch:  getEnv("GITHUB_DEFAULT_BASE_BRANCH", "main"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		LogFormat:            getEnv("LOG_FORMAT", "json"),
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("Warning: parsed invalid int for %s: %v, using default %d", key, err, fallback)
		return fallback
	}
	return value
}
