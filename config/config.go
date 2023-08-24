package config

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
)

type Config struct {
	Directories     []string
	IgnoredPatterns []string `yaml:",omitempty"`
	IgnoreSymlinks  bool     `yaml:",omitempty"`
	GitIgnore       bool     `yaml:",omitempty"`
	IgnoreHidden    bool     `yaml:",omitempty"`
}

var defaultConfigTemplate string = `directories:
  - "{{ .HomeDir }}"

ignoredPatterns:
  - "{{ .CacheDir }}"

ignoreSymlinks: true
gitIgnore: false
ignoreHidden: false`

func GetConfigFileName() string {
	confDir := xdg.ConfigHome
	home, err := os.UserHomeDir()

	if confDir == "" {
		if err != nil {
			log.Fatal(err)
		}
		confDir = filepath.Join(home, ".config")
	}

	confDir = filepath.Join(confDir, "glocate")

	configFileName := filepath.Join(confDir, "glocate.yaml")

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		os.MkdirAll(confDir, os.ModePerm)
		writeDefaultConfig(configFileName, home)
	}

	return configFileName
}

func writeDefaultConfig(configFileName, home string) {
	f, err := os.Create(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.Errorf("no config file found, creating default config file at %s", configFileName)

	t, err := template.New("config").Parse(defaultConfigTemplate)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(f, struct {
		HomeDir  string
		CacheDir string
	}{
		HomeDir:  home,
		CacheDir: xdg.CacheHome,
	})

	if err != nil {
		log.Fatal(err)
	}
}
