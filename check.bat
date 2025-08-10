@echo off
setlocal enabledelayedexpansion

set "FAILED=0"

echo === 1. Code Formatting ===
echo Checking gofmt...
for /f %%i in ('gofmt -l .') do (
    set "FAILED=1"
    echo FAIL: File needs formatting: %%i
)
if !FAILED! == 0 (
    echo PASS: All files are properly formatted
) else (
    echo To fix: gofmt -w .
)
echo.

echo Checking goimports...
go install golang.org/x/tools/cmd/goimports@latest >nul 2>&1
for /f %%i in ('goimports -l .') do (
    set "FAILED=1"
    echo FAIL: File has import issues: %%i
)
if !FAILED! == 0 (
    echo PASS: All imports are properly formatted
) else (
    echo To fix: goimports -w .
)
echo.

echo === 2. Go Vet ===
go vet ./...
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: go vet found issues
) else (
    echo PASS: go vet completed successfully
)
echo.

echo === 3. Static Analysis ===
echo Installing/updating analysis tools...
go install honnef.co/go/tools/cmd/staticcheck@latest >nul 2>&1
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest >nul 2>&1
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest >nul 2>&1
go install github.com/uudashr/gocognit/cmd/gocognit@latest >nul 2>&1
go install github.com/gordonklaus/ineffassign@latest >nul 2>&1
go install github.com/client9/misspell/cmd/misspell@latest >nul 2>&1
go install github.com/securego/gosec/v2/cmd/gosec@latest >nul 2>&1

echo Running staticcheck...
staticcheck ./...
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: staticcheck found issues
) else (
    echo PASS: staticcheck completed successfully
)
echo.

echo === 4. Linting (golangci-lint) ===
golangci-lint run --timeout=5m --config=.golangci.yml
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: golangci-lint found issues
) else (
    echo PASS: golangci-lint completed successfully
)
echo.

echo === 5. Security Scan ===
gosec -fmt text ./...
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: gosec found security issues
) else (
    echo PASS: gosec found no security issues
)
echo.

echo === 6. Code Complexity ===
echo Checking cyclomatic complexity...
for /f "delims=" %%i in ('gocyclo -over 15 . 2^>nul') do (
    set "FAILED=1"
    echo WARNING: High cyclomatic complexity: %%i
)
if !FAILED! == 0 (
    echo PASS: All functions have acceptable cyclomatic complexity (≤15)
)

echo Checking cognitive complexity...
for /f "delims=" %%i in ('gocognit -over 20 . 2^>nul') do (
    set "FAILED=1"
    echo WARNING: High cognitive complexity: %%i
)
if !FAILED! == 0 (
    echo PASS: All functions have acceptable cognitive complexity (≤20)
)
echo.

echo Checking ineffective assignments...
for /f "delims=" %%i in ('ineffassign ./... 2^>nul') do (
    set "FAILED=1"
    echo FAIL: Ineffective assignment: %%i
)
if !FAILED! == 0 (
    echo PASS: No ineffective assignments found
)
echo.

echo Checking spelling...
for /f "delims=" %%i in ('misspell . 2^>nul') do (
    set "FAILED=1"
    echo WARNING: Spelling error: %%i
)
if !FAILED! == 0 (
    echo PASS: No spelling errors found
)
echo.

echo === 7. Tests ===
echo Running unit tests...
go test -v -short -coverprofile="coverage.txt" -covermode=atomic ./...
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: Tests failed
) else (
    echo PASS: All tests passed
)
echo.

echo Running benchmarks...
go test -bench=. -benchmem ./...
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: Benchmarks failed
) else (
    echo PASS: Benchmarks completed successfully
)
echo.

echo === 8. Test Coverage ===
if exist coverage.txt (
    for /f "tokens=3" %%i in ('go tool cover -func^=coverage.txt ^| findstr "total:"') do (
        set "COVERAGE=%%i"
        echo Total coverage: !COVERAGE!
        for /f "tokens=1 delims=%%" %%j in ("!COVERAGE!") do (
            if %%j geq 80 (
                echo EXCELLENT: Coverage is %%j%% ^(≥80%%^)
            ) else if %%j geq 60 (
                echo GOOD: Coverage is %%j%% ^(≥60%%, but could be improved^)
            ) else (
                echo POOR: Coverage is %%j%% ^(needs improvement^)
                set "FAILED=1"
            )
        )
    )
) else (
    echo ERROR: Coverage file not found
    set "FAILED=1"
)
echo.

echo === 9. Build Test ===
echo Building for Windows AMD64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-H windowsgui -s -w" -o sleepnomore_amd64.exe main.go
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: AMD64 build failed
) else (
    echo PASS: AMD64 build successful
    del sleepnomore_amd64.exe >nul 2>&1
)

echo Building for Windows ARM64...
set GOARCH=arm64
go build -ldflags="-H windowsgui -s -w" -o sleepnomore_arm64.exe main.go
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: ARM64 build failed
) else (
    echo PASS: ARM64 build successful
    del sleepnomore_arm64.exe >nul 2>&1
)
echo.

echo === 10. Dependency Check ===
go mod verify
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: Module verification failed
) else (
    echo PASS: All dependencies verified
)

go mod tidy
if !errorlevel! neq 0 (
    set "FAILED=1"
    echo FAIL: go mod tidy found issues
) else (
    echo PASS: Dependencies are clean
)
echo.

echo ===============================================
echo FINAL RESULT
echo ===============================================
if !FAILED! == 0 (
    echo ALL CHECKS PASSED!
    echo Your code is ready for GitHub Actions
    exit /b 0
) else (
    echo SOME CHECKS FAILED
    echo Please fix the issues above before pushing
    exit /b 1
)