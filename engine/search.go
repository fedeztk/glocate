package engine

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/log"

	"github.com/fedeztk/glocate/config"
)

// Search searches for a pattern in the lz4 database
func Search(reg string, conf config.Config, smartcase bool, color bool) {
	if smartcase && isLowerCase(reg) {
		reg = "(?i)" + reg
	}

	shouldColor := color && !isPiped()

	re, err := regexp.Compile(reg)
	if err != nil {
		log.Fatal(err)
	}

	for path := range decompressPipe() {
		if re.MatchString(path) {
			// print colored matches if printing to a terminal
			if shouldColor {
				fmt.Println(re.ReplaceAllString(path, "\033[1;31m$0\033[0m"))
			} else {
				fmt.Println(path)
			}
		}
	}
}
