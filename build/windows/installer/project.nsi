!define PRODUCT_NAME "SpeedTest Tray"
!define PRODUCT_VERSION "1.0.2"
!define PRODUCT_PUBLISHER "Aarju Pal"
!define INSTALL_DIR "$LOCALAPPDATA\Programs\SpeedTest Tray"
!define UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\SpeedTest Tray"
!define RUN_KEY "Software\Microsoft\Windows\CurrentVersion\Run"
!define EXE_NAME "speedtest-tray.exe"

!ifdef ARG_WAILS_ARM64_BINARY
    OutFile "..\..\bin\speedtest-tray-arm64-installer.exe"
!else
    OutFile "..\..\bin\speedtest-tray-amd64-installer.exe"
!endif


!include "winmessages.nsh"
!include "WordFunc.nsh"
!include "LogicLib.nsh"
!include "nsDialogs.nsh"
!insertmacro WordReplace

Var LaunchAtLoginCheckbox

; Use per-user installation — no UAC required
RequestExecutionLevel user
InstallDir "${INSTALL_DIR}"

; Pages
Page directory
Page custom LaunchOptionsPage LaunchOptionsLeave
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Function LaunchOptionsPage
    nsDialogs::Create 1018
    Pop $0
    ${If} $0 == error
        Abort
    ${EndIf}

    ${NSD_CreateCheckbox} 0 12u 100% 12u "Launch SpeedTest Tray at login"
    Pop $LaunchAtLoginCheckbox
    ${NSD_SetState} $LaunchAtLoginCheckbox ${BST_UNCHECKED}

    nsDialogs::Show
FunctionEnd

Function LaunchOptionsLeave
    ${NSD_GetState} $LaunchAtLoginCheckbox $LaunchAtLoginCheckbox
FunctionEnd

Section "Install"
    SetOutPath "$INSTDIR"
    File "..\..\bin\speedtest-tray.exe"
    ; Start Menu shortcut
    CreateDirectory "$SMPROGRAMS\SpeedTest Tray"
    CreateShortcut "$SMPROGRAMS\SpeedTest Tray\SpeedTest Tray.lnk" "$INSTDIR\${EXE_NAME}"

    ; Add install dir to user PATH
    ReadRegStr $0 HKCU "Environment" "Path"
    StrCpy $0 "$0;$INSTDIR"
    WriteRegExpandStr HKCU "Environment" "Path" "$0"
    ; Broadcast WM_SETTINGCHANGE so open terminals pick up new PATH
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
    ; Optional launch at login
    ${If} $LaunchAtLoginCheckbox == ${BST_CHECKED}
        WriteRegStr HKCU "${RUN_KEY}" "${PRODUCT_NAME}" "$INSTDIR\${EXE_NAME}"
    ${EndIf}
    ; Register uninstaller
    WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayName" "${PRODUCT_NAME}"
    WriteRegStr HKCU "${UNINSTALL_KEY}" "UninstallString" "$INSTDIR\uninstall.exe"
    WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
    WriteRegStr HKCU "${UNINSTALL_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
    WriteUninstaller "$INSTDIR\uninstall.exe"
    ; Launch app after install
    ExecShell "" "$INSTDIR\${EXE_NAME}"
SectionEnd

Section "Uninstall"
    ; Remove install dir from user PATH
    ReadRegStr $0 HKCU "Environment" "Path"
    ; Strip "$INSTDIR;" from PATH string
    ${WordReplace} $0 "$INSTDIR;" "" "+" $0
    ${WordReplace} $0 ";$INSTDIR" "" "+" $0
    WriteRegExpandStr HKCU "Environment" "Path" "$0"
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

    ; Remove autostart Run key if present
    DeleteRegValue HKCU "${RUN_KEY}" "${PRODUCT_NAME}"

    MessageBox MB_YESNO "Remove configuration, logs, and history? This includes config.json, app.log, and history.json." IDNO SkipDataRemoval
    RMDir /r "$APPDATA\SpeedTest Tray"
SkipDataRemoval:

    Delete "$INSTDIR\${EXE_NAME}"
    Delete "$INSTDIR\uninstall.exe"
    Delete "$SMPROGRAMS\SpeedTest Tray\SpeedTest Tray.lnk"
    RMDir "$SMPROGRAMS\SpeedTest Tray"
    RMDir "$INSTDIR"
    DeleteRegKey HKCU "${UNINSTALL_KEY}"
SectionEnd
