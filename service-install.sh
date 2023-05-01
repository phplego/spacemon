#!/bin/bash

SERVICE_NAME="spacemond"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
EXECUTABLE_PATH="$(pwd)/${SERVICE_NAME}"

# Create the systemd service file
sudo bash -c "cat > ${SERVICE_FILE}" << EOL
[Unit]
Description=${SERVICE_NAME} service
After=network.target

[Service]
User=$USER
Group=${SERVICE_NAME}
WorkingDirectory=$(pwd)
ExecStart=${EXECUTABLE_PATH}
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOL


# Reload the systemd daemon
sudo systemctl daemon-reload

# Enable the service
sudo systemctl enable ${SERVICE_NAME}

# Start the service
sudo systemctl start ${SERVICE_NAME}

# Print a message indicating that the installation is complete
echo "The ${SERVICE_NAME} service has been installed and started."