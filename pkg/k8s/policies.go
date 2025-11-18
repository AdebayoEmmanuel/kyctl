package k8s

import (
    "context"
    "fmt"
    "path/filepath"

    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "sigs.k8s.io/yaml"
)

// Policy represents a Kyverno policy with only the fields we need.
type Policy struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec              PolicySpec `json:"spec,omitempty"`
}

// PolicyList is a list of Policy objects.
type PolicyList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Policy `json:"items"`
}

// PolicySpec contains the specification for a Policy.
type PolicySpec struct {
    ValidationFailureAction string   `json:"validationFailureAction"`
    Rules                   []Rule   `json:"rules"`
}

// Rule represents a rule in a Kyverno policy.
type Rule struct {
    Name string `json:"name"`
}

// Simplified struct for our CLI output
type SimplePolicy struct {
    Name                    string
    ValidationFailureAction string
    Rules                   []SimpleRule
}

type SimpleRule struct {
    Name string
}

// GetAllPolicies retrieves all Kyverno policies using a raw REST client.
func GetAllPolicies() ([]SimplePolicy, error) {
    restClient, err := getGenericRESTClient()
    if err != nil {
        return nil, fmt.Errorf("failed to create REST client: %v", err)
    }

    // Define the API path for cluster policies
    result := PolicyList{}
    err = restClient.Get().Resource("clusterpolicies").Do(context.TODO()).Into(&result)
    if err != nil {
        return nil, fmt.Errorf("failed to get policies: %v", err)
    }

    return convertToSimplePolicies(result.Items), nil
}

// GetPolicy retrieves a specific Kyverno policy by name.
func GetPolicy(name string) (*SimplePolicy, error) {
    restClient, err := getGenericRESTClient()
    if err != nil {
        return nil, fmt.Errorf("failed to create REST client: %v", err)
    }

    result := Policy{}
    err = restClient.Get().Resource("clusterpolicies").Name(name).Do(context.TODO()).Into(&result)
    if err != nil {
        return nil, fmt.Errorf("failed to get policy %s: %v", name, err)
    }

    simplePolicies := convertToSimplePolicies([]Policy{result})
    if len(simplePolicies) == 0 {
        return nil, fmt.Errorf("policy %s not found", name)
    }
    return &simplePolicies[0], nil
}

// Helper function to create a generic REST client
func getGenericRESTClient() (*rest.RESTClient, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        return nil, err
    }

    // Configure for generic API access
    config.GroupVersion = &schema.GroupVersion{Group: "kyverno.io", Version: "v1"}
    config.APIPath = "/apis"
    config.NegotiatedSerializer = yaml.NewYAMLSerializer(yaml.DefaultMetaFactory, nil, nil)

    return rest.RESTClientFor(config)
}

// Helper function to convert full Policy objects to our Simple struct
func convertToSimplePolicies(policies []Policy) []SimplePolicy {
    var result []SimplePolicy
    for _, p := range policies {
        sp := SimplePolicy{
            Name:                    p.Name,
            ValidationFailureAction: p.Spec.ValidationFailureAction,
        }
        for _, r := range p.Spec.Rules {
            sp.Rules = append(sp.Rules, SimpleRule{Name: r.Name})
        }
        result = append(result, sp)
    }
    return result
}
