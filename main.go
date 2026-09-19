package main

import (
	"os"
	"path/filepath"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/files"
	"github.com/portapps/portapps/v3/pkg/log"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: false,
	}

	// Init app
	if app, err = portapps.NewWithCfg("tabby-portable", "Tabby", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	if err := os.MkdirAll(app.DataPath, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("Cannot create data directory.")
	}
	app.Process = filepath.Join(app.AppPath, "Tabby.exe")
	app.Args = []string{
		"--user-data-dir=" + app.DataPath,
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			files.Cleanup(filepath.Join(os.Getenv("APPDATA"), "tabby"))
		}()
	}

	configFile := filepath.Join(app.DataPath, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Info().Msg("Creating default config...")
		if err := os.WriteFile(configFile, []byte(`enableAutomaticUpdates: false`), 0644); err != nil {
			log.Error().Err(err).Msg("Cannot write default config")
		}
	}
	if err := files.ReplaceByPrefix(configFile, "enableAutomaticUpdates:", "enableAutomaticUpdates: false"); err != nil {
		log.Fatal().Err(err).Msg("Cannot set enableAutomaticUpdates property")
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
