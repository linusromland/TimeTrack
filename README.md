<img src="docs/logo.png"  width="100" height="100" align="center"/>

# TimeTrack

TimeTrack is a time tracking tool designed. Originally a CLI-focused tool, it has evolved into a service-based platform built with Go, featuring an internal API, MongoDB for persistence, and Jira integration for automated time reporting.

## Installation

### Quick Install (Linux/macOS)

Install TimeTrack CLI with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.sh | bash
```

This will:
- Download the latest release for your platform
- Install to `/usr/local/bin/timetrack`
- Set up bash completion

#### Custom Installation Directory

```bash
# Install to user directory (no sudo required)
TIMETRACK_INSTALL_DIR="$HOME/.local/bin" curl -fsSL https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.sh | bash
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/linusromland/TimeTrack/releases)
2. Extract the archive
3. Move the binary to your PATH (e.g., `/usr/local/bin/timetrack`)
4. Make it executable: `chmod +x /usr/local/bin/timetrack`

### Bash Completion

If installed via the script, bash completion is automatically set up. For manual installations:

```bash
# Generate completion and add to your shell
timetrack --generate-bash-completion
```

## Features

-   **Task Management**: Easily add, start, and stop tasks.
-   **Analytics**: List all tasks to see time spent per task, day, week, month, and more.
-   **Bash Completion**: Auto-complete commands and options.

## Collaborating

Contributions are what make the open source community such an amazing place to be, learn, inspire, and create. Any contributions you make are **greatly appreciated**. For major changes, please open an issue first to discuss what you would like to change.

To contribute to TimeTrack, please fork the repository, create a new branch, and submit a pull request.

## License

TimeTrack is released under the [MIT License](https://choosealicense.com/licenses/mit/).
