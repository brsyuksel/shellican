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

// ExecutionContext holds the execution state.
type ExecutionContext struct {
	RunnablePath string
	Config       *config.RunnableConfig
	Environments map[string]string
}

// ResolveCommand resolves a runnable from a collection.
func ResolveCommand(collection string, pathComponents []string) (*ExecutionContext, error) {
	rootDir, err := getRoot()
	if err != nil {
		return nil, err
	}
	rootDir = filepath.Join(rootDir, collection)

	currentPath := rootDir

	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("collection not found: %s", collection)
	}

	colCfg, err := config.LoadCollectionConfig(currentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection config: %w", err)
	}

	if len(pathComponents) != 1 {
		return nil, fmt.Errorf("invalid command: expected exactly one runnable name, got %d components %v", len(pathComponents), pathComponents)
	}
	runName := pathComponents[0]

	if colCfg != nil {
		found := slices.Contains(colCfg.Runnables, runName)

		if !found {
			return nil, fmt.Errorf("runnable '%s' is not listed in collection.yml", runName)
		}
	} else {
		return nil, fmt.Errorf("collection.yml missing or runnables not listed")
	}

	currentPath = filepath.Join(rootDir, runName)

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
			mergedEnvs := make(map[string]string)

			maps.Copy(mergedEnvs, colCfg.Environments)

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

// ExecuteContext executes a runnable.
func ExecuteContext(ctx *ExecutionContext, args []string) error {
	cfg := ctx.Config

	if cfg.Before != "" {
		fmt.Println(">> Running pre-hooks...")
		if err := runShell(cfg.Before, args, ctx.Environments, ctx.RunnablePath); err != nil {
			return fmt.Errorf("pre-hook failed: %s: %w", cfg.Before, err)
		}
	}

	if cfg.Run == "" {
		return fmt.Errorf("no 'run' command specified in runnable.yml")
	}

	runCmdPath := filepath.Join(ctx.RunnablePath, cfg.Run)
	info, err := os.Stat(runCmdPath)
	isScript := err == nil && !info.IsDir()

	if isScript {
		err = runScript(runCmdPath, args, ctx.Environments)
	} else {
		err = runShell(cfg.Run, args, ctx.Environments, ctx.RunnablePath)
	}

	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	if cfg.After != "" {
		fmt.Println(">> Running post-hooks...")
		if err := runShell(cfg.After, args, ctx.Environments, ctx.RunnablePath); err != nil {
			fmt.Printf("Warning: post-hook failed: %s: %v\n", cfg.After, err)
		}
	}

	return nil
}

func runScript(path string, args []string, envs map[string]string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd.Run()
}

func runShell(command string, args []string, envs map[string]string, dir string) error {

	shellCmd := exec.Command("/bin/sh", "-c", command, "inline-script")
	shellCmd.Args = append(shellCmd.Args, args...)
	shellCmd.Dir = dir
	shellCmd.Stdin = os.Stdin
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr

	shellCmd.Env = os.Environ()
	for k, v := range envs {
		shellCmd.Env = append(shellCmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return shellCmd.Run()
}
