package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nagisa-inc/apidoc-md"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short:   "Convert apiDoc (http://apidocjs.com) documentation to markdown",
	Version: "1.0.0",
	RunE:    run,
}

var (
	input  string
	output string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input apidoc `DIR` path")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output markdown `FILE` path")
}

func run(_ *cobra.Command, _ []string) error {
	// input
	if input == "" {
		return fmt.Errorf("--input option is required")
	}
	inDirPath, err := filepath.Abs(input)
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
	if output == "" {
		return fmt.Errorf("--output option is required")
	}
	outPath, err := filepath.Abs(output)
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("[error] %s\n", err)
		os.Exit(1)
	}
}
