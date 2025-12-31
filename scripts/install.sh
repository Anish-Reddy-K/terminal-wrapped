#!/bin/bash
#
# Terminal Wrapped - Instant Developer Stats
# Run via: curl -fsSL arkr.ca/terminal-wrapped | bash
#
# This script downloads and runs the terminal-wrapped binary
# No installation required - just runs and displays your stats
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="Anish-Reddy-K/terminal-wrapped"
BINARY_NAME="terminal-wrapped"

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        *)        echo "unsupported" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64)   echo "amd64" ;;
        amd64)    echo "amd64" ;;
        arm64)    echo "arm64" ;;
        aarch64)  echo "arm64" ;;
        *)        echo "unsupported" ;;
    esac
}

# Main execution
main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unsupported" ]; then
        echo -e "${RED}Error: Unsupported operating system. Terminal Wrapped only supports macOS and Linux.${NC}"
        exit 1
    fi

    if [ "$ARCH" = "unsupported" ]; then
        echo -e "${RED}Error: Unsupported architecture. Terminal Wrapped supports x86_64 and arm64.${NC}"
        exit 1
    fi

    # Get latest release URL from GitHub
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_BINARY="${TMP_DIR}/${BINARY_NAME}"

    # Cleanup on exit
    trap "rm -rf ${TMP_DIR}" EXIT

    # Download binary
    echo -e "${YELLOW}⬇️  Downloading Terminal Wrapped for ${OS}/${ARCH}...${NC}"
    
    if command -v curl &> /dev/null; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_BINARY" 2>/dev/null || {
            echo -e "${RED}Error: Failed to download binary from ${DOWNLOAD_URL}${NC}"
            echo -e "${YELLOW}You may need to build from source: https://github.com/${REPO}${NC}"
            exit 1
        }
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TMP_BINARY" 2>/dev/null || {
            echo -e "${RED}Error: Failed to download binary from ${DOWNLOAD_URL}${NC}"
            echo -e "${YELLOW}You may need to build from source: https://github.com/${REPO}${NC}"
            exit 1
        }
    else
        echo -e "${RED}Error: Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi

    # Make executable
    chmod +x "$TMP_BINARY"

    # Run it!
    echo -e "${GREEN}✨ Running Terminal Wrapped...${NC}\n"
    "$TMP_BINARY" "$@"
}

# Run main function with all script arguments
main "$@"

