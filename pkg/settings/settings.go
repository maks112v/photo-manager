package settings

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Settings struct {
	BackupFolder string `toml:"backup_folder"`
	SourceFolder string `toml:"source_folder"`
}

//go:generate mockery --name SettingsProvider
type SettingsProvider interface {
	GetSettings() (*Settings, error)
	SaveSettings(settings *Settings) error
}

func GetSettings() (*Settings, error) {
	settingsFilePath := SettingsFilePath()

	fmt.Println("User Path", settingsFilePath)
	// Attempt to open the settings file
	file, err := os.Open(settingsFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the TOML content into a Settings struct
	var settings Settings
	err = toml.Unmarshal(content, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// SaveSettings writes the provided Settings struct to the settings.toml file.
func SaveSettings(settings *Settings) error {
	// Marshal the Settings struct to TOML format
	content, err := toml.Marshal(settings)
	if err != nil {
		return err
	}

	// Write the TOML content to the settings file
	settingsFilePath := SettingsFilePath()
	fmt.Println("User Path", settingsFilePath)
	err = os.WriteFile(settingsFilePath, content, 0644) // Using 0644 as the file permission
	if err != nil {
		return err
	}

	return nil
}

func SettingsFilePath() string {
	appDir, _ := os.UserConfigDir()
	settingsDir := filepath.Join(appDir, "photomanager")

	if _, err := os.Stat(settingsDir); os.IsNotExist(err) {
		os.Mkdir(settingsDir, os.ModePerm)
	}

	return filepath.Join(settingsDir, "settings.toml")
}
