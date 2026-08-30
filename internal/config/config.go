package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server struct {
		Addr        string `koanf:"addr"`
		BaseURL     string `koanf:"base_url"`
		BehindProxy bool   `koanf:"behind_proxy"`
		AdminPath   string `koanf:"admin_path"`
	} `koanf:"server"`
	Database struct {
		Driver string `koanf:"driver"`
		DSN    string `koanf:"dsn"`
	} `koanf:"database"`
	DataDir string `koanf:"data_dir"`
	Log     struct {
		Level      string `koanf:"level"`
		File       string `koanf:"file"`
		MaxSizeMB  int    `koanf:"max_size_mb"`
		MaxBackups int    `koanf:"max_backups"`
	} `koanf:"log"`
	Upload struct {
		MaxSizeMB      int64 `koanf:"max_size_mb"`
		ChunkSizeMB    int64 `koanf:"chunk_size_mb"`
		ImageMaxSizeMB int64 `koanf:"image_max_size_mb"`
	} `koanf:"upload"`
	RateLimit struct {
		PublicPerMin int `koanf:"public_per_min"`
		LoginPerMin  int `koanf:"login_per_min"`
	} `koanf:"ratelimit"`
	Theme struct {
		Dev bool `koanf:"dev"`
	} `koanf:"theme"`
	Timezone string `koanf:"timezone"`
}

func defaults() Config {
	var c Config
	c.Server.Addr = ":8080"
	c.Server.BaseURL = "http://localhost:8080"
	c.Server.AdminPath = "/admin"
	c.Database.Driver = "sqlite"
	c.Database.DSN = "data/netupdown.db"
	c.DataDir = "data"
	c.Log.Level = "info"
	c.Log.File = "data/logs/app.log"
	c.Log.MaxSizeMB = 50
	c.Log.MaxBackups = 5
	c.Upload.MaxSizeMB = 4096
	c.Upload.ChunkSizeMB = 5
	c.Upload.ImageMaxSizeMB = 10
	c.RateLimit.PublicPerMin = 60
	c.RateLimit.LoginPerMin = 5
	c.Timezone = "Asia/Shanghai"
	return c
}

func Load(path string) (Config, error) {
	c := defaults()
	k := koanf.New(".")
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
				return c, fmt.Errorf("load config %s: %w", path, err)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return c, fmt.Errorf("stat config %s: %w", path, err)
		}
	}
	if err := k.Load(env.Provider("NETUPDOWN_", ".", func(key string) string {
		return strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "NETUPDOWN_"), "__", "."))
	}), nil); err != nil {
		return c, fmt.Errorf("load environment: %w", err)
	}
	if err := k.Unmarshal("", &c); err != nil {
		return c, fmt.Errorf("decode config: %w", err)
	}
	if c.Server.Addr == "" || c.DataDir == "" {
		return c, errors.New("server.addr and data_dir are required")
	}
	c.Server.BaseURL = strings.TrimRight(c.Server.BaseURL, "/")
	if !filepath.IsAbs(c.Database.DSN) && c.Database.Driver == "sqlite" && !strings.HasPrefix(c.Database.DSN, "file:") {
		c.Database.DSN = filepath.Clean(c.Database.DSN)
	}
	return c, nil
}
