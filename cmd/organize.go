/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/maks112v/photomanager/pkg/file"
	"github.com/maks112v/photomanager/pkg/photomanager"
	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

// organizeCmd represents the organize command
var organizeCmd = &cobra.Command{
	Use:     "organize",
	Aliases: []string{"o"},
	Short:   "Run the organizer",
	Long:    `Run the organizer to move photos to the right location by date / time & event name.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		photomanager := photomanager.New(file.New(), settings.New())

		photomanager.PreRunValidation()
	},
	Run: func(cmd *cobra.Command, args []string) {
		photomanager := photomanager.New(file.New(), settings.New())

		photomanager.Organize()
	},
}

func init() {

	rootCmd.AddCommand(organizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// organizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	organizeCmd.Flags().BoolP("duration", "d", false, "Duration required to create a new album")
}
