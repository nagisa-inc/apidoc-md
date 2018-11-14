package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Value: "",
			Usage: "input `DIR` path",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "output `FILE` path",
		},
	}
	app.Action = func(c *cli.Context) error {
		// input
		inputDirPath, err := filepath.Abs(c.String("i"))
		if err != nil {
			return err
		}
		info, err := os.Stat(inputDirPath)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("input path is not directory")
		}

		// output
		outputPath, err := filepath.Abs(c.String("o"))
		if err != nil {
			return err
		}

		fmt.Println(inputDirPath)
		fmt.Println(outputPath)

		//return exec(inputDirPath, outputPath)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
