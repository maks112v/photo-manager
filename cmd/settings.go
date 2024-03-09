/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

var s *settings.Settings

var settingsCmd = &cobra.Command{
	Use:     "settings",
	Aliases: []string{"s"},
	Short:   "View and modify settings",
}

var settingsViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"v"},
	Short:   "View current settings",
	Run: func(cmd *cobra.Command, args []string) {
		setting, err := s.GetSettings()
		if err != nil {
			fmt.Println("Error reading settings: ", err)
			return
		}

		fmt.Println("Settings Path: ", settings.SettingsFilePath())
		fmt.Println("Backup Folder Path: ", setting.BackupFolder)
		fmt.Println("Source Folder Path: ", setting.SourceFolder)
		fmt.Println("Album Path Pattern: ", setting.AlbumPathPattern)
		fmt.Println("Photo Name Pattern: ", setting.PhotoNamePattern)
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
	settingsCmd.AddCommand(settingsViewCmd)

	s = settings.New()

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// settingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// settingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
