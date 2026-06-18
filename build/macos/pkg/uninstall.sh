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

# Remove LaunchAgent
PLIST_PATH="$HOME/Library/LaunchAgents/dev.aarju.speedtest-tray.plist"
if [ -f "$PLIST_PATH" ]; then
    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    rm "$PLIST_PATH"
    echo "Removed LaunchAgent."
fi

# Remove app
rm -rf "$APP"

read -r -p "Remove configuration, logs, and history? This includes config.json, app.log, and history.json. [y/N] " REMOVE_DATA
case "$REMOVE_DATA" in
    [yY][eE][sS]|[yY])
        rm -rf "$HOME/Library/Application Support/SpeedTest Tray"
        echo "Removed app data."
        ;;
esac

echo "SpeedTest Tray has been uninstalled."
