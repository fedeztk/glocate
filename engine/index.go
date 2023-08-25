package engine

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/fedeztk/glocate/config"
	"github.com/opencoff/go-walk"
)

// Index indexes the directories specified in the config struct
func Index(conf config.Config) {
	if len(conf.Directories) == 0 {
		log.Fatal("no directories to index")
	}

	directories := cleanDirectories(conf.Directories)
	regexes := buildRegexes(conf.IgnoredPatterns)

	wg := sync.WaitGroup{}
	wg.Add(1)

	filePipe, errChan := walk.Walk(
		directories,
		&walk.Options{
			FollowSymlinks: !conf.IgnoreSymlinks,
			// Excludes:       conf.IgnoredPatterns,
			Type: walk.FILE | walk.DIR | walk.SYMLINK,
			Filter: func(filePath string, info os.FileInfo) bool {
				if conf.IgnoreHidden {
					hidden, err := isHidden(filepath.Base(filePath))
					if err != nil {
						log.Error(err)
					}
					if hidden {
						log.Debugf("Ignoring hidden file: %s", filePath)
						return true
					}
				}

				for _, re := range regexes {
					if re.MatchString(filePath) {
						log.Debugf("Regex macthed: %s for file %s", re, filePath)
						return true
					}
				}

				return false
			},
		},
	)

	// read errors from the channel
	go func() {
		for e := range errChan {
			// log error only in debug mode since most of the errors are permission denied
			log.Debugf("Error while walking: %s", e)
		}
		wg.Done()
	}()

	// create the db file
	f, err := os.Create(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// compress all entries from the channel to a lz4 file
	err = compressPipe(filePipe, f)
	if err != nil {
		log.Fatal(err)
	}
}
