/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

// settingsViewCmd represents the settingsView command
var settingsViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"v"},
	Short:   "View current settings",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		setting, err := settings.GetSettings()
		if err != nil {
			fmt.Println("Error reading settings: ", err)
			return
		}

		fmt.Println("Backup Folder Path: ", setting.BackupFolder)
		fmt.Println("Source Folder Path: ", setting.SourceFolder)
	},
}

func init() {
	settingsCmd.AddCommand(settingsViewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// settingsViewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// settingsViewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
