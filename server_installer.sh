#!/usr/bin/env bash
#
# This script downloads the chiSSL binary for your current platform from Github artifacts and installs it.
# It also creates the necessary systemd service scripts and initial configuration file.
# NOTE: This script generates a randomized admin password for chiSSL which can be found at /etc/chissl.json.
#
#
# Usage:
#   ./script_name.sh <domain_name> [port]
#
# Arguments:
#   domain_name: A fully qualified domain name (required)
#   port: An optional port number (default is 443 if not provided)
#
# Example:
#   ./script_name.sh subdomain.example.com
#  or
#   ./script_name.sh subdomain.example.com 8443
#
# To download and execute this script from a GitHub public repository in a single line:
#   bash <(curl -s https://raw.githubusercontent.com/NextChapterSoftware/chissl/v1.1/server_installer.sh) <domain_name> [port]

# Target chiSSL version
LATEST_VERSION="1.1"

# Function to display usage
usage() {
    echo "Usage: $0 <domain_name> [port]"
    exit 1
}

# Function to validate FQDN
is_valid_fqdn() {
    local fqdn=$1
    if [[ $fqdn =~ ^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$ && ${#fqdn} -le 253 ]]; then
        return 0
    else
        return 1
    fi
}

# Function to validate port number
is_valid_port() {
    local port=$1
    if [[ $port -ge 1 && $port -le 65535 ]]; then
        return 0
    else
        return 1
    fi
}

# Check if at least one argument is provided
if [ "$#" -lt 1 ]; then
    usage
fi

FQDN=$1
PORT=${2:-"443"}

# Validate the domain name
if ! is_valid_fqdn "$FQDN"; then
    echo "Error: $FQDN is not a valid fully qualified domain name (FQDN)."
    exit 1
fi

# Validate the port number
if ! is_valid_port "$PORT"; then
    echo "Error: $PORT is not a valid port number. It should be between 1 and 65535."
    exit 1
fi

# Define variables
REPO_OWNER="NextChapterSoftware"
REPO_NAME="chissl"
BASE_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/v$LATEST_VERSION"
INSTALL_PATH="/usr/local/bin/chissl"
SERVICE_NAME="chissl"
ADMIN_PASS=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13)

# Detect OS and architecture
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to supported binaries
case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64)
    ARCH="arm64"
    ;;
  armv5*)
    ARCH="armv5"
    ;;
  armv6*)
    ARCH="armv6"
    ;;
  armv7*)
    ARCH="armv7"
    ;;
  i386)
    ARCH="386"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Construct the download URL
ARCHIVE_NAME="${REPO_NAME}_${LATEST_VERSION}_${OS}_${ARCH}"
DOWNLOAD_URL="$BASE_URL/$ARCHIVE_NAME"

# Download the binary archive
echo "Downloading $DOWNLOAD_URL" $DOWNLOAD_URL
sudo curl -L -o $INSTALL_PATH $DOWNLOAD_URL
sudo chmod +x $INSTALL_PATH

# Create default configuration file with a randomized password for admin
echo "Creating /etc/chissl.json"
sudo tee /etc/chissl.json > /dev/null <<EOL
[
  {
    "username": "admin",
    "password": "$ADMIN_PASS",
    "addresses": [".*"],
    "is_admin": true
  }
]
EOL
if [ $? -ne 0 ]; then
    echo "Failed to create default /etc/chissl.json config file"
    exit 1
fi

# Create a systemd service file
echo "Creating systemd service"
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null <<EOL
[Unit]
Description=Chissl Service
After=network.target

[Service]
ExecStart=$INSTALL_PATH server --port $PORT --tls-domain $FQDN --authfile /etc/chissl.json
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOL

if [ $? -ne 0 ]; then
    echo "Failed to create /etc/systemd/system/$SERVICE_NAME.service script"
    exit 1
fi


# Reload systemd, enable and start the service
echo "Reloading systemd and starting $SERVICE_NAME"
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl stop $SERVICE_NAME || true
sudo systemctl start $SERVICE_NAME

if ! sudo systemctl status $SERVICE_NAME --no-pager ; then
     echo "$SERVICE_NAME startup failed. Check 'journalctl -xe' for more info "
     exit 1
fi

# Verify health endpoint - Checks health and SSL cert
echo
echo "Verifying service health endpoint at https://$FQDN:$PORT/health"
if ! curl https://$FQDN:$PORT/health -m 5 ; then
     echo "Failed to reach health endpoint via https://$FQDN:$PORT/health"
     exit 1
fi

echo
echo "============================================================================================="
echo "$SERVICE_NAME installation and setup complete."
echo "Admin user: admin"
echo "Admin password: $ADMIN_PASS"
echo "All user credentials stored at /etc/chissl.json"
echo "You can add additional users by editing /etc/chissl.json file or using the chiSSL rest API"
echo "chiSSL is watching /etc/chissl.json and will auto-reload upon any changes"
echo "============================================================================================="