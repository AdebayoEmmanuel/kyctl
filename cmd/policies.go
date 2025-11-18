package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/AdebayoEmmanuel/kyctl/pkg/k8s"
)

var (
    listAll bool
    policyName string
    
    policiesCmd = &cobra.Command{
        Use:   "policies",
        Short: "List and inspect Kyverno policies",
        Long:  `List all Kyverno policies or get details of a specific policy`,
        Run: func(cmd *cobra.Command, args []string) {
            if listAll {
                policies, err := k8s.GetAllPolicies()
                if err != nil {
                    fmt.Printf("Error getting policies: %v\n", err)
                    return
                }
                
                fmt.Println("Available Kyverno Policies:")
                fmt.Println("===========================")
                for _, policy := range policies {
                    fmt.Printf("Name: %s\n", policy.Name)
                    fmt.Printf("Validation Failure Action: %s\n", policy.ValidationFailureAction)
                    fmt.Println("---------------------------")
                }
            } else if policyName != "" {
                policy, err := k8s.GetPolicy(policyName)
                if err != nil {
                    fmt.Printf("Error getting policy: %v\n", err)
                    return
                }
                
                fmt.Printf("Policy Details for: %s\n", policy.Name)
                fmt.Println("=============================")
                fmt.Printf("Validation Failure Action: %s\n", policy.ValidationFailureAction)
                fmt.Println("Rules:")
                for _, rule := range policy.Rules {
                    fmt.Printf("- %s\n", rule.Name)
                }
            } else {
                fmt.Println("Please specify either --all or --policy <name>")
                cmd.Help()
            }
        },
    }
)

func init() {
    policiesCmd.Flags().BoolVarP(&listAll, "all", "a", false, "List all policies")
    policiesCmd.Flags().StringVarP(&policyName, "policy", "p", "", "Get details of a specific policy")
}
