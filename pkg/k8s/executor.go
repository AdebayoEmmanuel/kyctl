package k8s

// CommandExecutor defines how we run commands
type CommandExecutor interface {
    Run(args ...string) ([]byte, error)
}

// RealExecutor runs actual kubectl commands
type RealExecutor struct{}

func (r RealExecutor) Run(args ...string) ([]byte, error) {
    return RunKubectlCommand(args...)
}

// Executor is the variable we use to run commands. 
var Executor CommandExecutor = RealExecutor{}