package output

import (
    "fmt"
)

// PrintSuccess prints a success message.
func PrintSuccess(message string) {
    fmt.Printf("[SUCCESS] %s\n", message)
}

// PrintError prints an error message.
func PrintError(message string) {
    fmt.Printf("[ERROR] %s\n", message)
}

// PrintInfo prints an info message.
func PrintInfo(message string) {
    fmt.Printf("[INFO] %s\n", message)
}
