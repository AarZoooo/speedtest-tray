!define PRODUCT_NAME "SpeedTest Tray"
!define PRODUCT_VERSION "1.1.1"
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

; Request user or admin execution level (user level is fine for local appdata installs)
RequestExecutionLevel user
InstallDir "${INSTALL_DIR}"

;--------------------------------
; Include Modern UI

!include "MUI2.nsh"
!include "winmessages.nsh"
!include "WordFunc.nsh"
!include "LogicLib.nsh"
!include "nsDialogs.nsh"
!insertmacro WordReplace

;--------------------------------
; Interface Settings

!define MUI_ABORTWARNING
!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"

;--------------------------------
; Pages

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY

Var LaunchAtLoginCheckbox
Page custom LaunchOptionsPage LaunchOptionsLeave

!insertmacro MUI_PAGE_INSTFILES

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

;--------------------------------
; Languages

!insertmacro MUI_LANGUAGE "English"

;--------------------------------
; Custom Launch Options Page

Function LaunchOptionsPage
    nsDialogs::Create 1018
    Pop $0
    ${If} $0 == error
        Abort
    ${EndIf}

    ${NSD_CreateLabel} 0 0 100% 24u "Select additional options for SpeedTest Tray installation."
    Pop $0

    ${NSD_CreateCheckbox} 0 30u 100% 12u "Launch SpeedTest Tray at login"
    Pop $LaunchAtLoginCheckbox
    ${NSD_SetState} $LaunchAtLoginCheckbox ${BST_UNCHECKED}

    nsDialogs::Show
FunctionEnd

Function LaunchOptionsLeave
    ${NSD_GetState} $LaunchAtLoginCheckbox $LaunchAtLoginCheckbox
FunctionEnd

;--------------------------------
; Install Section

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

;--------------------------------
; Uninstall Section

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
