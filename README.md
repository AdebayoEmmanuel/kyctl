# kyctl

`kyctl` is a command-line interface (CLI) tool for interacting with Kyverno policies and policy reports within a Kubernetes cluster. It allows you to easily list policies, inspect their details, and view reports without needing to write complex `kubectl` commands.

## Features

-   Check the current Kubernetes context.
-   List all Kyverno `ClusterPolicies`.
-   Get details for a specific `ClusterPolicy`.
-   View all Kyverno policy reports (`PolicyReports` and `ClusterPolicyReports`).
-   Filter policy reports by status (`pass`, `fail`, etc.).
-   View reports for a specific policy.

## Prerequisites

-   Go (version 1.24 or later).
-   Access to a Kubernetes cluster.
-   `kubeconfig` file correctly set up (e.g., at `~/.kube/config`).
-   Kyverno and the Policy Reporter UI installed in the cluster.

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

3.  Move the binary to a directory in your `PATH`:
    ```sh
    sudo mv kyctl /usr/local/bin/
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