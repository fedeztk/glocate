package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/fedeztk/glocate/config"
	"github.com/fedeztk/glocate/engine"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var appName = "glocate"
var configFileName = config.GetConfigFileName()
var conf = config.Config{}
var debugLevel int

var app = &cli.App{
	Flags: flags,
	Authors: []*cli.Author{
		{
			Email: "fedeztk@tutanota.com",
			Name:  "Federico Serra",
		},
	},
	Action: func(cCtx *cli.Context) error {
		log.Debug(fmt.Sprintf("starting with debug level: %d", debugLevel))
		if cCtx.Bool("index") {
			engine.Index(conf)
		} else {
			if cCtx.NArg() != 1 {
				return fmt.Errorf("wrong number of arguments, expected 1, got %d", cCtx.NArg())
			}
			engine.Search(cCtx.Args().First(), conf)
		}
		return nil
	},
	Name:  appName,
	Usage: "a cli tool for searching files in your filesystem",
	Description: `glocate is a fast and lightweight alternative to locate(1) and updatedb(8) written in Go.
Uses smartcase by default, i.e. case-sensitive if the pattern contains uppercase characters, case-insensitive otherwise`,
	UseShortOptionHandling: true,
	EnableBashCompletion:   true,
	Before: func(cCtx *cli.Context) error {
		log.SetLevel(log.ErrorLevel)
		return altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))(cCtx)
	},
	HideHelpCommand: true,
}

var flags = []cli.Flag{
	// cli flags
	&cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		Usage:   "show verbose output, can be used multiple times. Increase value to show more info, by default only errors are shown",
		Count:   &debugLevel,
		Action: func(cCtx *cli.Context, b bool) error {
			log.SetLevel(log.FatalLevel - log.Level(debugLevel))
			return nil
		},
		Category: "cli only",
		Value:    false,
	},
	&cli.BoolFlag{
		Name:     "index",
		Aliases:  []string{"i"},
		Usage:    "index directories",
		Category: "cli only",
	},
	&cli.StringFlag{
		// TODO: implement this
		Name:     "config",
		Value:    configFileName,
		Usage:    "config file to use",
		Category: "cli only",
	},
	// cli/env/config file flags
	altsrc.NewBoolFlag(
		&cli.BoolFlag{
			Name:        "ignore-symlinks",
			Usage:       "ignore symlinks",
			Aliases:     []string{"ignoreSymlinks"},
			Category:    "cli and config file",
			Destination: &conf.IgnoreSymlinks,
			EnvVars:     []string{"GLOCATE_IGNORE_SYMLINKS"},
			Value:       true,
		},
	),
	altsrc.NewBoolFlag(
		&cli.BoolFlag{
			Name:        "gitignore",
			Usage:       "ignore files and directories specified in .gitignore",
			Aliases:     []string{"gitIgnore"},
			Category:    "cli and config file",
			Destination: &conf.GitIgnore,
			EnvVars:     []string{"GLOCATE_GITIGNORE"},
		},
	),
	altsrc.NewBoolFlag(
		&cli.BoolFlag{
			Name:        "ignorehidden",
			Usage:       "ignore hidden files and directories",
			Aliases:     []string{"ignoreHidden"},
			Category:    "cli and config file",
			Destination: &conf.IgnoreHidden,
			EnvVars:     []string{"GLOCATE_IGNORE_HIDDEN"},
		},
	),
	altsrc.NewStringSliceFlag(
		&cli.StringSliceFlag{
			Name:     "directories",
			Usage:    "directories to index",
			Category: "cli and config file",
			Action: func(cCtx *cli.Context, s []string) error {
				conf.Directories = s
				return nil
			},
			EnvVars: []string{"GLOCATE_DIRECTORIES"},
		},
	),
	altsrc.NewStringSliceFlag(
		&cli.StringSliceFlag{
			Name:     "ignored-patterns",
			Usage:    "patterns to ignore",
			Category: "cli and config file",
			Aliases:  []string{"ignoredPatterns"},
			Action: func(cCtx *cli.Context, s []string) error {
				conf.IgnoredPatterns = s
				return nil
			},
			EnvVars: []string{"GLOCATE_IGNORED_PATTERNS"},
		},
	),
}

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
FLAGS PRECEDENCE:
	flags > env vars > config file > default value
`, cli.AppHelpTemplate)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
