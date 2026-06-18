!define PRODUCT_NAME "SpeedTest Tray"
!define PRODUCT_VERSION "1.0.2"
!define PRODUCT_PUBLISHER "Aarju Pal"
!define INSTALL_DIR "$LOCALAPPDATA\Programs\SpeedTest Tray"
!define UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\SpeedTest Tray"
!define EXE_NAME "speedtest-tray.exe"

!include "winmessages.nsh"
!include "WordFunc.nsh"
!insertmacro StrRep

; Use per-user installation — no UAC required
RequestExecutionLevel user
InstallDir "${INSTALL_DIR}"

; Pages
Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install"
    SetOutPath "$INSTDIR"
    File "${EXE_NAME}"
    ; Start Menu shortcut
    CreateDirectory "$SMPROGRAMS\SpeedTest Tray"
    CreateShortcut "$SMPROGRAMS\SpeedTest Tray\SpeedTest Tray.lnk" "$INSTDIR\${EXE_NAME}"

    ; Add install dir to user PATH
    ReadRegStr $0 HKCU "Environment" "Path"
    StrCpy $0 "$0;$INSTDIR"
    WriteRegExpandStr HKCU "Environment" "Path" "$0"
    ; Broadcast WM_SETTINGCHANGE so open terminals pick up new PATH
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
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
    ${StrRep} $0 $0 "$INSTDIR;" ""
    ${StrRep} $0 $0 ";$INSTDIR" ""
    WriteRegExpandStr HKCU "Environment" "Path" "$0"
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

    ; TODO: Run key removal added by autostart branch
    ; TODO: Data cleanup dialog added by autostart branch
    Delete "$INSTDIR\${EXE_NAME}"
    Delete "$INSTDIR\uninstall.exe"
    Delete "$SMPROGRAMS\SpeedTest Tray\SpeedTest Tray.lnk"
    RMDir "$SMPROGRAMS\SpeedTest Tray"
    RMDir "$INSTDIR"
    DeleteRegKey HKCU "${UNINSTALL_KEY}"
SectionEnd
