# TimeTrack Changelog

## [Unreleased]

- Fixed SSL/TLS certificate verification in Docker container for secure external API communication
- Updated CLI to integrate with new API-centric architecture
- Removed debug logging for cleaner production output

## 0.5.3 (2025-12-18)

- Sets version in the API service using build flags in the Dockerfile.

## 0.5.2 (2025-12-18)

- Sets GIN_MODE to release and PORT to 8080 in the Dockerfile for the API service.

## 0.5.1 (2025-12-18)

- Fixed issue with docker image not starting.

## 0.5.0 (2025-10-25)

- Major architectural overhaul: transitioned from a CLI-centric design to a clean, service-based architecture.
- Introduced a new internal API built with Go and the Gin web framework. The CLI now communicates with this API for all persistence and business logic.
- Replaced Google Calendar storage with a MongoDB backend, improving reliability, scalability, and data control.
- Added automatic time reporting to Jira, enabling better integration and workflow automation.
- Established a foundation for future web application support using the same API and account system as the CLI.

## 0.4.0 (2024-01-19)

- Implemented overtime/undertime calculation to `list` command.

## 0.3.1 (2023-12-03)

-   Fixed bug where `list` command would not use the unit specified by the user.
-   Added `--next` and `--last` flags to `list` command to show the next or last unit of time.

## 0.3.0 (2023-12-03)

-   Added automatic check for updates on startup.
-   Added `update` command to manually check and update.

## 0.2.0 (2023-11-01)

-   Added `change` command to change the current task.
-   Added `info` command to show information about the current task.
-   Fixed issue where timezones were not being handled correctly.
-   Improved logging.
-   Made description optional when ending a task.
-   Fixed type where `selectCalendar` was being called `selectDatabase`.

## 0.1.3 (2023-10-22)

-   Fixed bug where --version flag would always show version as "dev".
-   Added error message when no events are found.

## 0.1.2 (2023-10-22)

-   Updated Google Credentials.

## 0.1.1 (2023-10-22)

-   Added installation scripts for Windows & macOS.
-   Added installation instructions to README.md.
-   Added update instructions to README.md.

## 0.1.0 (2023-10-22)

Initial Release:

-   Introduced CLI tool "TimeTrack".
-   Added task creation support.
-   Implemented start and end task functionality.
-   Feature to list tasks by day, week, and month.
-   View time spent on each task.
-   Google Calendar integration for task tracking.
