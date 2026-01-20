@echo off
setlocal

echo ==========================================
echo Start compiling Authos for Linux (AMD64)...
echo ==========================================

REM Set Go environment variables
REM CGO_ENABLED=0: Disable CGO for pure Go build
REM GOOS=linux: Target OS
REM GOARCH=amd64: Target Architecture
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

REM Print environment info for verification
echo [Debug] Current Build Environment:
go env GOOS GOARCH CGO_ENABLED

REM Execute build
echo [Info] Building...
go build -ldflags="-s -w" -trimpath -o authos-linux-amd64 .

if %errorlevel% neq 0 (
    echo [ERROR] Compilation failed!
    exit /b %errorlevel%
)

echo ==========================================
echo Compilation successful!
echo Output file: authos-linux-amd64
echo ==========================================
echo Tip: verify file type on Linux using 'file authos-linux-amd64'
echo Tip: verify architecture on Linux using 'uname -m'
pause
