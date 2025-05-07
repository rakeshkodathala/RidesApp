package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	GoogleMapsKey  string
	KafkaBrokers   string
	RedisURL       string
	WebSocketPort  string
	ServerPort     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "ridesapp"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		GoogleMapsKey:  getEnv("GOOGLE_MAPS_API_KEY", ""),
		KafkaBrokers:   getEnv("KAFKA_BROKERS", "localhost:9092"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
		WebSocketPort:  getEnv("WS_PORT", "8081"),
		ServerPort:     getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
