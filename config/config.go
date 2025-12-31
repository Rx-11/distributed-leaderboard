package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RegionID     string
	HTTPPort     int
	GRPCPort     int
	Peers        []string
	FreshTTL     time.Duration
	StaleTTL     time.Duration
	TotalRegions int
	DataDir      string
}

var cfg *Config

func GetConfig() *Config {
	return cfg
}

func Load() {
	_ = godotenv.Load()

	cfg = &Config{
		RegionID:     getEnv("REGION_ID", "us-east"),
		HTTPPort:     getEnvAsInt("HTTP_PORT", 8080),
		GRPCPort:     getEnvAsInt("GRPC_PORT", 9090),
		Peers:        getEnvAsSlice("PEERS", ""),
		FreshTTL:     getEnvAsDuration("FRESH_TTL", 2*time.Second),
		StaleTTL:     getEnvAsDuration("STALE_TTL", 30*time.Second),
		TotalRegions: getEnvAsInt("TOTAL_REGIONS", 1),
		DataDir:      getEnv("DATA_DIR", "./data"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal string) []string {
	valStr := getEnv(key, defaultVal)
	if valStr == "" {
		return []string{}
	}
	return strings.Split(valStr, ",")
}
