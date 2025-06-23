package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

const (
	algorandDir      = "/opt/algorand"
	defaultsDir      = "/opt/algorand/algorand-defaults"
	defaultConfigDir = "/opt/algorand/default-config"
)

func main() {
	if err := initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command provided\n")
		os.Exit(1)
	}

	if err := execCommand(args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		os.Exit(1)
	}
}

func execCommand(args []string) error {
	return syscall.Exec(args[0], args, os.Environ())
}

func initialize() error {
	genesisPath := filepath.Join(algorandDir, ".algorand", "genesis.json")
	if _, err := os.Stat(genesisPath); os.IsNotExist(err) {
		if err := setupGenesis(); err != nil {
			return fmt.Errorf("Failed to set up genesis: %w", err)
		}

		configPath := filepath.Join(algorandDir, ".algorand", "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			if err := copyFile(
				filepath.Join(defaultConfigDir, "config.json"),
				configPath,
			); err != nil {
				return fmt.Errorf("Failed to copy config file: %w", err)
			}
		}
	}
	return nil
}

func setupGenesis() error {
	network := os.Getenv("ALGOD_NETWORK")
	if network == "" {
		network = "mainnet"
	}

	sourceFile := filepath.Join(defaultsDir, fmt.Sprintf("genesis-%s.json", network))
	destFile := filepath.Join(algorandDir, ".algorand", "genesis.json")

	return copyFile(sourceFile, destFile)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
