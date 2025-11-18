package k8s

import (
    "context"
    "fmt"
    "path/filepath"
    "strings"

    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "sigs.k8s.io/yaml"
)

// PolicyReport represents a policy report entry with only the fields we need.
type PolicyReport struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Results           []PolicyReportResult `json:"results"`
}

// PolicyReportList is a list of PolicyReport objects.
type PolicyReportList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []PolicyReport `json:"items"`
}

// PolicyReportResult contains the result of a policy application.
type PolicyReportResult struct {
    Policy string           `json:"policy"`
    Result PolicyResultEntry `json:"result"`
    Scope  ScopeEntry        `json:"scope"`
}

// PolicyResultEntry contains the details of the result.
type PolicyResultEntry struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

// ScopeEntry contains the scope of the result.
type ScopeEntry struct {
    Kind      string `json:"kind"`
    Name      string `json:"name"`
    Namespace string `json:"namespace"`
}

// Simplified struct for our CLI output
type SimplePolicyReport struct {
    Policy    string
    Namespace string
    Resource  string
    Status    string
    Message   string
}

// GetAllPolicyReports retrieves all policy reports from the cluster.
func GetAllPolicyReports(filter string) ([]SimplePolicyReport, error) {
    var allReports []SimplePolicyReport

    // Get ClusterPolicyReports
    clusterReports, err := getReports("clusterpolicyreports", "")
    if err != nil {
        return nil, fmt.Errorf("failed to get cluster policy reports: %v", err)
    }
    allReports = append(allReports, clusterReports...)

    // Get namespaced PolicyReports
    namespacedReports, err := getReports("policyreports", "")
    if err != nil {
        return nil, fmt.Errorf("failed to get namespaced policy reports: %v", err)
    }
    allReports = append(allReports, namespacedReports...)

    // Apply filter if needed
    if filter != "" {
        var filteredReports []SimplePolicyReport
        for _, report := range allReports {
            if strings.ToLower(report.Status) == strings.ToLower(filter) {
                filteredReports = append(filteredReports, report)
            }
        }
        return filteredReports, nil
    }

    return allReports, nil
}

// GetPolicyResources retrieves all resources affected by a specific policy.
func GetPolicyResources(policyName string) ([]SimplePolicyReport, error) {
    allReports, err := GetAllPolicyReports("")
    if err != nil {
        return nil, err
    }

    var result []SimplePolicyReport
    for _, report := range allReports {
        if report.Policy == policyName {
            result = append(result, report)
        }
    }

    return result, nil
}

// Helper function to get reports from a specific API endpoint
func getReports(resource, namespace string) ([]SimplePolicyReport, error) {
    restClient, err := getPolicyReportRESTClient()
    if err != nil {
        return nil, fmt.Errorf("failed to create REST client: %v", err)
    }

    request := restClient.Get().Resource(resource)
    if namespace != "" {
        request = request.Namespace(namespace)
    }

    result := PolicyReportList{}
    err = request.Do(context.TODO()).Into(&result)
    if err != nil {
        return nil, fmt.Errorf("failed to get %s: %v", resource, err)
    }

    return convertToSimpleReports(result.Items, namespace), nil
}

// Helper function to create a REST client for Policy Reports
func getPolicyReportRESTClient() (*rest.RESTClient, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        return nil, err
    }

    // Configure for PolicyReport API access
    config.GroupVersion = &schema.GroupVersion{Group: "wgpolicyk8s.io", Version: "v1alpha2"}
    config.APIPath = "/apis"
    config.NegotiatedSerializer = yaml.NewYAMLSerializer(yaml.DefaultMetaFactory, nil, nil)

    return rest.RESTClientFor(config)
}

// Helper function to convert full PolicyReport objects to our Simple struct
func convertToSimpleReports(reports []PolicyReport, defaultNamespace string) []SimplePolicyReport {
    var result []SimplePolicyReport
    for _, report := range reports {
        for _, r := range report.Results {
            ns := r.Scope.Namespace
            if ns == "" {
                ns = defaultNamespace
            }
            sr := SimplePolicyReport{
                Policy:    r.Policy,
                Namespace: ns,
                Resource:  fmt.Sprintf("%s/%s", r.Scope.Kind, r.Scope.Name),
                Status:    r.Result.Status,
                Message:   r.Result.Message,
            }
            result = append(result, sr)
        }
    }
    return result
}
