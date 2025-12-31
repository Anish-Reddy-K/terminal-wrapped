#!/bin/bash
#
# terminal-wrapped - instant developer stats
# run via: curl -fsSL arkr.ca/terminal-wrapped | bash
#
# this script downloads and runs the terminal-wrapped binary
# no installation required - just runs and displays your stats
#

set -e

# colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # no color

# configuration
REPO="Anish-Reddy-K/terminal-wrapped"
BINARY_NAME="terminal-wrapped"

# detect os
detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        *)        echo "unsupported" ;;
    esac
}

# detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64)   echo "amd64" ;;
        amd64)    echo "amd64" ;;
        arm64)    echo "arm64" ;;
        aarch64)  echo "arm64" ;;
        *)        echo "unsupported" ;;
    esac
}

# main execution
main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unsupported" ]; then
        echo -e "${RED}Error: Unsupported operating system. terminal-wrapped only supports macOS and Linux.${NC}"
        exit 1
    fi

    if [ "$ARCH" = "unsupported" ]; then
        echo -e "${RED}Error: Unsupported architecture. terminal-wrapped supports x86_64 and arm64.${NC}"
        exit 1
    fi

    # get latest release URL from GitHub
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}-${OS}-${ARCH}"
    
    # create temp directory
    TMP_DIR=$(mktemp -d)
    TMP_BINARY="${TMP_DIR}/${BINARY_NAME}"

    # cleanup on exit
    trap "rm -rf ${TMP_DIR}" EXIT

    # download binary
    echo -e "${YELLOW} downloading terminal-wrapped for ${OS}/${ARCH}...${NC}"
    
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

    # make executable
    chmod +x "$TMP_BINARY"

    # run it!
    echo -e "${GREEN}✨ Running terminal-wrapped...${NC}\n"
    "$TMP_BINARY" "$@"
}

# run main function with all script arguments
main "$@"

