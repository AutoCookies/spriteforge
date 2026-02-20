!define APPNAME "pixelc"
!define VERSION "${VERSION}"
OutFile "${OUT_EXE}"
InstallDir "$PROGRAMFILES64\pixelc"
Page directory
Page instfiles
Section "Install"
  SetOutPath "$INSTDIR"
  File "${BIN_PATH}"
  WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd
Section "Uninstall"
  Delete "$INSTDIR\pixelc.exe"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"
SectionEnd
