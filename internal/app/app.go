package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/moli-xia/netupdown/internal/auth"
	"github.com/moli-xia/netupdown/internal/config"
	"github.com/moli-xia/netupdown/internal/database"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/moli-xia/netupdown/internal/pkg/cryptoutil"
	"github.com/moli-xia/netupdown/internal/pkg/validate"
	"github.com/moli-xia/netupdown/internal/repo"
	"github.com/moli-xia/netupdown/internal/server"
	"golang.org/x/term"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Serve(ctx context.Context, cfg config.Config, version string) error {
	if err := setupLogger(cfg); err != nil {
		return err
	}
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	jwtKey, err := cryptoutil.LoadOrCreateKey(cfg.DataDir, "jwt.key", "NETUPDOWN_JWT_KEY")
	if err != nil {
		return err
	}
	masterKey, err := cryptoutil.LoadOrCreateKey(cfg.DataDir, "master.key", "NETUPDOWN_MASTER_KEY")
	if err != nil {
		return err
	}
	storageSealer, err := cryptoutil.NewSealer(masterKey, "storage-config")
	if err != nil {
		return err
	}
	store := repo.New(db)
	srv, err := server.New(cfg, version, store, auth.New(store, jwtKey), storageSealer)
	if err != nil {
		return err
	}
	return server.RunHTTP(ctx, cfg.Server.Addr, srv.Handler())
}
func setupLogger(cfg config.Config) error {
	level := slog.LevelInfo
	switch strings.ToLower(cfg.Log.Level) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	var out io.Writer = os.Stdout
	if cfg.Log.File != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.Log.File), 0o700); err != nil {
			return err
		}
		out = io.MultiWriter(os.Stdout, &lumberjack.Logger{Filename: cfg.Log.File, MaxSize: cfg.Log.MaxSizeMB, MaxBackups: cfg.Log.MaxBackups, Compress: true})
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})))
	return nil
}

func CreateAdmin(ctx context.Context, cfg config.Config, username string, passwordStdin bool) error {
	username = strings.ToLower(strings.TrimSpace(username))
	if !validate.Username(username) {
		return errors.New("username must match [a-z0-9_]{3,50}")
	}
	var password string
	if passwordStdin {
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		password = strings.TrimSpace(line)
	} else {
		fmt.Fprint(os.Stderr, "Password: ")
		raw, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return err
		}
		password = string(raw)
	}
	if len([]rune(password)) < 12 {
		return errors.New("password must contain at least 12 characters")
	}
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	row := model.User{Username: username, PasswordHash: hash, Nickname: username, Role: 1, Status: 1}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("create administrator: %w", err)
	}
	fmt.Printf("Administrator %s created.\n", username)
	return nil
}
