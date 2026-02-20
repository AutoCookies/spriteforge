param([string]$Version="1.0.0",[string]$Commit="dev",[string]$BuildDate="1970-01-01T00:00:00Z")
$ErrorActionPreference='Stop'
New-Item -ItemType Directory -Force -Path dist | Out-Null
$ld="-X pixelc/internal/version.Version=$Version -X pixelc/internal/version.Commit=$Commit -X pixelc/internal/version.BuildDate=$BuildDate"
go build -ldflags $ld -o dist/pixelc.exe ./cmd/pixelc
Compress-Archive -Path dist/pixelc.exe -DestinationPath "dist/pixelc-$Version-windows-x64.zip" -Force
$makensis=Get-Command makensis -ErrorAction SilentlyContinue
if ($makensis) {
  $env:VERSION=$Version; $env:OUT_EXE="dist/pixelc-setup-$Version-windows-x64.exe"; $env:BIN_PATH=(Resolve-Path dist/pixelc.exe)
  makensis /DVERSION=$Version /DOUT_EXE=$env:OUT_EXE /DBIN_PATH=$env:BIN_PATH scripts/package/windows_installer.nsi
}
Write-Host "Artifacts:"; Get-ChildItem dist | ForEach-Object { Write-Host $_.FullName }
