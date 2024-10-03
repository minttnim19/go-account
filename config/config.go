package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI             string
	MongoDBName            string
	MongoDBUser            string
	MongoDBPassword        string
	TokenExpireTime        string
	TokenRefreshExpireTime string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from a .env file (optional)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Environment variable not set")
	}

	config := &Config{
		MongoDBURI:             GetEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName:            GetEnv("MONGODB_NAME", "myapp"),
		MongoDBUser:            GetEnv("MONGODB_USER", "myuser"),
		MongoDBPassword:        GetEnv("MONGODB_PASSWORD", "mypwd"),
		TokenExpireTime:        GetEnv("TOKEN_EXPIRE_TIME", "86400"),
		TokenRefreshExpireTime: GetEnv("TOKEN_REFRESH_EXPIRE_TIME", "604800"),
	}

	return config, nil
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
