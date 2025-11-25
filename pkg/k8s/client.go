package k8s

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
)

// RunKubectlCommand executes a kubectl command and returns the output
func RunKubectlCommand(args ...string) ([]byte, error) {
    // Check if kubectl is installed
    _, err := exec.LookPath("kubectl")
    if err != nil {
        return nil, fmt.Errorf("kubectl not found in PATH: %w", err)
    }

    cmd := exec.Command("kubectl", args...)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    
    err = cmd.Run()
    if err != nil {
        return nil, fmt.Errorf("kubectl error: %s: %w", stderr.String(), err)
    }
    return out.Bytes(), nil
}

// GetCurrentContext returns the current Kubernetes context and cluster name
func GetCurrentContext() (string, string, error) {
    // Get context
    ctxOut, err := RunKubectlCommand("config", "current-context")
    if err != nil {
        return "", "", err
    }
    context := strings.TrimSpace(string(ctxOut))

    // Get cluster name
    // We use jsonpath to extract the cluster name for the current context
    clusterOut, err := RunKubectlCommand("config", "view", "--minify", "-o", "jsonpath={.clusters[0].name}")
    if err != nil {
        // Fallback if we can't get the cluster name, just return context
        return context, "unknown", nil
    }
    cluster := strings.TrimSpace(string(clusterOut))

    return context, cluster, nil
}