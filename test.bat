@echo off
echo Running Sleep No More Tests...
echo.

echo [1/3] Running unit tests...
echo ----------------------------------------
go test -v -short ./...
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Unit tests failed!
    echo.
    exit /b 1
)
echo Unit tests passed!

echo.
echo [2/3] Running integration tests...
echo ----------------------------------------
go test -v -tags=integration ./...
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Integration tests failed!
    echo.
    exit /b 1
)
echo Integration tests passed!

echo.
echo [3/3] Running benchmarks...
echo ----------------------------------------
go test -bench=. -benchmem ./...
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Benchmarks failed!
    echo.
    exit /b 1
)
echo Benchmarks completed successfully!

echo.
echo ========================================
echo Summary:
echo - Unit tests: PASS
echo - Integration tests: PASS  
echo - Benchmarks: PASS
echo ========================================

