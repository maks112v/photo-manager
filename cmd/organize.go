/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/maks112v/photomanager/pkg/file"
	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

type File struct {
	Path      string
	CreatedAt time.Time
}

type Album struct {
	Name       string
	Year       int
	Month      int
	PhotoCount int
	FirstPhoto time.Time
	Photos     []File
}

// organizeCmd represents the organize command
var organizeCmd = &cobra.Command{
	Use:     "organize",
	Aliases: []string{"o"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("organize called")

		setting, err := settings.GetSettings()

		if err != nil {
			fmt.Println("Error reading settings: ", err)
			return
		}

		foundFiles, err := file.GetAllFiles(setting.SourceFolder)
		if err != nil {
			fmt.Println("Error reading files: ", err)
			return
		}

		var files []File

		for _, filePath := range foundFiles {
			createdAt, err := file.CreatedAt(filePath)
			if err != nil {
				fmt.Println("Error reading created at: ", err)
				return
			}
			files = append(files, File{Path: filePath, CreatedAt: createdAt})
		}

		var albums []Album

		sort.Slice(files, func(i, j int) bool {
			return files[i].CreatedAt.Before(files[j].CreatedAt)
		})

		for _, file := range files {
			fmt.Println("File: ", file.Path)
			fmt.Println("Created at: ", file.CreatedAt)

			// Start a new album if this is the first file or if the gap exceeds 24 hours
			if len(albums) == 0 || file.CreatedAt.Sub(albums[len(albums)-1].Photos[len(albums[len(albums)-1].Photos)-1].CreatedAt).Hours() > 24 {
				newAlbum := Album{
					Name:       fmt.Sprintf("Album %d", len(albums)+1),
					Year:       file.CreatedAt.Year(),
					Month:      int(file.CreatedAt.Month()),
					FirstPhoto: file.CreatedAt,
					Photos:     []File{file},
				}
				albums = append(albums, newAlbum)
			} else {
				// Add to the current album
				currentAlbum := &albums[len(albums)-1]
				currentAlbum.Photos = append(currentAlbum.Photos, file)
				currentAlbum.PhotoCount = len(currentAlbum.Photos)
			}
		}

		// Print album details
		for i, album := range albums {
			fmt.Printf("Album %d: %+v\n", i+1, album)
			fmt.Println("Photos in album:", len(album.Photos))
			for _, photo := range album.Photos {
				fmt.Println("Photo:", photo.Path, "Created at:", photo.CreatedAt)
			}
		}
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
	// organizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
