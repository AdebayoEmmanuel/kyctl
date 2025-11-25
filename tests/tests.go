package tests

import (
    "errors"
    "testing"

    "github.com/AdebayoEmmanuel/kyctl/pkg/k8s"
)

// FakeExecutor is our test implementation of the CommandExecutor interface.
// It returns whatever output and error we tell it to, without running any real commands.
type FakeExecutor struct {
    Output []byte
    Err    error
}

// Run is the method that satisfies the interface.
func (f *FakeExecutor) Run(args ...string) ([]byte, error) {
    return f.Output, f.Err
}

// --- Test for client.go ---

func TestGetCurrentContext(t *testing.T) {
    // Setup: Replace the real executor with our fake one.
    k8s.Executor = &FakeExecutor{Output: []byte("minikube\n")}

    // Run the function we want to test.
    ctx, _, err := k8s.GetCurrentContext()

    // Assert: Check if the result is what we expect.
    if err != nil {
        t.Fatalf("Expected no error, but got: %v", err)
    }
    if ctx != "minikube" {
        t.Errorf("Expected context 'minikube', but got '%s'", ctx)
    }
}

// --- Tests for policies.go ---

func TestGetAllPolicies(t *testing.T) {
    fakeJSON := `
    {
        "items": [
            {
                "metadata": { "name": "disallow-latest-tag" },
                "spec": { "validationFailureAction": "enforce" }
            },
            {
                "metadata": { "name": "require-labels" },
                "spec": { "validationFailureAction": "audit" }
            }
        ]
    }`
    k8s.Executor = &FakeExecutor{Output: []byte(fakeJSON)}

    policies, err := k8s.GetAllPolicies()

    if err != nil {
        t.Fatalf("Expected no error, but got: %v", err)
    }
    if len(policies) != 2 {
        t.Fatalf("Expected 2 policies, but got %d", len(policies))
    }
    if policies[0].Name != "disallow-latest-tag" {
        t.Errorf("Expected first policy name 'disallow-latest-tag', but got '%s'", policies[0].Name)
    }
}

func TestGetPolicy(t *testing.T) {
    fakeJSON := `
    {
        "metadata": { "name": "require-labels" },
        "spec": {
            "validationFailureAction": "audit",
            "rules": [{ "name": "check-for-app-label" }]
        }
    }`
    k8s.Executor = &FakeExecutor{Output: []byte(fakeJSON)}

    policy, err := k8s.GetPolicy("require-labels")

    if err != nil {
        t.Fatalf("Expected no error, but got: %v", err)
    }
    if policy.Name != "require-labels" {
        t.Errorf("Expected policy name 'require-labels', but got '%s'", policy.Name)
    }
    if len(policy.Rules) != 1 {
        t.Errorf("Expected 1 rule, but got %d", len(policy.Rules))
    }
}

func TestGetPolicy_Error(t *testing.T) {
    k8s.Executor = &FakeExecutor{Err: errors.New("policy not found")}

    _, err := k8s.GetPolicy("non-existent-policy")

    if err == nil {
        t.Fatal("Expected an error, but got none")
    }
}

// --- Tests for reports.go ---

func TestGetAllPolicyReports(t *testing.T) {
    fakeJSON := `
    {
        "items": [
            {
                "metadata": { "name": "report-1", "namespace": "default" },
                "results": [
                    { "policy": "require-labels", "result": "fail" },
                    { "policy": "disallow-latest-tag", "result": "pass" }
                ]
            }
        ]
    }`
    k8s.Executor = &FakeExecutor{Output: []byte(fakeJSON)}

    // Test without filter
    reports, err := k8s.GetAllPolicyReports("")
    if err != nil {
        t.Fatalf("Expected no error, but got: %v", err)
    }
    if len(reports) != 2 {
        t.Fatalf("Expected 2 reports, but got %d", len(reports))
    }

    // Test with filter
    failedReports, err := k8s.GetAllPolicyReports("fail")
    if err != nil {
        t.Fatalf("Expected no error on filter, but got: %v", err)
    }
    if len(failedReports) != 1 {
        t.Fatalf("Expected 1 failed report, but got %d", len(failedReports))
    }
    if failedReports[0].Status != "fail" {
        t.Errorf("Expected status 'fail', but got '%s'", failedReports[0].Status)
    }
}

func TestGetPolicyResources(t *testing.T) {
    fakeJSON := `
    {
        "items": [
            {
                "metadata": { "name": "report-1", "namespace": "default" },
                "results": [
                    { "policy": "require-labels", "result": "fail" },
                    { "policy": "disallow-latest-tag", "result": "pass" },
                    { "policy": "require-labels", "result": "pass" }
                ]
            }
        ]
    }`
    k8s.Executor = &FakeExecutor{Output: []byte(fakeJSON)}

    resources, err := k8s.GetPolicyResources("require-labels")

    if err != nil {
        t.Fatalf("Expected no error, but got: %v", err)
    }
    if len(resources) != 2 {
        t.Fatalf("Expected 2 resources for policy 'require-labels', but got %d", len(resources))
    }
}