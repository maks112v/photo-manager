/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup the photo manager",
	Run: func(cmd *cobra.Command, args []string) {
		newSettings := settings.SettingStruct{}

		input := textinput.New("Source Folder Path")
		input.Placeholder = "The source folder cannot be empty"
		input.Validate = func(input string) error {
			if input == "" {
				return fmt.Errorf("the source folder cannot be empty")
			}

			if _, err := os.Stat(input); os.IsNotExist(err) {
				return fmt.Errorf("the source folder does not exist")
			}

			return nil
		}

		folderPath, err := input.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			os.Exit(1)
		}

		newSettings.SourceFolder = folderPath

		input = textinput.New("Backup Folder Path")
		input.Placeholder = "The source folder cannot be empty"
		input.Validate = func(input string) error {
			if input == "" {
				return fmt.Errorf("the source folder cannot be empty")
			}

			if _, err := os.Stat(input); os.IsNotExist(err) {
				return fmt.Errorf("the source folder does not exist")
			}

			if input == newSettings.SourceFolder {
				return fmt.Errorf("the backup folder cannot be the same as the source folder")
			}

			return nil
		}

		folderPath, err = input.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			os.Exit(1)
		}

		newSettings.BackupFolder = folderPath

		input = textinput.New("Album Folder Pattern")
		input.Placeholder = "The album pattern for the photos"
		input.InitialValue = "{{.Year}}-{{.Month}} {{.Name}}"

		albumPath, err := input.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			os.Exit(1)
		}

		newSettings.AlbumPathPattern = albumPath

		input = textinput.New("File Name Pattern")
		input.Placeholder = "The file pattern for the photos"
		input.InitialValue = "{{.Name}}{{.Ext}}"

		fileNamePattern, err := input.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			os.Exit(1)
		}

		newSettings.PhotoNamePattern = fileNamePattern

		settings.New().SaveSettings(&newSettings)
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
