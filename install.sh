#!/bin/bash

set -e

# TimeTrack CLI installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/linusromland/TimeTrack/master/install.sh | bash

REPO="linusromland/TimeTrack"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="timetrack"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for system install
check_sudo() {
    if [[ $EUID -eq 0 ]]; then
        SUDO=""
    else
        if command -v sudo >/dev/null 2>&1; then
            SUDO="sudo"
            print_info "Will use sudo for system installation"
        else
            print_error "This script requires sudo for system installation"
            print_info "Please install to user directory with: INSTALL_DIR=\"\$HOME/.local/bin\" $0"
            exit 1
        fi
    fi
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $OS in
        linux*)
            PLATFORM="linux"
            ;;
        darwin*)
            PLATFORM="darwin"
            ;;
        *)
            print_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    print_info "Detected platform: $PLATFORM-$ARCH"
}

# Get latest release version
get_latest_version() {
    print_info "Fetching latest release..."
    LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi
    
    print_info "Latest version: $LATEST_VERSION"
}

# Download and install binary
install_binary() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/timetrack-cli-$PLATFORM-$ARCH.tar.gz"
    TEMP_DIR=$(mktemp -d)
    
    print_info "Downloading from: $DOWNLOAD_URL"
    
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/timetrack.tar.gz"; then
        print_error "Failed to download TimeTrack CLI"
        exit 1
    fi
    
    cd "$TEMP_DIR"
    tar -xzf timetrack.tar.gz
    
    # Make sure install directory exists
    $SUDO mkdir -p "$INSTALL_DIR"
    
    # Remove existing binary if it exists
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        print_info "Removing existing TimeTrack binary..."
        $SUDO rm -f "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Install binary
    $SUDO cp "timetrack-cli-$PLATFORM-$ARCH" "$INSTALL_DIR/$BINARY_NAME"
    $SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # Cleanup
    rm -rf "$TEMP_DIR"
    
    print_info "TimeTrack CLI installed to $INSTALL_DIR/$BINARY_NAME"
}

# Setup bash completion
setup_bash_completion() {
    COMPLETION_DIR=""
    COMPLETION_SUDO=""
    
    # Remove existing completions first
    print_info "Removing any existing bash completions..."
    
    # Check common completion directories and remove existing files
    for dir in "/etc/bash_completion.d" "/usr/local/etc/bash_completion.d" "$HOME/.local/share/bash-completion/completions"; do
        if [ -f "$dir/timetrack" ]; then
            if [[ "$dir" == "$HOME"* ]]; then
                rm -f "$dir/timetrack"
            else
                $SUDO rm -f "$dir/timetrack"
            fi
            print_info "Removed existing completion from $dir"
        fi
    done
    
    # Try different completion directories
    if [ -d "/etc/bash_completion.d" ]; then
        COMPLETION_DIR="/etc/bash_completion.d"
        COMPLETION_SUDO="$SUDO"
    elif [ -d "/usr/local/etc/bash_completion.d" ]; then
        COMPLETION_DIR="/usr/local/etc/bash_completion.d"
        COMPLETION_SUDO="$SUDO"
    elif [ -d "$HOME/.local/share/bash-completion/completions" ]; then
        COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
        COMPLETION_SUDO="" # Don't use sudo for user directory
    fi
    
    if [ -n "$COMPLETION_DIR" ]; then
        print_info "Setting up bash completion..."
        
        # Generate completion script
        COMPLETION_SCRIPT="# TimeTrack bash completion
_timetrack_completion() {
    local cur prev words cword
    _init_completion || return

    # Define available commands
    local commands=\"add list login register settings dashboard\"
    
    # If we're completing the first argument (command)
    if [[ \$cword -eq 1 ]]; then
        COMPREPLY=(\$(compgen -W \"\$commands\" -- \"\$cur\"))
        return
    fi
    
    # For subcommands, we can add more specific completion later
    # For now, just don't complete anything after the command
    COMPREPLY=()
}

complete -F _timetrack_completion timetrack"
        
        $COMPLETION_SUDO mkdir -p "$COMPLETION_DIR"
        echo "$COMPLETION_SCRIPT" | $COMPLETION_SUDO tee "$COMPLETION_DIR/timetrack" >/dev/null
        
        print_info "Bash completion installed to $COMPLETION_DIR/timetrack"
        print_warn "Restart your shell or run 'source $COMPLETION_DIR/timetrack' to enable completion"
    else
        print_warn "Could not find bash completion directory. Completion not installed."
    fi
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        VERSION_OUTPUT=$("$BINARY_NAME" --version)
        print_info "Installation successful!"
        print_info "Installed version: $VERSION_OUTPUT"
        print_info "Run '$BINARY_NAME --help' to get started"
    else
        print_error "Installation failed - binary not found in PATH"
        print_info "Make sure $INSTALL_DIR is in your PATH"
        exit 1
    fi
}

# Main installation flow
main() {
    print_info "Installing TimeTrack CLI..."
    
    # Allow user to override install directory
    if [ -n "$TIMETRACK_INSTALL_DIR" ]; then
        INSTALL_DIR="$TIMETRACK_INSTALL_DIR"
        print_info "Using custom install directory: $INSTALL_DIR"
    fi
    
    detect_platform
    
    # Check sudo requirements for system directories
    if [[ "$INSTALL_DIR" == "/usr"* ]] || [[ "$INSTALL_DIR" == "/opt"* ]]; then
        check_sudo
    else
        SUDO=""
    fi
    
    get_latest_version
    install_binary
    setup_bash_completion
    verify_installation
    
    print_info "TimeTrack CLI installation complete!"
}

# Run main function
main "$@"