package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	NodeURL        string
	NodeUser       string
	NodePass       string
	Checks         []string
	Addr           string
	MinConnections int
}

const MIN_CONNECTIONS = 8

func configFromEnv() (Config, error) {
	nodeURL := os.Getenv("NODE_URL")
	if nodeURL == "" {
		return Config{}, fmt.Errorf("NODE_URL is required")
	}

	checks := []string{"blockdownload"}
	if raw := os.Getenv("CHECKS"); raw != "" {
		checks = strings.Split(raw, ",")
	}

	for _, name := range checks {
		if _, ok := registry[name]; !ok {
			return Config{}, fmt.Errorf("unknown check %q (available: %s)", name, availableChecks())
		}
	}

	minConns := MIN_CONNECTIONS
	if raw := os.Getenv("MIN_CONNECTIONS"); raw != "" {
		n, err := strconv.Atoi(raw)

		if err != nil || n < 0 {
			return Config{}, fmt.Errorf("MIN_CONNECTIONS must be a non-negative integer")
		}

		minConns = n
	}

	return Config{
		NodeURL:        nodeURL,
		NodeUser:       os.Getenv("NODE_USER"),
		NodePass:       os.Getenv("NODE_PASS"),
		Checks:         checks,
		Addr:           ":8080",
		MinConnections: minConns,
	}, nil
}

func availableChecks() string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return strings.Join(names, ", ")
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cfg, err := configFromEnv()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("starting", "addr", cfg.Addr, "checks", cfg.Checks, "node", cfg.NodeURL)

	client := &http.Client{Timeout: 10 * time.Second}

	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", handleReadyz(client, cfg, logger))

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func handleReadyz(client *http.Client, cfg Config, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		for _, name := range cfg.Checks {
			if err := registry[name](ctx, client, cfg); err != nil {
				logger.Warn("check failed", "check", name, "error", err)
				http.Error(w, fmt.Sprintf("check %q failed: %v", name, err), http.StatusServiceUnavailable)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
