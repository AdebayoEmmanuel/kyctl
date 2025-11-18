package k8s

import (
    "fmt"
    "os"
    "path/filepath"

    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

// GetKubernetesRESTClient returns a REST client that can make raw HTTP requests to the Kubernetes API.
// This is a minimal approach to avoid pulling in the full client-go library.
func GetKubernetesRESTClient() (*rest.RESTClient, error) {
    // Get the kubeconfig file path
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

    // Use the kubeconfig file to create a config
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        return nil, err
    }

    // Set the API group and version for Kyverno
    config.GroupVersion = &schema.GroupVersion{Group: "kyverno.io", Version: "v1"}
    config.APIPath = "/apis"
    
    // Create a new RESTClient
    restClient, err := rest.RESTClientFor(config)
    if err != nil {
        return nil, err
    }

    return restClient, nil
}

// GetCurrentContext returns the current Kubernetes context and cluster name
func GetCurrentContext() (string, string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", "", err
    }
    kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

    config, err := clientcmd.LoadFromFile(kubeconfigPath)
    if err != nil {
        return "", "", err
    }

    currentContext := config.CurrentContext
    context, exists := config.Contexts[currentContext]
    if !exists {
        return "", "", fmt.Errorf("context %s not found", currentContext)
    }

    return currentContext, context.Cluster, nil
}
