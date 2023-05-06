#!/bin/bash

# Import service name and paths
source ./service-config.sh

# Show logs
sudo journalctl -u $SERVICE_NAME -f -n 1000
