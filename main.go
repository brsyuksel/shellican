package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/brsyuksel/shellican/pkg/core"
	"github.com/spf13/cobra"
)

// version is injected at build time.
var version = "dev"

func getVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return version
}

// main is the entry point for the application.
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "shellican",
	Short: "A CLI tool for managing shell script collections",
	Long:  `shellican is a CLI tool that helps you organize, run, and share shell script collections.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getVersion())
	},
}

var createShellCmd = &cobra.Command{
	Use:   "create-shell <collection> [name]",
	Short: "Create a shell helper for a collection",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		var name string
		if len(args) > 1 {
			name = args[1]
		}

		if err := core.CreateShell(collection, name); err != nil {
			fmt.Printf("Error creating shell helper: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Shell helper for '%s' created successfully.\n", collection)
	},
}

var newCmd = &cobra.Command{
	Use:   "new <collection> [runnable]",
	Short: "Create a new collection or runnable",
	Long: `Create a new collection or runnable.
  If only collection is provided, creates a new collection.
  If both collection and runnable are provided, creates a new runnable in the collection.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			// creates collection
			name := args[0]
			if err := core.CreateCollection(name); err != nil {
				fmt.Printf("Error creating collection: %v\n", err)
				os.Exit(1)
			}
		} else {
			// creates runnable
			collection := args[0]
			name := args[1]
			if err := core.CreateRunnable(collection, name); err != nil {
				fmt.Printf("Error creating runnable: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show <collection> [runnable]",
	Short: "Show details of a collection or runnable",
	Long: `Show details of a collection or runnable.
  If only collection is provided, shows the collection details.
  If both collection and runnable are provided, shows the runnable details.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		showReadme, _ := cmd.Flags().GetBool("readme")

		if len(args) == 1 {
			// show collection
			name := args[0]
			if err := core.ShowCollection(name, showReadme); err != nil {
				fmt.Printf("Error showing collection: %v\n", err)
				os.Exit(1)
			}
		} else {
			// show runnable
			collection := args[0]
			name := args[1]
			if err := core.ShowRunnable(collection, name, showReadme); err != nil {
				fmt.Printf("Error showing runnable: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list [collection]",
	Short: "List collections or runnables",
	Long: `List collections or runnables in a collection.
  If no collection is provided, lists all collections.
  If collection is provided, lists all runnables in that collection.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// list collections
			if err := core.ListCollections(); err != nil {
				fmt.Printf("Error listing collections: %v\n", err)
				os.Exit(1)
			}
		} else {
			// list runnables
			collection := args[0]
			if err := core.ListRunnables(collection); err != nil {
				fmt.Printf("Error listing runnables: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import <source> [name]",
	Short: "Import a collection from a source",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		var name string
		if len(args) > 1 {
			name = args[1]
		}

		if err := core.ImportCollection(source, name); err != nil {
			fmt.Printf("Error importing collection: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Collection imported successfully.\n")
	},
}

var exportCmd = &cobra.Command{
	Use:   "export <collection> [output]",
	Short: "Export a collection",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		var output string
		if len(args) > 1 {
			output = args[1]
		}

		if err := core.ExportCollection(collection, output); err != nil {
			fmt.Printf("Error exporting collection: %v\n", err)
			os.Exit(1)
		}
	},
}

var runCmd = &cobra.Command{
	Use:   "run <collection> <runnable> [args...]",
	Short: "Run a runnable from a collection",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		scriptName := args[1]
		scriptArgs := args[2:]

		ctx, err := core.ResolveCommand(collection, []string{scriptName})
		if err != nil {
			fmt.Printf("Error resolving command: %v\n", err)
			os.Exit(1)
		}

		if err := core.ExecuteContext(ctx, scriptArgs); err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add flags
	showCmd.Flags().Bool("readme", false, "Show README content")

	// Add commands to root
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createShellCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(runCmd)
}
