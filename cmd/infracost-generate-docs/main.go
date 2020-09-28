package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/infracost/infracost/internal/docs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var logFormatter log.TextFormatter = log.TextFormatter{
	DisableTimestamp:       true,
	DisableLevelTruncation: true,
}

func customError(c *cli.Context, msg string) error {
	color.HiRed(fmt.Sprintf("%v\n", msg))
	_ = cli.ShowAppHelp(c)

	return fmt.Errorf("")
}

func getcwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Warn(err)
		cwd = ""
	}

	return cwd
}

func main() {

	app := &cli.App{
		Name:                 "infracost-geneate-docs",
		Usage:                "Generate infracost documentations",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "input",
				Usage:     "Path to docs templates",
				TakesFile: true,
			},
			&cli.StringFlag{
				Name:      "output",
				Usage:     "Path to output of docs",
				TakesFile: true,
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Log level (TRACE, DEBUG, INFO, WARN, ERROR)",
				Value: "WARN",
			},
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			return customError(c, err.Error())
		},
		Action: func(c *cli.Context) error {
			log.SetFormatter(&logFormatter)

			switch strings.ToUpper(c.String("log-level")) {
			case "TRACE":
				log.SetLevel(log.TraceLevel)
			case "DEBUG":
				log.SetLevel(log.DebugLevel)
			case "WARN":
				log.SetLevel(log.WarnLevel)
			case "ERROR":
				log.SetLevel(log.ErrorLevel)
			default:
				log.SetLevel(log.InfoLevel)
			}

			templatesPath := c.String("input")
			if templatesPath == "" {
				templatesPath = getcwd() + "/docs/templates"
			}

			outputPath := c.String("output")
			if outputPath == "" {
				outputPath = getcwd() + "/docs/generated"
			}

			err := docs.GenerateDocs(templatesPath, outputPath)
			if err != nil {
				return errors.Wrap(err, "")
			}

			return nil

		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}