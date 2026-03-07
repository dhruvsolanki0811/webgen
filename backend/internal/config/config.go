package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	MongoURI             string
	MongoDB              string
	JWTSecret            string
	ClaudeAPIKey         string
	GitHubAppID          string
	GitHubPrivateKeyPath string
	GitHubOrgName        string
	VercelToken          string
	FrontendURL          string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:                 envOrDefault("PORT", "8080"),
		MongoURI:             os.Getenv("MONGO_URI"),
		MongoDB:              envOrDefault("MONGO_DB", "webgen"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		ClaudeAPIKey:         os.Getenv("CLAUDE_API_KEY"),
		GitHubAppID:          os.Getenv("GITHUB_APP_ID"),
		GitHubPrivateKeyPath: os.Getenv("GITHUB_PRIVATE_KEY_PATH"),
		GitHubOrgName:        os.Getenv("GITHUB_ORG_NAME"),
		VercelToken:          os.Getenv("VERCEL_TOKEN"),
		FrontendURL:          envOrDefault("FRONTEND_URL", "http://localhost:3000"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	required := map[string]string{
		"MONGO_URI":  c.MongoURI,
		"JWT_SECRET": c.JWTSecret,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}

	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
