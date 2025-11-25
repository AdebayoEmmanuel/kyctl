package k8s

import (
    "encoding/json"
    "fmt"
    "strings"
)

// SimplePolicyReport for CLI output
type SimplePolicyReport struct {
    Policy    string
    Namespace string
    Resource  string
    Status    string
    Message   string
}

// Internal structs for JSON parsing
type k8sReportList struct {
    Items []rawReport `json:"items"`
}

type rawReport struct {
    Metadata struct {
        Name      string `json:"name"`
        Namespace string `json:"namespace"`
    } `json:"metadata"`
    Results []struct {
        Policy    string `json:"policy"`
        Result    string `json:"result"`
        Message   string `json:"message"`
        Resources []struct {
            Kind      string `json:"kind"`
            Name      string `json:"name"`
            Namespace string `json:"namespace"`
        } `json:"resources"`
    } `json:"results"`
}

// GetAllPolicyReports retrieves all policy reports from the cluster
func GetAllPolicyReports(filter string) ([]SimplePolicyReport, error) {
    // Get both PolicyReports and ClusterPolicyReports
    output, err := RunKubectlCommand("get", "policyreports,clusterpolicyreports", "-A", "-o", "json")
    if err != nil {
        return nil, fmt.Errorf("failed to get reports (ensure Kyverno is installed): %w", err)
    }

    var list k8sReportList
    if err := json.Unmarshal(output, &list); err != nil {
        return nil, fmt.Errorf("failed to parse reports json: %w", err)
    }

    var reports []SimplePolicyReport
    for _, item := range list.Items {
        for _, res := range item.Results {
            // Apply filter if provided
            if filter != "" && !strings.EqualFold(res.Result, filter) {
                continue
            }

            // Construct resource name
            resourceName := "unknown"
            if len(res.Resources) > 0 {
                r := res.Resources[0]
                resourceName = fmt.Sprintf("%s/%s", r.Kind, r.Name)
                if r.Namespace != "" {
                    resourceName = fmt.Sprintf("%s/%s", r.Namespace, resourceName)
                }
            }

            reports = append(reports, SimplePolicyReport{
                Policy:    res.Policy,
                Namespace: item.Metadata.Namespace,
                Resource:  resourceName,
                Status:    res.Result,
                Message:   res.Message,
            })
        }
    }
    return reports, nil
}

// GetPolicyResources retrieves all resources affected by a specific policy
func GetPolicyResources(policyName string) ([]SimplePolicyReport, error) {
    all, err := GetAllPolicyReports("")
    if err != nil {
        return nil, err
    }
    
    var filtered []SimplePolicyReport
    for _, r := range all {
        if r.Policy == policyName {
            filtered = append(filtered, r)
        }
    }
    return filtered, nil
}