package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "kyctl",
    Short: "A CLI tool for checking Kyverno policies in Kubernetes",
    Long: `Kyctl is a CLI tool that helps you check and manage Kyverno policies 
in your Kubernetes cluster. It provides commands to list policies, check their 
status, and view policy reports.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(versionCmd)
    rootCmd.AddCommand(contextCmd)
    rootCmd.AddCommand(policiesCmd)
    rootCmd.AddCommand(reportsCmd)
}
