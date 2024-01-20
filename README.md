<img src="docs/logo.png"  width="100" height="100" align="center"/>

# TimeTrack CLI

TimeTrack is a Command Line Interface (CLI) application designed for efficient time tracking. With a simple command, you can add and start tasks. Once you're done, TimeTrack automatically syncs your tasks with Google Calendar, offering a visual representation of your time spent on each task.

## Features

-   **Task Management**: Easily add, start, and stop tasks.
-   **Google Calendar Integration**: Visualize your tasks and time spent using Google Calendar.
-   **Analytics**: List all tasks to see time spent per task, day, week, month, and more.

## Installation

### Windows:

For Windows users, TimeTrack can be installed using the provided installer found on the [releases page](https://github.com/linusromland/TimeTrack/releases).
Simply download the latest .exe file and run it to install TimeTrack.

### Other Operating Systems:

Currently, TimeTrack is officially supported on Windows so support for other operating systems is limited. However, TimeTrack can be built from source for other operating systems. But success is not guaranteed.

To build from source follow these steps:

1. Clone the repository:

```bash
git clone git@github.com:linusromland/TimeTrack.git
```

2. Create a `.env`

```bash
cd TimeTrack
cp .env.example .env
```

Then fill in the required fields in the `.env` file. For more information on how to get the required credentials, see the [Google Calendar API documentation](https://developers.google.com/calendar/quickstart).

3. Build the project:

```bash
./cli-build.sh VERSION
```

Where `VERSION` is the version number you want to build. For example, `./cli-build.sh 1.0.0`.

4. Move the built binary to your PATH:

```bash
mv ./dist/TimeTrack /usr/local/bin
```

5. Run TimeTrack:

```bash
timetrack
```

That should be it! You should now be able to run TimeTrack from anywhere using the `timetrack` command.

### Update

To update TimeTrack to the latest version, simply re-run the installation command for your respective operating system.

## Collaborating

Contributions are what make the open source community such an amazing place to be, learn, inspire, and create. Any contributions you make are **greatly appreciated**. For major changes, please open an issue first to discuss what you would like to change.

To contribute to TimeTrack, please fork the repository, create a new branch, and submit a pull request.

## License

TimeTrack is released under the [MIT License](https://choosealicense.com/licenses/mit/).
