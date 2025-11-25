# kyctl

[![Go Test](https://github.com/AdebayoEmmanuel/kyctl/actions/workflows/go.yml/badge.svg)](https://github.com/AdebayoEmmanuel/kyctl/actions/workflows/go.yml)

`kyctl` is a command-line interface (CLI) tool for interacting with Kyverno policies and policy reports within a Kubernetes cluster. It allows you to easily list policies, inspect their details, and view reports without needing to write complex `kubectl` commands.

The tool is designed to be a lightweight wrapper around `kubectl`, simplifying common Kyverno-related queries into easy-to-remember commands.

## Features

-   Check the current Kubernetes context.
-   List all Kyverno `ClusterPolicies`.
-   Get details for a specific `ClusterPolicy`.
-   View all Kyverno policy reports (`PolicyReports` and `ClusterPolicyReports`).
-   Filter policy reports by status (`pass`, `fail`, etc.).
-   View reports for a specific policy.

## Prerequisites

-   Go (version 1.24 or later) for building from source.
-   `kubectl` installed and available in your system's `PATH`.
-   A configured `kubeconfig` file (e.g., at `~/.kube/config`) pointing to a Kubernetes cluster.
-   Kyverno installed in the cluster to provide the `clusterpolicies` and `policyreports` resources.

## Building from Source

1.  Clone the repository:
    ```sh
    git clone https://github.com/AdebayoEmmanuel/kyctl.git
    cd kyctl
    ```

2.  Build the binary:
    ```sh
    go build -o kyctl main.go
    ```

3.  (Optional) Move the binary to a directory in your `PATH` for system-wide access:
    ```sh
    sudo mv kyctl /usr/local/bin/
    ```

## Testing

The project is fully unit-tested. The tests use an interface-based mocking approach to simulate `kubectl` calls, allowing the test suite to run without needing a live Kubernetes cluster.

To run the tests, execute the following command from the project root:

```sh
go test ./... -v
```

## Usage

`kyctl` provides several commands to interact with Kyverno resources.

### `context`

Check the current Kubernetes context and cluster information.

```sh
kyctl context
```

### `policies`

List and inspect Kyverno policies.

```sh
# List all available ClusterPolicies
kyctl policies --all

# Get details for a specific policy
kyctl policies --policy <policy-name>
```

### `reports`

View Kyverno policy reports.

```sh
# View all policy reports
kyctl reports --all

# Filter reports by status (e.g., "fail")
kyctl reports --all --filter fail

# View reports for a specific policy
kyctl reports --policy <policy-name>
```

### `version`

Print the version of `kyctl`.

```sh
kyctl version
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.