package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"TimeTrack/src/database"
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

func CheckForUpdate(version string, skipChecks bool) (bool, error) {
	db, err := database.OpenDB()
	if err != nil {
		database.CloseDB(db)
		return false, err
	}

	LAST_UPDATE_CHECK_DB_KEY := "last_update_check"

	if !skipChecks {
		// Get the last time the user checked for updates.
		lastUpdateCheck := database.GetData(db, LAST_UPDATE_CHECK_DB_KEY)

		// Check if the last check is older than 6 hours.
		if lastUpdateCheck != "" {
			// Parse the last update check to a time.
			lastUpdateCheckTime, err := time.Parse("2006-01-02 15:04", lastUpdateCheck)
			if err != nil {
				return false, err
			}

			// Get the current time.
			currentTime := time.Now()

			// Calculate the difference between the last update check and the current time.
			diff := currentTime.Sub(lastUpdateCheckTime)

			// If the difference is less than 6 hours, return.
			if diff.Hours() < 6 {
				database.CloseDB(db)
				return false, err
			}
		}
	}

	// Parse the version number to an integer.
	currentVersion := parseVersion("v" + version)

	githubUser, githubRepo := githubRepo()

	// Get the latest release from Github.
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubUser, githubRepo))
	if err != nil {
		database.CloseDB(db)
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		database.CloseDB(db)
		return false, err
	}

	// Create a struct for the response body.
	var release struct {
		TagName string `json:"tag_name"`
	}

	// Unmarshal the response body into the struct.
	err = json.Unmarshal(body, &release)
	if err != nil {
		database.CloseDB(db)
		return false, err
	}

	// Set the last update check to the current time.
	database.SetData(db, LAST_UPDATE_CHECK_DB_KEY, time.Now().Format("2006-01-02 15:04"))

	tagName := release.TagName

	// Parse the version number to an integer.
	latestVersion := parseVersion(tagName)

	// If the latest version is bigger than the current version, return the latest version.
	if latestVersion > currentVersion {
		SKIP_DB_KEY := "skip_update"

		if !skipChecks {
			// Check if the user has skipped this update.
			skipUpdate := database.GetData(db, SKIP_DB_KEY)

			// If the user has skipped this update, return.
			if skipUpdate == tagName {
				database.CloseDB(db)
				return true, nil
			}
		}

		fmt.Printf("There is a new version of TimeTrack available: %s\n", tagName)

		if Confirm("Do you want to update?") {
			err = UpdateVersion(tagName)
			if err != nil {
				database.CloseDB(db)
				return false, err
			}
		} else {
			fmt.Println("Skipping this update, you can update later by running 'timetrack update'.")
			database.SetData(db, SKIP_DB_KEY, tagName)
		}

		database.CloseDB(db)
		return true, nil
	}

	database.CloseDB(db)
	return false, nil
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
