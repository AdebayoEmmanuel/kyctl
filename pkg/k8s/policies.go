package k8s

import (
    "encoding/json"
    "fmt"
)

// SimplePolicy matches the interface expected by the cmd package
type SimplePolicy struct {
    Name                    string
    ValidationFailureAction string
    Rules                   []SimpleRule
}

type SimpleRule struct {
    Name string
}

// Internal structs for JSON parsing
type k8sPolicyList struct {
    Items []rawPolicy `json:"items"`
}

type rawPolicy struct {
    Metadata struct {
        Name string `json:"name"`
    } `json:"metadata"`
    Spec struct {
        ValidationFailureAction string `json:"validationFailureAction"`
        Rules                   []struct {
            Name string `json:"name"`
        } `json:"rules"`
    } `json:"spec"`
}

// GetAllPolicies retrieves all Kyverno policies using kubectl
func GetAllPolicies() ([]SimplePolicy, error) {
    output, err := Executor.Run("get", "clusterpolicies", "-o", "json")
    if err != nil {
        return nil, err
    }

    var list k8sPolicyList
    if err := json.Unmarshal(output, &list); err != nil {
        return nil, fmt.Errorf("failed to parse policies json: %w", err)
    }

    var policies []SimplePolicy
    for _, item := range list.Items {
        policies = append(policies, convertToSimplePolicy(item))
    }
    return policies, nil
}

// GetPolicy retrieves a specific Kyverno policy by name using kubectl
func GetPolicy(name string) (*SimplePolicy, error) {
    output, err := Executor.Run("get", "clusterpolicy", name, "-o", "json")
    if err != nil {
        return nil, err
    }

    var item rawPolicy
    if err := json.Unmarshal(output, &item); err != nil {
        return nil, fmt.Errorf("failed to parse policy json: %w", err)
    }

    p := convertToSimplePolicy(item)
    return &p, nil
}

func convertToSimplePolicy(p rawPolicy) SimplePolicy {
    sp := SimplePolicy{
        Name:                    p.Metadata.Name,
        ValidationFailureAction: p.Spec.ValidationFailureAction,
    }
    for _, r := range p.Spec.Rules {
        sp.Rules = append(sp.Rules, SimpleRule{Name: r.Name})
    }
    return sp
}