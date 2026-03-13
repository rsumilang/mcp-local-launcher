package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// runCommand runs cmd and returns a formatted error if it fails (stderr or err message).
func runCommand(cmd *exec.Cmd, errPrefix string) error {
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	msg := strings.TrimSpace(string(out))
	if msg == "" {
		msg = err.Error()
	}
	return fmt.Errorf("%s: %s", errPrefix, msg)
}

// openApp opens a local application by name using the appropriate OS command.
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
	default:
		cmd = exec.Command(appName)
	}
	if err := runCommand(cmd, fmt.Sprintf("failed to open app %q", appName)); err != nil {
		return "", err
	}
	return fmt.Sprintf("Opened application: %s", appName), nil
}

// openURL opens a URL in the user's default browser using the appropriate OS command.
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
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := runCommand(cmd, fmt.Sprintf("failed to open URL %q", url)); err != nil {
		return "", err
	}
	return fmt.Sprintf("Opened URL: %s", url), nil
}

// openPath opens a file or folder with the default application.
func openPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	path = expandPath(path)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("powershell.exe", "Start-Process", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	if err := runCommand(cmd, fmt.Sprintf("failed to open path %q", path)); err != nil {
		return "", err
	}
	return fmt.Sprintf("Opened path: %s", path), nil
}

// revealInFinder reveals a file or folder in the system file manager (Finder, Explorer, etc.).
func revealInFinder(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	path = expandPath(path)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-R", path)
	case "windows":
		// Quote path so Explorer handles paths with spaces (e.g. C:\Program Files\foo).
		cmd = exec.Command("explorer.exe", "/select,"+`"`+path+`"`)
	default:
		cmd = exec.Command("xdg-open", filepath.Dir(path))
	}
	if err := runCommand(cmd, fmt.Sprintf("failed to reveal path %q", path)); err != nil {
		return "", err
	}
	return fmt.Sprintf("Revealed in file manager: %s", path), nil
}

// openWithApp opens a URL or file with a specific application.
func openWithApp(appName, target string) (string, error) {
	if strings.TrimSpace(appName) == "" {
		return "", fmt.Errorf("app_name must not be empty")
	}
	if strings.TrimSpace(target) == "" {
		return "", fmt.Errorf("target must not be empty (path or URL)")
	}
	target = expandPath(target)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-a", appName, target)
	case "windows":
		cmd = exec.Command("powershell.exe", "Start-Process", appName, "-ArgumentList", target)
	default:
		cmd = exec.Command(appName, target)
	}
	if err := runCommand(cmd, fmt.Sprintf("failed to open %q with %q", target, appName)); err != nil {
		return "", err
	}
	return fmt.Sprintf("Opened %s with %s", target, appName), nil
}

// expandPath converts a path with ~ to the user's home directory.
func expandPath(path string) string {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	return filepath.Join(home, path[2:])
}
