package config

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
)

// Config is the configuration struct
type Config struct {
	Directories     []string
	IgnoredPatterns []string `yaml:",omitempty"`
	IgnoreSymlinks  bool     `yaml:",omitempty"`
	IgnoreHidden    bool     `yaml:",omitempty"`
}

var defaultConfigTemplate string = `directories:
  - "{{ .HomeDir }}"

ignoredPatterns:
  - "{{ .CacheDir }}"

ignoreSymlinks: true
ignoreHidden: false`

// GetConfigFileName returns the config file name
func GetConfigFileName() string {
	confDir := xdg.ConfigHome
	home := xdg.Home

	if confDir == "" {
		confDir = filepath.Join(home, ".config")
	}

	confDir = filepath.Join(confDir, "glocate")

	configFileName := filepath.Join(confDir, "glocate.yaml")

	return configFileName
}

// WriteDefaultConfigIfNotExist writes the default config file if it doesn't exist
func WriteDefaultConfigIfNotExist(configFileName string) {
	if _, err := os.Stat(configFileName); err == nil {
		return
	}

	log.Errorf("no config file found, creating default config file at %s", configFileName)

	os.MkdirAll(filepath.Dir(configFileName), os.ModePerm)

	f, err := os.Create(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t, err := template.New("config").Parse(defaultConfigTemplate)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(f, struct {
		HomeDir  string
		CacheDir string
	}{
		HomeDir:  xdg.Home,
		CacheDir: xdg.CacheHome,
	})

	if err != nil {
		log.Fatal(err)
	}
}
