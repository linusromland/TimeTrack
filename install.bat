@echo off
setlocal

:: Fetch the latest version tag from GitHub API
FOR /F "tokens=*" %%i IN ('curl -s https://api.github.com/repos/linusromland/TimeTrack/releases/latest ^| find "tag_name"') DO SET "VERSION=%%i"
SET "VERSION=%VERSION:~10,-1%"

:: Construct the download URL
SET "DOWNLOAD_URL=https://github.com/linusromland/TimeTrack/releases/download/%VERSION%/TimeTrack_Windows_x86_64.zip"

:: Download the archive
echo Downloading TimeTrack from %DOWNLOAD_URL%...
curl -L -o TimeTrack.zip %DOWNLOAD_URL%

:: Extract the archive
tar -xf TimeTrack.zip

:: Move the binary to the appropriate directory
move TimeTrack.exe %USERPROFILE%\AppData\Local\Microsoft\WindowsApps\

:: Clean up
del TimeTrack.zip

echo TimeTrack has been installed!
endlocal
