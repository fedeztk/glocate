package engine

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/log"

	"github.com/fedeztk/glocate/config"
)

// Search searches for a pattern in the lz4 database
func Search(reg string, conf config.Config, smartcase bool) {
	if smartcase && isLowerCase(reg) {
		reg = "(?i)" + reg
	}

	re, err := regexp.Compile(reg)
	if err != nil {
		log.Fatal(err)
	}

	for path := range decompressPipe() {
		if re.MatchString(path) {
			fmt.Println(path)
		}
	}
}
