package cmd

import (
    "github.com/Theo-Hafsaoui/Anemon/internal/adapters/input"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "anemon",
    Short: "A CV generator",
    Long:  `This CLI tool, automates the generation of customized CVs from Markdown files based on a specified configuration.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        threshold, err := cmd.Flags().GetInt("threshold")
        if err != nil { return err }

        root, err := os.Getwd()
        if err != nil { return err }

        if err != nil {
            return err
        }
        input.ChangeOverflowThreshold(threshold)

        generate, err := cmd.Flags().GetBool("generate")
        if err != nil { return err }

        if generate {
            return input.GenerateCVFromMarkDownToLatex(root)
        }

        info, err := cmd.Flags().GetBool("cvInfo")
        if err != nil { return err }
        if info{
            input.PrintAllCvs(root)
            return nil
        }

        return nil
    },
}

func Execute() {
    rootCmd.Flags().IntP("threshold", "t", 1, "Set the page overflow threshold (default 1)")
    rootCmd.Flags().BoolP("generate", "g", false, "Generate a CV")
    rootCmd.Flags().BoolP("cvInfo", "i", false, "Get all the info of all the cvs")
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
