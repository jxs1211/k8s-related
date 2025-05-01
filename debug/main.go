package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "A simple CLI application",
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new resource",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("Creating resource with name: %s\n", name)
		},
	}

	createCmd.Flags().StringP("name", "n", "", "Name of the resource to create")
	createCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
