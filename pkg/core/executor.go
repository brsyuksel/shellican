package core

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/brsyuksel/shellican/pkg/config"
)

// ExecutionContext holds the resolved state for execution
type ExecutionContext struct {
	RunnablePath string // Directory containing runnable.yml or resolved file
	Config       *config.RunnableConfig
	Environments map[string]string
}

// ResolveCommand traverses the collection to find the target runnable.
// It supports both legacy file access (optional) and new structured access.
func ResolveCommand(collection string, pathComponents []string) (*ExecutionContext, error) {
	rootDir, err := getRoot()
	if err != nil {
		return nil, err
	}
	rootDir = filepath.Join(rootDir, collection)

	currentPath := rootDir

	// Verify collection exists
	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("collection not found: %s", collection)
	}

	// Load Collection Config
	colCfg, err := config.LoadCollectionConfig(currentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection config: %w", err)
	}

	// Check if we have exactly one path component (the runnable name)
	if len(pathComponents) != 1 {
		return nil, fmt.Errorf("invalid command: expected exactly one runnable name, got %d components %v", len(pathComponents), pathComponents)
	}
	runName := pathComponents[0]

	// Check if explicit listing requirement is met
	if colCfg != nil {
		found := slices.Contains(colCfg.Runnables, runName)

		if !found {
			return nil, fmt.Errorf("runnable '%s' is not listed in collection.yml", runName)
		}
	} else {
		return nil, fmt.Errorf("collection.yml missing or runnables not listed")
	}

	// Target path is directly under collection root
	currentPath = filepath.Join(rootDir, runName)

	// Check if currentPath is a directory with runnable.yml
	info, err := os.Stat(currentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("runnable directory not found: %s", currentPath)
		}
		return nil, err
	}

	if info.IsDir() {
		runCfg, err := config.LoadRunnableConfig(currentPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load runnable config: %w", err)
		}
		if runCfg != nil {
			// Merge environments
			mergedEnvs := make(map[string]string)

			// 1. Collection environments (if any)
			maps.Copy(mergedEnvs, colCfg.Environments)

			// 2. Runnable environments (override collection)
			maps.Copy(mergedEnvs, runCfg.Environments)

			return &ExecutionContext{
				RunnablePath: currentPath,
				Config:       runCfg,
				Environments: mergedEnvs,
			}, nil
		}
		return nil, fmt.Errorf("directory found but no runnable.yml: %s", currentPath)
	}

	return nil, fmt.Errorf("target is a file, expected a directory with runnable.yml: %s", currentPath)
}

// ExecuteContext runs the resolved command with hooks
func ExecuteContext(ctx *ExecutionContext, args []string) error {
	cfg := ctx.Config

	// Run Before Hooks
	if cfg.Before != "" {
		fmt.Println(">> Running pre-hooks...")
		if err := runShell(cfg.Before, args, ctx.Environments, ctx.RunnablePath); err != nil {
			return fmt.Errorf("pre-hook failed: %s: %w", cfg.Before, err)
		}
	}

	// Run Main Command
	if cfg.Run == "" {
		return fmt.Errorf("no 'run' command specified in runnable.yml")
	}

	// Check if the run command points to an executable file
	runCmdPath := filepath.Join(ctx.RunnablePath, cfg.Run)
	info, err := os.Stat(runCmdPath)
	isScript := err == nil && !info.IsDir()

	if isScript {
		// Run as script file directly (Passes arguments correctly)
		err = runScript(runCmdPath, args, ctx.Environments)
	} else {
		// Run as inline shell command
		err = runShell(cfg.Run, args, ctx.Environments, ctx.RunnablePath)
	}

	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	// Run After Hooks
	if cfg.After != "" {
		fmt.Println(">> Running post-hooks...")
		if err := runShell(cfg.After, args, ctx.Environments, ctx.RunnablePath); err != nil {
			fmt.Printf("Warning: post-hook failed: %s: %v\n", cfg.After, err)
			// Don't fail the whole execution for post-hooks?
		}
	}

	return nil
}

func runScript(path string, args []string, envs map[string]string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Inject environments
	cmd.Env = os.Environ()
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd.Run()
}

func runShell(command string, args []string, envs map[string]string, dir string) error {
	// Basic shell execution.
	// We pass args as environment variables or just appended?
	// "inline scripts... commands"
	// A simple way is sh -c "command" -- args...

	shellCmd := exec.Command("/bin/sh", "-c", command, "inline-script")
	shellCmd.Args = append(shellCmd.Args, args...)
	shellCmd.Dir = dir
	shellCmd.Stdin = os.Stdin
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr

	// Inject environments
	shellCmd.Env = os.Environ()
	for k, v := range envs {
		shellCmd.Env = append(shellCmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return shellCmd.Run()
}
