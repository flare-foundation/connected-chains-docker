package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
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

	// Switch to non-root
	if os.Geteuid() == 0 {
		if err := switchToAlgoUserAndExec(args); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to switch user: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := execCommand(args); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
			os.Exit(1)
		}
	}
}

func execCommand(args []string) error {
	return syscall.Exec(args[0], args, os.Environ())
}

func switchToAlgoUserAndExec(args []string) error {
	gosuArgs := append([]string{"/usr/bin/gosu", "algo"}, args...)
	return syscall.Exec(gosuArgs[0], gosuArgs, os.Environ())
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

		algoUID, err := strconv.Atoi(os.Getenv("USER_ID"))
		if err != nil {
			return fmt.Errorf("failed to parse USER_ID: %w", err)
		}
		algoGID, err := strconv.Atoi(os.Getenv("GROUP_ID"))
		if err != nil {
			return fmt.Errorf("failed to parse GROUP_ID: %w", err)
		}

		if err := chownRecursive(algorandDir, algoUID, algoGID); err != nil {
			return fmt.Errorf("failed to change ownership: %w", err)
		}
	}
	return nil
}

func chownRecursive(path string, uid, gid int) error {
	return filepath.WalkDir(path, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Chown(name, uid, gid)
	})
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
