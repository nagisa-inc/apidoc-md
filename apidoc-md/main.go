package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nagisa-inc/apidoc-md"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = ""
	app.Description = "Convert apiDoc (http://apidocjs.com) documentation to markdown"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "input, i",
			Usage: "input apidoc `DIR` path",
		},
		&cli.StringFlag{
			Name:  "output, o",
			Usage: "output markdown `FILE` path",
		},
	}
	app.Action = func(c *cli.Context) error {
		// input
		inOpt := c.String("i")
		if inOpt == "" {
			return fmt.Errorf("--input option is required")
		}
		inDirPath, err := filepath.Abs(inOpt)
		if err != nil {
			return err
		}
		inInfo, err := os.Stat(inDirPath)
		if err != nil {
			return err
		}
		if !inInfo.IsDir() {
			return fmt.Errorf("input path is not directory")
		}

		// output
		outOpt := c.String("o")
		if outOpt == "" {
			return fmt.Errorf("--output option is required")
		}
		outPath, err := filepath.Abs(outOpt)
		if err != nil {
			return err
		}
		outInfo, err := os.Stat(outPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			if outInfo.IsDir() {
				return fmt.Errorf("output path is directory")
			}
		}

		return apidocmd.Convert(inDirPath, outPath)
	}

	// run
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("[error] %s\n", err)
		os.Exit(1)
	}
}
