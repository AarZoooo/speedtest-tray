#!/bin/bash
set -e

APP="/Applications/SpeedTest Tray.app"

if [ ! -d "$APP" ]; then
    echo "SpeedTest Tray is not installed in /Applications."
    exit 1
fi

# Placeholder: symlink removal added by feature/installer/cli-alias
# Placeholder: LaunchAgent removal added by feature/installer/autostart

# Remove app
rm -rf "$APP"

# Placeholder: data cleanup dialog added by feature/installer/autostart

echo "SpeedTest Tray has been uninstalled."
