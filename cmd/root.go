package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anemon",
	Short: "a CV genrator",
	Long:  `This CLI tool, written in Go, automates the generation of customized CVs from Markdown files based on a specified configuration. It parses CV sections in
    multiple languages, prioritizes key skills or features as defined in an output.yml file, and outputs LaTeX files for each CV version, ready for compilation.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
