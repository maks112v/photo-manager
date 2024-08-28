package photomanager

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"sort"
	"time"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/maks112v/photomanager/pkg/file"
	"github.com/maks112v/photomanager/pkg/settings"
)

type PhotoManager struct {
	file     file.FileProvider
	settings settings.SettingsProvider
}

func New(fileprovider file.FileProvider, settings settings.SettingsProvider) *PhotoManager {
	return &PhotoManager{
		file:     fileprovider,
		settings: settings,
	}
}

type Album struct {
	Name       string
	Year       int
	Month      int
	PhotoCount int
	FirstPhoto time.Time
	Photos     []file.PhotoFile
	Path       string
}

func (p *PhotoManager) PreRunValidation() error {

	settings, err := p.settings.GetSettings()
	if err != nil {
		return err
	}

	// Validate source and target folders are set
	if settings.SourceFolder == "" {
		return errors.New("source folder not set")
	}

	if settings.BackupFolder == "" {
		return errors.New("backup folder not set")
	}

	return nil
}

func (p *PhotoManager) Organize(durationOverride int) error {
	setting, err := p.settings.GetSettings()

	var runDuration int
	if durationOverride > 0 {
		runDuration = durationOverride
	} else {
		runDuration = setting.DurationBetweenAlbums
	}

	fmt.Printf("Searching for albums with a duration of %d hours\n", runDuration)

	if err != nil {
		fmt.Println("Error reading settings: ", err)
		return err
	}

	photoFiles, err := p.file.GetAllFiles(setting.SourceFolder)
	if err != nil {
		fmt.Println("Error reading files: ", err)
		return err
	}

	sort.Slice(photoFiles, func(i, j int) bool {
		return photoFiles[i].CreatedAt.Before(photoFiles[j].CreatedAt)
	})

	var albums []Album
	for _, photoFile := range photoFiles {
		if len(albums) == 0 || photoFile.CreatedAt.Sub(albums[len(albums)-1].Photos[len(albums[len(albums)-1].Photos)-1].CreatedAt).Hours() > float64(runDuration) {
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

	input := confirmation.New("Save new albums?", confirmation.Yes)

	if ready, err := input.RunPrompt(); !ready || err != nil {
		fmt.Println("Aborted")
		return err
	}

	for _, album := range albums {
		fmt.Println("Creating album", album.Path)

		if _, err := os.Stat(setting.BackupFolder + "/" + album.Path); os.IsNotExist(err) {
			if err := os.Mkdir(setting.BackupFolder+"/"+album.Path, os.ModePerm); err != nil {
				fmt.Println("Error creating album folder: ", err)
			}
		}

		basePath := setting.BackupFolder + "/" + album.Path
		for i, photo := range album.Photos {
			fmt.Println("Copying photo ", i+1, "/", len(album.Photos))
			photo.Number = i + 1
			tmpl, err := template.New("photo-file-name").Parse(setting.PhotoNamePattern)
			if err != nil {
				panic(err)
			}
			var buf bytes.Buffer
			err = tmpl.Execute(&buf, photo)
			if err != nil {
				panic(err)
			}

			fileName := buf.String()
			p.file.CopyFile(photo.Path, path.Join(basePath, fileName))
		}
	}

	return nil
}
