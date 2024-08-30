/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"os"

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
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		// fmt.Printf("You entered: %s\n", input)
		// cmd.InOrStdin()
		fmt.Println("bloom called")
		runCommand("npm", "create", "vite@latest", input, "--", "--template", "react-swc-ts")
		// err := runCommand("npm", "install")
		// if err != nil {
		// 	fmt.Printf("Error: %v\n", err)
		// }
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
