package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/AdebayoEmmanuel/kyctl/pkg/k8s"
)

var contextCmd = &cobra.Command{
    Use:   "context",
    Short: "Check the current Kubernetes context",
    Long:  `Check the current Kubernetes context and cluster information`,
    Run: func(cmd *cobra.Command, args []string) {
        context, cluster, err := k8s.GetCurrentContext()
        if err != nil {
            fmt.Printf("Error getting current context: %v\n", err)
            fmt.Println("Please ensure you're authenticated to a Kubernetes cluster")
            return
        }
        
        fmt.Printf("Current context: %s\n", context)
        fmt.Printf("Cluster name: %s\n", cluster)
        fmt.Println("You're connected to a Kubernetes cluster")
    },
}
