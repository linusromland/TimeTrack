@echo off
setlocal

:: Fetch the latest version tag from GitHub API
FOR /F "tokens=2 delims=:" %%i IN ('curl -s https://api.github.com/repos/linusromland/TimeTrack/releases/latest ^| findstr "tag_name"') DO SET "RAW_VERSION=%%i"

:: Trim and clean the version string
FOR /F "tokens=1 delims=," %%j IN ("%RAW_VERSION%") DO SET "VERSION=%%j"
SET "VERSION=%VERSION:~2,-1%"

:: Construct the download URL
SET "DOWNLOAD_URL=https://github.com/linusromland/TimeTrack/releases/download/%VERSION%/TimeTrack_Windows_x86_64.zip"

:: Download the archive
echo Downloading TimeTrack from %DOWNLOAD_URL%...
curl -L -o TimeTrack.zip %DOWNLOAD_URL%

:: Extract the archive using PowerShell
powershell -command "Expand-Archive -Path .\TimeTrack.zip -DestinationPath ./timetrack_tmp"

:: Move the binary to the appropriate directory
move timetrack_tmp\TimeTrack.exe %USERPROFILE%\AppData\Local\Microsoft\WindowsApps\

:: Clean up
del TimeTrack.zip
rmdir /s /q timetrack_tmp

echo TimeTrack has been installed!
endlocal
