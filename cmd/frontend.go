/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Start the frontend static webserver",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("frontend called")
	},
}

func init() {
	runCmd.AddCommand(frontendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// frontendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// frontendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
