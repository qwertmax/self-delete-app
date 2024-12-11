package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	maxAttempts = 3
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	attempts := 0
	for attempts < maxAttempts {
		fmt.Print("Enter password: ")
		var password string
		fmt.Scanln(&password)

		if password == "secure123" {
			fmt.Println("Access granted. Application running...")
			// Continue with your application logic here
			return
		}

		attempts++
		remaining := maxAttempts - attempts
		if remaining > 0 {
			fmt.Printf("Incorrect password. %d attempts remaining.\n", remaining)
		} else {
			fmt.Println("Access denied. Self-destructing...")

			// Wait a moment before deletion
			time.Sleep(2 * time.Second)

			// Self-deletion process
			deleteFile(execPath)
			os.Exit(1)
		}
	}
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func deleteFile(path string) {
	// On Windows, we need to handle file deletion differently
	if runtime.GOOS == "windows" {
		// First remove read-only attribute if present
		os.Chmod(path, 0666)

		// Create a batch file to delete the executable
		batchContent := fmt.Sprintf(`
@echo off
:retry
del "%s"
if exist "%s" goto retry
del "%%~f0"
`, path, path)

		batchFile := path + ".bat"
		os.WriteFile(batchFile, []byte(batchContent), 0666)

		// Execute the batch file
		os.StartProcess(batchFile, []string{}, &os.ProcAttr{})
	} else {
		// For Unix-like systems, we can delete directly
		os.Remove(path)
	}
}
