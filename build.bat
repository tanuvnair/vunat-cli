@echo off
setlocal

set BIN_DIR=.\bin
set CMD_DIR=.\cmd\vunat

if not exist "%BIN_DIR%" (
    mkdir "%BIN_DIR%"
)

echo Building vunat-cli...

go build -o "%BIN_DIR%\vunat.exe" "%CMD_DIR%"

if errorlevel 1 (
    echo Build failed.
    exit /b 1
)

echo Build successful: %BIN_DIR%\vunat.exe
