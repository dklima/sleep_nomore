@echo off
setlocal enabledelayedexpansion

echo ===============================================
echo Sleep No More - Quick Check
echo ===============================================
echo.

set "FAILED=0"
set "GREEN=[92m"
set "RED=[91m"
set "YELLOW=[93m"
set "NC=[0m"

echo Formatting...
for /f %%i in ('gofmt -l .') do (
    set "FAILED=1"
    echo %RED%Format issue:%NC% %%i
)

echo Imports...
go install golang.org/x/tools/cmd/goimports@latest >nul 2>&1
for /f %%i in ('goimports -l .') do (
    set "FAILED=1"
    echo %RED%Import issue:%NC% %%i
)

echo Vet...
go vet ./... >nul 2>&1
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo %RED%go vet found issues%NC%
)

echo Tests...
go test -short ./... >nul 2>&1
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo %RED%Tests failed%NC%
)

echo Build...
go build -o temp.exe main.go >nul 2>&1
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo %RED%Build failed%NC%
) else (
    del temp.exe >nul 2>&1
)

echo.
if !FAILED! == 0 (
    echo %GREEN%✓ Quick check passed!%NC%
    exit /b 0
) else (
    echo %RED%✗ Issues found. Run 'check.bat' for details%NC%
    exit /b 1
)