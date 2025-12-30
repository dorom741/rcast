package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultUUID     = "uuid:0199ffd9-6856-74cc-a2f2-4c74af0161b1"
	DefaultPort     = 8200
	DefaultUUIDPath = ".local/rcast/dmr_uuid.txt"
)

type Config struct {
	UUIDPath               string
	AllowSessionPreempt    bool
	LinkSystemOutputVolume bool
	HTTPPort               int
	IINAFullscreen         bool
	PlayerType             string // Player type: iina | aria2
	Aria2RPCURL            string // Aria2 RPC URL
	Aria2RPCPassword       string // Aria2 RPC password
	Aria2DownloadPath      string // Aria2 download path
}

// loadEnvFile loads environment variables from .env file if it exists
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return // Ignore if .env file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// If the value is not already set in the environment, set it
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func Load() Config {
	// Load .env file if it exists
	loadEnvFile()

	cfg := Config{
		UUIDPath:               envVar("DMR_UUID_PATH", os.Getenv("HOME")+"/"+DefaultUUIDPath),
		AllowSessionPreempt:    envVar("DMR_ALLOW_PREEMPT", true),
		LinkSystemOutputVolume: envVar("DMR_LINK_SYSTEM_VOLUME", false),
		HTTPPort:               envVar("DMR_HTTP_PORT", DefaultPort),
		IINAFullscreen:         envVar("DMR_IINA_FULLSCREEN", false),
		
		PlayerType:             envVar("DMR_PLAYER_TYPE", "iina"), // Default to iina player
		Aria2RPCURL:            envVar("DMR_ARIA2_RPC_URL", "http://localhost:6800/jsonrpc"),
		Aria2RPCPassword:       envVar("DMR_ARIA2_RPC_PASSWORD", ""),
		Aria2DownloadPath:      envVar("DMR_ARIA2_DOWNLOAD_PATH", ""), // Default to empty string
	}

	// Validate configuration
	cfg.validate()

	return cfg
}

func envVar[T ~string | ~bool | ~int](key string, def T) T {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	switch any(def).(type) {
	case string:
		return any(v).(T)
	case bool:
		if b, err := strconv.ParseBool(v); err == nil {
			return any(b).(T)
		}
	case int:
		if i, err := strconv.Atoi(v); err == nil {
			return any(i).(T)
		}
	case int64:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return any(i).(T)
		}
	case float64:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return any(f).(T)
		}
	}
	return def
}

// validate performs validation on configuration values
func (c *Config) validate() {
	// Validate HTTP port range
	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		c.HTTPPort = DefaultPort
	}

	// Validate player type
	if c.PlayerType != "iina" && c.PlayerType != "aria2" {
		c.PlayerType = "iina" // Default to iina if invalid
	}

	// Ensure UUID path directory exists
	if c.UUIDPath != "" {
		// Extract directory from path
		if idx := strings.LastIndex(c.UUIDPath, "/"); idx > 0 {
			dir := c.UUIDPath[:idx]
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				// Create directory if it doesn't exist
				os.MkdirAll(dir, 0755)
			}
		}
	}
}
