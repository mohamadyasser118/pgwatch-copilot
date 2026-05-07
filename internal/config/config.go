package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	SinkConn    string 
	DBConn      string 
	SysID       string 
	LLMProvider string 
	LLMModel    string 
	APIKey      string 
	Window      string 
}


func Load(envFile string) (*Config, error) {
	if envFile != "" {
		_ = godotenv.Load(envFile)
	} else {
		_ = godotenv.Load(".env")  // default
	}

	cfg := &Config{
		SinkConn:    os.Getenv("SINK_CONN"),
		DBConn:      os.Getenv("DB_CONN"),
		SysID:       os.Getenv("SYS_ID"),
		LLMProvider: os.Getenv("LLM_PROVIDER"),
		LLMModel:    os.Getenv("LLM_MODEL"),
		APIKey:      os.Getenv("LLM_API_KEY"),
		Window:      os.Getenv("WINDOW"),
	}

	// Set defaults
	if cfg.LLMProvider == "" {
		cfg.LLMProvider = "openai"
	}
	if cfg.LLMModel == "" {
		cfg.LLMModel = "gpt-4o-mini"
	}
	if cfg.Window == "" {
		cfg.Window = "5m"
	}

	// Validate required fields
	if cfg.SinkConn == "" {
		return nil, fmt.Errorf("SINK_CONN is required")
	}
	if cfg.SysID == "" {
		return nil, fmt.Errorf("SYS_ID is required — run: SELECT DISTINCT sys_id FROM metrics.pgwatch2_pgwatch2")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY is required")
	}

	return cfg, nil
}