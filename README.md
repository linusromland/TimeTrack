<img src="docs/logo.png"  width="100" height="100" align="center"/>

# TimeTrack CLI

TimeTrack is a Command Line Interface (CLI) application designed for efficient time tracking. With a simple command, you can add and start tasks. Once you're done, TimeTrack automatically syncs your tasks with Google Calendar, offering a visual representation of your time spent on each task.

## Features

-   **Task Management**: Easily add, start, and stop tasks.
-   **Google Calendar Integration**: Visualize your tasks and time spent using Google Calendar.
-   **Analytics**: List all tasks to see time spent per task, day, week, month, and more.

## Installation

macOS:

Open your terminal.

Run the following command to download and install TimeTrack:

```bash
curl -sSL https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.sh | bash
```

Windows:

For users with PowerShell (Windows 10 and newer):

1. Open PowerShell as Administrator.

2. Run the following command:

```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.bat" -OutFile "install.bat"; .\install.bat
```

Linux:

For Linux users, manual installation is required:

1. Navigate to the releases page of the TimeTrack repository.

2. Download the latest .tar.gz archive suitable for your OS and architecture.

3. Extract the archive using the following command:

```bash
tar -xzf TimeTrack\_<YOUR_OS_AND_ARCH>.tar.gz
```

Move the TimeTrack binary to a directory in your PATH, typically `/usr/local/bin:`

```bash
sudo mv TimeTrack /usr/local/bin/
```

### Update

To update TimeTrack to the latest version, simply re-run the installation command for your respective operating system.

## Collaborating

Contributions are what make the open source community such an amazing place to be, learn, inspire, and create. Any contributions you make are **greatly appreciated**. For major changes, please open an issue first to discuss what you would like to change.

To contribute to TimeTrack, please fork the repository, create a new branch, and submit a pull request.

## License

TimeTrack is released under the [MIT License](https://choosealicense.com/licenses/mit/).
