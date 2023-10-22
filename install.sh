#!/bin/bash

# Fetch the latest version tag from GitHub API
VERSION=$(curl -s https://api.github.com/repos/linusromland/TimeTrack/releases/latest | grep "tag_name" | cut -d '"' -f 4)

# Construct the download URL based on the OS and architecture
DOWNLOAD_URL="https://github.com/linusromland/TimeTrack/releases/download/$VERSION/TimeTrack_$(uname)_$(uname -m | sed 's/x86_64/x86_64/;s/i686/i386/')"

# Check for ARM architecture and adjust the download URL if needed
if [ "$(uname -m)" = "armv7l" ]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}v7"
fi
DOWNLOAD_URL="${DOWNLOAD_URL}.tar.gz"

# Download the archive
echo "Downloading TimeTrack from $DOWNLOAD_URL..."
curl -L -o TimeTrack.tar.gz $DOWNLOAD_URL

# Extract the archive
tar xzf TimeTrack.tar.gz

# Move the binary to the appropriate directory
sudo mv TimeTrack /usr/local/bin/

# Clean up
rm TimeTrack.tar.gz

echo "TimeTrack has been installed!"
