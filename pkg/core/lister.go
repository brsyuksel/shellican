package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/brsyuksel/shellican/pkg/config"
)

// ListCollections prints available collections
func ListCollections() error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No collections found.")
			return nil
		}
		return fmt.Errorf("failed to list directory: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION")

	var names []string
	configs := make(map[string]*config.CollectionConfig)

	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
			cfg, _ := config.LoadCollectionConfig(filepath.Join(rootDir, entry.Name()))
			configs[entry.Name()] = cfg
		}
	}
	sort.Strings(names)

	for _, name := range names {
		desc := "No description"
		if cfg := configs[name]; cfg != nil && cfg.Help != "" {
			desc = cfg.Help
		}
		fmt.Fprintf(w, "%s\t%s\n", name, desc)
	}
	w.Flush()
	return nil
}

// ListRunnables prints available runnables for a collection
func ListRunnables(collectionName string) error {
	rootDir, err := getRoot()
	if err != nil {
		return err
	}
	collectionPath := filepath.Join(rootDir, collectionName)

	// Load collection to see explicit list
	colCfg, err := config.LoadCollectionConfig(collectionPath)
	if err != nil {
		return fmt.Errorf("failed to load collection config: %w", err)
	}
	if colCfg == nil {
		return fmt.Errorf("collection '%s' not found or invalid", collectionName)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION")

	// We only list what is in Runnables list
	if len(colCfg.Runnables) == 0 {
		fmt.Println("No runnables found.")
		return nil
	}

	for _, name := range colCfg.Runnables {
		runnablePath := filepath.Join(collectionPath, name)
		desc := "No description"

		runCfg, _ := config.LoadRunnableConfig(runnablePath)
		if runCfg != nil && runCfg.Help != "" {
			desc = runCfg.Help
		}

		fmt.Fprintf(w, "%s\t%s\n", name, desc)
	}

	w.Flush()
	return nil
}
