/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"os/exec"

	"fmt"

	"github.com/spf13/cobra"
)

func runCommand(cmdName string, args ...string) error {
	// Create a new command object
	cmd := exec.Command(cmdName, args...)

	// Get the standard output and standard error
	out, err := cmd.CombinedOutput()

	// Print the output of the command
	fmt.Printf("Output: %s\n", out)
	return err
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
		reader := bufio.NewReader(os.Stdin)

		// Prompt the user for input
		fmt.Print("Enter project name: ")

		// Read the input from the user
		input, err := reader.ReadString('\n')
		fmt.Println(input == "testtt")

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
		// fmt.Printf("You entered: %s\n", input)
		// cmd.InOrStdin()
		fmt.Println("bloom called")
		err = runCommand("npm", "create", "vite@latest", input, "--", "--template", "react-swc-ts")
		if err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			return
		}
		// time.Sleep(3000 * time.Millisecond)
		fmt.Println(os.ReadDir("./testbest"))
		err = os.Chdir(input)
		if err != nil {
			fmt.Printf("Error changing directory: %v\n", err)
			return
		}
		// Delete index.css and App.css
		filesToDelete := []string{"src/index.css", "src/App.css"}
		for _, file := range filesToDelete {
			err = os.Remove(file)
			if err != nil {
				fmt.Printf("Error deleting %s: %v\n", file, err)
			} else {
				fmt.Printf("Deleted %s\n", file)
			}
		}
		err = runCommand("npm", "install", "-D", "tailwindcss", "postcss", "autoprefixer")
		if err != nil {
			fmt.Printf("Error installing tailwind or dependencies: %v\n", err)
			return
		}
		err = runCommand("npx", "tailwindcss", "init", "-p")
		if err != nil {
			fmt.Printf("Error initializing tailwind: %v\n", err)
			return
		}

	},
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
