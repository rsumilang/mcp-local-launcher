package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// openApp opens a local application by name using the appropriate OS command.
// It returns a human-readable result message and any error encountered.
func openApp(appName string) (string, error) {
	if strings.TrimSpace(appName) == "" {
		return "", fmt.Errorf("app_name must not be empty")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-a", appName)
	case "windows":
		cmd = exec.Command("powershell.exe", "Start-Process", appName)
	default: // linux and others
		cmd = exec.Command(appName)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("failed to open app %q: %s", appName, msg)
	}
	return fmt.Sprintf("Opened application: %s", appName), nil
}

// openURL opens a URL in the user's default browser using the appropriate OS command.
// It returns a human-readable result message and any error encountered.
func openURL(url string) (string, error) {
	if strings.TrimSpace(url) == "" {
		return "", fmt.Errorf("url must not be empty")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("powershell.exe", "Start-Process", url)
	default: // linux and others
		cmd = exec.Command("xdg-open", url)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("failed to open URL %q: %s", url, msg)
	}
	return fmt.Sprintf("Opened URL: %s", url), nil
}
