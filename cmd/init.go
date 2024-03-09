/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup the photo manager",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		newSettings := settings.Settings{}

		fmt.Println("Backup Folder Path:")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input")
		}
		newSettings.BackupFolder = strings.TrimSuffix(line, "\n")

		fmt.Println("Source Folder Path:")
		reader = bufio.NewReader(os.Stdin)
		line, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input")
		}
		newSettings.SourceFolder = strings.TrimSuffix(line, "\n")

		settings.SaveSettings(&newSettings)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
