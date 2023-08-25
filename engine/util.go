package engine

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
)

var dbFileName = getDBFileName()

// getDBFileName returns the path to the database file
func getDBFileName() string {
	glocateDir := xdg.DataHome + "/glocate"
	os.MkdirAll(glocateDir, os.ModePerm)
	return glocateDir + "/glocate.lz4"
}

func cleanPath(dir string) (string, error) {
	// expand variables if present
	dir = os.ExpandEnv(dir)
	// remove trailing slash
	dir = filepath.Clean(dir)
	// expand ~ if present
	dir = strings.Replace(dir, "~", os.Getenv("HOME"), 1)

	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return "", err
	}

	return dir, nil
}

// cleanDirectories removes trailing slashes and expands ~ and $var
func cleanDirectories(directories []string) []string {
	sanitizedDirs := make([]string, 0, len(directories))
	for _, dir := range directories {
		var err error
		if dir, err = cleanPath(dir); err != nil {
			log.Error(err)
			continue
		}
		sanitizedDirs = append(sanitizedDirs, dir)
	}
	return sanitizedDirs
}

// buildRegexes compiles all regexes and returns them
func buildRegexes(ignoredPatterns []string) []*regexp.Regexp {
	regexes := make([]*regexp.Regexp, 0, len(ignoredPatterns))
	for _, pattern := range ignoredPatterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			log.Error(err)
			continue
		}
		regexes = append(regexes, regex)
	}
	return regexes
}

func isLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func isPiped() bool {
	fi, _ := os.Stdout.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}
