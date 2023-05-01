#!/bin/bash

SERVICE_NAME="spacemond"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# Stop the service
sudo systemctl stop ${SERVICE_NAME}

# Disable the service
sudo systemctl disable ${SERVICE_NAME}

# Remove the service file
sudo rm ${SERVICE_FILE}

# Reload the systemd daemon
sudo systemctl daemon-reload

# Remove the user created for the service
sudo userdel ${SERVICE_NAME}

# Print a message indicating that the uninstallation is complete
echo "The ${SERVICE_NAME} service has been uninstalled."