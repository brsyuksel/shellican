package main

import (
	"fmt"
	"os"

	"github.com/brsyuksel/shellican/pkg/core"
)

// main is the entry point for the application.
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create-shell":
		if len(os.Args) < 3 {
			fmt.Println("Usage: shellican create-shell <collection> [name]")
			os.Exit(1)
		}
		collection := os.Args[2]
		var name string
		if len(os.Args) > 3 {
			name = os.Args[3]
		}

		if err := core.CreateShell(collection, name); err != nil {
			fmt.Printf("Error creating shell helper: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Shell helper for '%s' created successfully.\n", collection)

	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Usage:")
			fmt.Println("  shellican new collection <name>")
			fmt.Println("  shellican new runnables <collection> <name>")
			os.Exit(1)
		}

		subCmd := os.Args[2]
		switch subCmd {
		case "collection":
			if len(os.Args) < 4 {
				fmt.Println("Usage: shellican new collection <name>")
				os.Exit(1)
			}
			name := os.Args[3]
			if err := core.CreateCollection(name); err != nil {
				fmt.Printf("Error creating collection: %v\n", err)
				os.Exit(1)
			}
		case "runnable":
			if len(os.Args) < 5 {
				fmt.Println("Usage: shellican new runnable <collection> <name>")
				os.Exit(1)
			}
			collection := os.Args[3]
			name := os.Args[4]
			if err := core.CreateRunnable(collection, name); err != nil {
				fmt.Printf("Error creating runnable: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unknown type: %s. Use 'collection' or 'runnable'\n", subCmd)
			os.Exit(1)
		}

	case "show":
		if len(os.Args) < 3 {
			fmt.Println("Usage:")
			fmt.Println("  shellican show collection <name> [--readme]")
			fmt.Println("  shellican show runnable <collection> <name> [--readme]")
			os.Exit(1)
		}

		subCmd := os.Args[2]

		showReadme := false
		args := os.Args[3:]
		if len(args) > 0 && args[len(args)-1] == "--readme" {
			showReadme = true
			args = args[:len(args)-1]
		}

		switch subCmd {
		case "collection":
			if len(args) < 1 {
				fmt.Println("Usage: shellican show collection <name> [--readme]")
				os.Exit(1)
			}
			name := args[0]
			if err := core.ShowCollection(name, showReadme); err != nil {
				fmt.Printf("Error showing collection: %v\n", err)
				os.Exit(1)
			}
		case "runnable":
			if len(args) < 2 {
				fmt.Println("Usage: shellican show runnable <collection> <name> [--readme]")
				os.Exit(1)
			}
			collection := args[0]
			name := args[1]
			if err := core.ShowRunnable(collection, name, showReadme); err != nil {
				fmt.Printf("Error showing runnable: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unknown type: %s. Use 'collection' or 'runnable'\n", subCmd)
			os.Exit(1)
		}

	case "list":
		if len(os.Args) < 3 {
			fmt.Println("Usage:")
			fmt.Println("  shellican list collections")
			fmt.Println("  shellican list runnables <collection>")
			os.Exit(1)
		}

		subCmd := os.Args[2]
		switch subCmd {
		case "collections":
			if err := core.ListCollections(); err != nil {
				fmt.Printf("Error listing collections: %v\n", err)
				os.Exit(1)
			}
		case "runnables":
			if len(os.Args) < 4 {
				fmt.Println("Usage: shellican list runnables <collection>")
				os.Exit(1)
			}
			collection := os.Args[3]
			if err := core.ListRunnables(collection); err != nil {
				fmt.Printf("Error listing runnables: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unknown type: %s. Use 'collections' or 'runnables'\n", subCmd)
			os.Exit(1)
		}

	case "import":
		if len(os.Args) < 3 {
			fmt.Println("Usage: shellican import <source> [name]")
			os.Exit(1)
		}
		source := os.Args[2]
		var name string
		if len(os.Args) > 3 {
			name = os.Args[3]
		}

		if err := core.ImportCollection(source, name); err != nil {
			fmt.Printf("Error importing collection: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Collection imported successfully.\n")

	case "export":
		if len(os.Args) < 3 {
			fmt.Println("Usage: shellican export <collection> [output]")
			os.Exit(1)
		}
		collection := os.Args[2]
		var output string
		if len(os.Args) > 3 {
			output = os.Args[3]
		}

		if err := core.ExportCollection(collection, output); err != nil {
			fmt.Printf("Error exporting collection: %v\n", err)
			os.Exit(1)
		}

	case "run":
		if len(os.Args) < 4 {
			fmt.Println("Usage: shellican run <collection> <runnable> [args...]")
			os.Exit(1)
		}

		collection := os.Args[2]
		scriptName := os.Args[3]
		scriptArgs := os.Args[4:]

		ctx, err := core.ResolveCommand(collection, []string{scriptName})
		if err != nil {
			fmt.Printf("Error resolving command: %v\n", err)
			os.Exit(1)
		}

		if err := core.ExecuteContext(ctx, scriptArgs); err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			os.Exit(1)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

// printUsage displays the usage information for the CLI.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  shellican run <collection> <runnable> [args...]")
	fmt.Println("  shellican create-shell <collection> [name]")
	fmt.Println("  shellican new collection <name>")
	fmt.Println("  shellican new runnable <collection> <name>")
	fmt.Println("  shellican list collections")
	fmt.Println("  shellican list runnables <collection>")
	fmt.Println("  shellican show collection <name> [--readme]")
	fmt.Println("  shellican show runnable <collection> <name> [--readme]")
	fmt.Println("  shellican import <source> [name]")
	fmt.Println("  shellican export <collection> [output]")
}
