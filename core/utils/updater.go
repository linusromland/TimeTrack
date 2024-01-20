package utils

import (
	"TimeTrack/core/database"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func parseVersion(version string) int {
	if version == "vdev" {
		// Max integer value.
		return 2147483647
	}

	versionRegex := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	if !versionRegex.MatchString(version) {
		fmt.Println("Invalid version format. Version is: ", version)
		return 0
	}

	// Remove the "v" from the version.
	version = strings.ReplaceAll(version, "v", "")

	// Remove all dots from the version.
	version = strings.ReplaceAll(version, ".", "")

	// Convert the version to an integer.
	versionInt := 0
	fmt.Sscanf(version, "%d", &versionInt)

	return versionInt
}

func githubRepo() (string, string) {
	githubUser := "linusromland"
	githubRepo := "TimeTrack"

	return githubUser, githubRepo
}

func CheckForUpdate(version string) (string, error) {
	db, err := database.OpenDB()
	if err != nil {
		database.CloseDB(db)
		return "", err
	}

	// Parse the version number to an integer.
	currentVersion := parseVersion("v" + version)

	githubUser, githubRepo := githubRepo()

	// Get the latest release from Github.
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubUser, githubRepo))
	if err != nil {
		database.CloseDB(db)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		database.CloseDB(db)
		return "", err
	}

	// Create a struct for the response body.
	var release struct {
		TagName string `json:"tag_name"`
	}

	// Unmarshal the response body into the struct.
	err = json.Unmarshal(body, &release)
	if err != nil {
		database.CloseDB(db)
		return "", err
	}

	tagName := release.TagName

	// Parse the version number to an integer.
	latestVersion := parseVersion(tagName)

	// If the latest version is newer than the current version, return.
	if latestVersion > currentVersion {
		database.CloseDB(db)
		return tagName, nil
	}

	database.CloseDB(db)
	return "", nil
}

func UpdateVersion(version string) error {
	// Get the OS type.
	osType := runtime.GOOS
	osType = strings.Title(osType)

	// Define the command to run on Windows and macOS, and a default message for other systems.
	var cmd *exec.Cmd
	switch osType {
	case "Windows":
		cmd = exec.Command("powershell", "-Command", "Invoke-WebRequest", "-Uri", "https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.bat", "-OutFile", "install.bat", ";", ".\\install.bat")
	case "Darwin":
		cmd = exec.Command("bash", "-c", "curl -sSL https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.sh | bash")
	default:
		fmt.Println("Sorry, no auto-update available for your system.")
		return nil
	}

	// Set cmd.Stdout and cmd.Stderr to os.Stdout and os.Stderr to see the command's output.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command.
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
