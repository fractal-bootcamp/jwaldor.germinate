/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"embed"
	"errors"
	"io"
	"os"
	"strings"
	"path/filepath"

	"os/exec"

	"fmt"

	"github.com/spf13/cobra"
)

//go:embed assets/*

var f embed.FS

func runCommand(cmdName string, args ...string) error {
	// Create a new command object
	cmd := exec.Command(cmdName, args...)

	// Get the standard output and standard error
	out, err := cmd.CombinedOutput()

	// Print the output of the command
	fmt.Printf("Output: %s\n", out)
	return err
}

func writeEmbeddedToDisk(embeddedPath string, outputPath string) error {
	// Attempt to open the embedded file or directory
	file, err := f.Open(embeddedPath)
	if err != nil {
		return fmt.Errorf("failed to open embedded path: %w", err)
	}
	defer file.Close()

	// If it's a directory, you can read its entries
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get info for embedded path: %w", err)
	}

	if info.IsDir() {
		// If it's a directory, create it and copy its contents
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		entries, err := f.ReadDir(embeddedPath)
		if err != nil {
			return fmt.Errorf("failed to read embedded directory: %w", err)
		}

		for _, entry := range entries {
			srcPath := filepath.Join(embeddedPath, entry.Name())
			dstPath := filepath.Join(outputPath, entry.Name())
			if err := writeEmbeddedToDisk(srcPath, dstPath); err != nil {
				return err
			}
		}
	} else {
		// If it's a file, copy it
		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			return fmt.Errorf("failed to copy content to output file: %w", err)
		}
	}

	return nil
}
func askYesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " (y/n): ")
	answer, _ := reader.ReadString('\n')
	answer = answer[:len(answer)-1] // Remove the newline character

	switch answer {
	case "y", "Y", "yes", "Yes":
		return true
	case "n", "N", "no", "No":
		return false
	default:
		fmt.Println("Invalid response. Please answer 'y' or 'n'.")
		return askYesNo(question) // Recursive call if the answer is invalid
	}
}

// bloomCmd represents the bloom command
var bloomCmd = &cobra.Command{
	Use:   "bloom",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		completedSetupSteps := []string{}

		reader := bufio.NewReader(os.Stdin)

		// Prompt the user for input
		fmt.Print("Enter project name: ")

		// Read the input from the user
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = strings.TrimSpace(input)

		if _, err := os.Stat(input); err == nil {
			fmt.Println("Project directory already exists")
			return
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
		} else {
			fmt.Println(err, "Provided directory could not be evaluated to already exist or not, aborting to be safe")
			return
			// Schrodinger: file may or may not exist. See err for details.

			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		}

		fmt.Println("bloom called")
		err = runCommand("npm", "create", "vite@latest", input, "--", "--template", "react-swc-ts")
		if err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			return
		}
		completedSetupSteps = append(completedSetupSteps, "✅ Created Vite project with React and TypeScript")

		err = os.Chdir(input)
		if err != nil {
			fmt.Printf("Error changing directory: %v\n", err)
			return
		}
		completedSetupSteps = append(completedSetupSteps, "✅ Changed to project directory")

		//Tailwind install section
		err = runCommand("npm", "install", "-D", "tailwindcss", "postcss", "autoprefixer")
		if err != nil {
			fmt.Printf("Error installing tailwind or dependencies: %v\n", err)
			return
		}
		completedSetupSteps = append(completedSetupSteps, "✅ Installed Tailwind CSS and dependencies")

		err = runCommand("npx", "tailwindcss", "init", "-p")
		if err != nil {
			fmt.Printf("Error initializing tailwind: %v\n", err)
			return
		}
		completedSetupSteps = append(completedSetupSteps, "✅ Initialized Tailwind CSS")

		// Delete index.css and App.css
		filesToDelete := []string{"src/index.css", "src/App.css", "tailwind.config.js"}
		for _, file := range filesToDelete {
			err = os.Remove(file)
			if err != nil {
				fmt.Printf("Error deleting %s: %v\n", file, err)
			} else {
				fmt.Printf("Deleted %s\n", file)
				completedSetupSteps = append(completedSetupSteps, fmt.Sprintf("✅ Removed default file for %s", file))
			}
		}

		indexCss := "assets/index.css"
		tailwindConfig := "assets/tailwind.config.js"
		output_index := "src/index.css"
		if err := writeEmbeddedToDisk(indexCss, output_index); err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Printf("File written successfully to %s\n", output_index)
			completedSetupSteps = append(completedSetupSteps, "✅ Created custom index.css")
		}

		output_tailwindconf := "tailwind.config.js"
		if err := writeEmbeddedToDisk(tailwindConfig, output_tailwindconf); err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Printf("File written successfully to %s\n", output_tailwindconf)
			completedSetupSteps = append(completedSetupSteps, "✅ Created custom tailwind.config.js")
		}
		if askYesNo("Would you like to set up a Prisma database?") {
			fmt.Println("Setting up Prisma database...")
			
			// Write the embedded database boilerplate folder to disk
			boilerplatePath := "assets/database-boilerplate-main"
			outputBoilerplatePath := "backend"
			if err := writeEmbeddedToDisk(boilerplatePath, outputBoilerplatePath); err != nil {
				fmt.Printf("Error writing database boilerplate folder: %v\n", err)
				return
			}
			fmt.Println("Created database boilerplate folder")
			
			completedSetupSteps = append(completedSetupSteps, "✅ Set up Prisma database")
			//clone boilerplate from Github
			// fmt.Println("Downloading boilerplate...")
			// err = runCommand("git", "clone", "git@github.com:fractal-bootcamp/database-boilerplate.git")
			// if err != nil {
			// 	fmt.Printf("Error cloning boilerplate: %v\n", err)
			// 	return
			// }
			// err = os.Chdir("database-boilerplate")
			// if err != nil {
			// 	fmt.Printf("Error changing directory: %v\n", err)
			// 	return
			// }
			// // Delete .git and .gitignore in the database-boilerplate repo
			// filesToDelete := []string{".gitignore"}
			// for _, file := range filesToDelete {
			// 	err = os.Remove(file)
			// 	if err != nil {
			// 		fmt.Printf("Error deleting %s: %v\n", file, err)
			// 	} else {
			// 		fmt.Printf("Deleted %s\n", file)
			// 	}
			// }
			// a_file := ".git"
			// err = os.RemoveAll(a_file)
			// if err != nil {
			// 	fmt.Printf("Error deleting %s: %v\n", a_file, err)
			// } else {
			// 	fmt.Printf("Deleted %s\n", a_file)
			// }

			// fmt.Println("Initializing database...")
			// err = runCommand("npx", "prisma", "generate")
			// if err != nil {
			// 	fmt.Printf("Error with npx prisma generate: %v\n", err)
			// 	return
			// }
			// fmt.Println("Performing initial migration...")
			// err = runCommand("npx", "prisma", "migrate", "reset")
			// if err != nil {
			// 	fmt.Printf("Error with npx prisma generate: %v\n", err)
			// 	return
			// }

		} else {
			fmt.Println("Proceeding without database...")
			completedSetupSteps = append(completedSetupSteps, "✅ Skipped database setup")
		}
		// Ask if the user wants to set up an Express backend
		// if askYesNo("Do you want to set up an Express backend?") {
		// 	fmt.Println("Setting up Express backend...")

		// 	// Create a backend directory
		// 	backendDir := filepath.Join(projectName, "backend")
		// 	if err := os.MkdirAll(backendDir, 0755); err != nil {
		// 		fmt.Printf("Error creating backend directory: %v\n", err)
		// 		return
		// 	}

		// 	// Change to the backend directory
		// 	if err := os.Chdir(backendDir); err != nil {
		// 		fmt.Printf("Error changing to backend directory: %v\n", err)
		// 		return
		// 	}

		// 	// Initialize npm and install express
		// 	if err := runCommand("npm", "init", "-y"); err != nil {
		// 		fmt.Printf("Error initializing npm: %v\n", err)
		// 		return
		// 	}

		// 	if err := runCommand("npm", "install", "express"); err != nil {
		// 		fmt.Printf("Error installing Express: %v\n", err)
		// 		return
		// 	}

		// 	// Copy the server.js file from embedded assets
		// 	if err := writeEmbeddedToDisk("assets/server.js", "server.js"); err != nil {
		// 		fmt.Printf("Error creating server.js: %v\n", err)
		// 		return
		// 	}

		// 	fmt.Println("Express backend set up successfully!")
		// 	completedSetupSteps = append(completedSetupSteps, "✅ Set up Express backend")

		// 	// Change back to the project root directory
		// 	if err := os.Chdir(".."); err != nil {
		// 		fmt.Printf("Error changing back to project root: %v\n", err)
		// 		return
		// 	}
		// } else {
		// 	fmt.Println("Skipping Express backend setup...")
		// 	completedSetupSteps = append(completedSetupSteps, "✅ Skipped Express backend setup")
		// }

		// Print completed setup steps
		fmt.Println("\nCompleted the following setup steps:")
		for _, step := range completedSetupSteps {
			fmt.Println(step)
		}
	},
	//react router option
	//organizing routes option
	//auth option
	//database option
	//make frontend and backend folders
}

func init() {
	rootCmd.AddCommand(bloomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bloomCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bloomCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
