package cmd

import (
    "fmt"

    "github.com/AdebayoEmmanuel/kyctl/pkg/k8s"
    "github.com/spf13/cobra"
)

var (
    reportsAll       bool
    reportPolicyName string
    filterStatus     string

    reportsCmd = &cobra.Command{
        Use:   "reports",
        Short: "View Kyverno policy reports",
        Long:  `View policy reports for all policies or a specific policy`,
        Run: func(cmd *cobra.Command, args []string) {
            if reportsAll {
                reports, err := k8s.GetAllPolicyReports(filterStatus)
                if err != nil {
                    fmt.Printf("Error getting policy reports: %v\n", err)
                    return
                }

                fmt.Println("Policy Reports:")
                fmt.Println("===============")
                for _, report := range reports {
                    fmt.Printf("Policy: %s\n", report.Policy)
                    // Fixed: Use .Resource directly (it includes namespace if applicable)
                    fmt.Printf("Resource: %s\n", report.Resource)
                    fmt.Printf("Status: %s\n", report.Status)
                    if report.Message != "" {
                        fmt.Printf("Message: %s\n", report.Message)
                    }
                    fmt.Println("---------------------------")
                }
            } else if reportPolicyName != "" {
                resources, err := k8s.GetPolicyResources(reportPolicyName)
                if err != nil {
                    fmt.Printf("Error getting policy resources: %v\n", err)
                    return
                }

                fmt.Printf("Resources affected by policy: %s\n", reportPolicyName)
                fmt.Println("=====================================")
                for _, resource := range resources {
                    // Fixed: Changed resource.Name to resource.Resource to fix build error
                    fmt.Printf("Resource: %s\n", resource.Resource)
                    fmt.Printf("Status: %s\n", resource.Status)
                    if resource.Message != "" {
                        fmt.Printf("Message: %s\n", resource.Message)
                    }
                    fmt.Println("---------------------------")
                }
            } else {
                fmt.Println("Please specify either --all or --policy <name>")
                _ = cmd.Help()
            }
        },
    }
)

func init() {
    reportsCmd.Flags().BoolVarP(&reportsAll, "all", "a", false, "Show all policy reports")
    reportsCmd.Flags().StringVarP(&reportPolicyName, "policy", "p", "", "Show report for a specific policy")
    reportsCmd.Flags().StringVarP(&filterStatus, "filter", "f", "", "Filter reports by status (pass/fail)")
}