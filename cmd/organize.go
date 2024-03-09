/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/maks112v/photomanager/pkg/file"
	"github.com/maks112v/photomanager/pkg/settings"
	"github.com/spf13/cobra"
)

type Album struct {
	Name       string
	Year       int
	Month      int
	PhotoCount int
	FirstPhoto time.Time
	Photos     []file.PhotoFile
	Path       string
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
		setting, err := settings.GetSettings()

		if err != nil {
			fmt.Println("Error reading settings: ", err)
			return
		}

		photoFiles, err := file.GetAllFiles(setting.SourceFolder)
		if err != nil {
			fmt.Println("Error reading files: ", err)
			return
		}

		var albums []Album

		sort.Slice(photoFiles, func(i, j int) bool {
			return photoFiles[i].CreatedAt.Before(photoFiles[j].CreatedAt)
		})

		for _, photoFile := range photoFiles {
			if len(albums) == 0 || photoFile.CreatedAt.Sub(albums[len(albums)-1].Photos[len(albums[len(albums)-1].Photos)-1].CreatedAt).Hours() > 24 {
				newAlbum := Album{
					Name:       fmt.Sprintf("Album %d", len(albums)+1),
					Year:       photoFile.CreatedAt.Year(),
					Month:      int(photoFile.CreatedAt.Month()),
					FirstPhoto: photoFile.CreatedAt,
					Photos:     []file.PhotoFile{photoFile},
				}

				tmpl, err := template.New("album-path").Parse("{{.Year}}-{{.Month}} {{.Name}}")
				if err != nil {
					panic(err)
				}
				var buf bytes.Buffer
				err = tmpl.Execute(&buf, newAlbum)
				if err != nil {
					panic(err)
				}

				newAlbum.Path = buf.String()
				albums = append(albums, newAlbum)
			} else {
				// Add to the current album
				currentAlbum := &albums[len(albums)-1]
				currentAlbum.Photos = append(currentAlbum.Photos, photoFile)
				currentAlbum.PhotoCount = len(currentAlbum.Photos)
			}
		}

		for _, album := range albums {
			fmt.Println(album.Path, "Photos in album:", len(album.Photos))
		}

		input := confirmation.New("Run Files?", confirmation.Yes)

		if ready, err := input.RunPrompt(); !ready || err != nil {
			fmt.Println("Aborted")
			return
		}

		for _, album := range albums {
			fmt.Println("Creating album", album.Path)

			if _, err := os.Stat(setting.BackupFolder + "/" + album.Path); os.IsNotExist(err) {
				if err := os.Mkdir(setting.BackupFolder+"/"+album.Path, os.ModePerm); err != nil {
					fmt.Println("Error creating album folder: ", err)
				}
			}

			for i, photo := range album.Photos {
				filePath := fmt.Sprintf("%s/%s/%s-%d%s", setting.BackupFolder, album.Path, album.Path, i+1, photo.Ext)
				CopyFile(photo.Path, filePath)
			}
		}

	},
}

func CopyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file for writing
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure that any writes to the destination file are committed
	err = destinationFile.Sync()
	return err
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
