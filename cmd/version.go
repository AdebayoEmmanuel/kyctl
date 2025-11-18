package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var (
    versionCmd = &cobra.Command{
        Use:   "version",
        Short: "Print the version number of Kyctl",
        Long:  `All software has versions. This is Kyctl's`,
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Kyctl v0.1.0-minimal")
        },
    }
)
