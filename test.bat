@echo off
echo Running Sleep No More Tests...
echo.

echo [1/3] Running unit tests...
go test -v -short ./...
if %ERRORLEVEL% NEQ 0 (
    echo Unit tests failed!
    exit /b 1
)

echo.
echo [2/3] Running integration tests...
go test -v -tags=integration ./...
if %ERRORLEVEL% NEQ 0 (
    echo Integration tests failed!
    exit /b 1
)

echo.
echo [3/3] Running benchmarks...
go test -bench=. -benchmem ./...
if %ERRORLEVEL% NEQ 0 (
    echo Benchmarks failed!
    exit /b 1
)

echo.
