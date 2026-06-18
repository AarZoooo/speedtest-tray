#!/bin/bash
set -e

APP="/Applications/SpeedTest Tray.app"

if [ ! -d "$APP" ]; then
    echo "SpeedTest Tray is not installed in /Applications."
    exit 1
fi

# Remove CLI symlink
CLI_LINK="/usr/local/bin/speedtest-tray"
if [ -L "$CLI_LINK" ]; then
    rm "$CLI_LINK"
    echo "Removed CLI symlink."
fi

# Placeholder: LaunchAgent removal added by feature/installer/autostart

# Remove app
rm -rf "$APP"

# Placeholder: data cleanup dialog added by feature/installer/autostart

echo "SpeedTest Tray has been uninstalled."
