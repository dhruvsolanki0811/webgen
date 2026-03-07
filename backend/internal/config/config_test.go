package config

import (
	"os"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("JWT_SECRET", "test-secret-32-chars-long-enough")
	defer os.Unsetenv("MONGO_URI")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.MongoURI != "mongodb://localhost:27017" {
		t.Errorf("expected MongoURI = mongodb://localhost:27017, got %s", cfg.MongoURI)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default Port = 8080, got %s", cfg.Port)
	}
}

func TestLoad_MissingMongoURI(t *testing.T) {
	os.Unsetenv("MONGO_URI")
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing MONGO_URI")
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("MONGO_URI")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing JWT_SECRET")
	}
}
