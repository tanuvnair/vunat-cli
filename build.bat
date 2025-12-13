@echo off
setlocal

set BIN_DIR=.\bin
set CMD_DIR=.\cmd\vunat

if not exist "%BIN_DIR%" (
    mkdir "%BIN_DIR%"
)

echo Building vunat-cli for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o "%BIN_DIR%\vunat-windows-amd64.exe" "%CMD_DIR%"

if errorlevel 1 (
    echo Windows build failed.
    exit /b 1
)

echo Building vunat-cli for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o "%BIN_DIR%\vunat-linux-amd64" "%CMD_DIR%"

if errorlevel 1 (
    echo Linux build failed.
    exit /b 1
)

echo Build successful:
echo   - Windows: %BIN_DIR%\vunat-windows-amd64.exe
echo   - Linux: %BIN_DIR%\vunat-linux-amd64

