#!/usr/bin/env pwsh
param(
    [switch]$Quick,
    [switch]$NoColor,
    [switch]$Verbose
)

# Color functions
function Write-Success { param([string]$Message) if ($NoColor) { Write-Host "✓ $Message" } else { Write-Host "✓ $Message" -ForegroundColor Green } }
function Write-Error { param([string]$Message) if ($NoColor) { Write-Host "✗ $Message" } else { Write-Host "✗ $Message" -ForegroundColor Red } }
function Write-Warning { param([string]$Message) if ($NoColor) { Write-Host "⚠ $Message" } else { Write-Host "⚠ $Message" -ForegroundColor Yellow } }
function Write-Info { param([string]$Message) if ($NoColor) { Write-Host "ℹ $Message" } else { Write-Host "ℹ $Message" -ForegroundColor Cyan } }
function Write-Section { param([string]$Message) if ($NoColor) { Write-Host "`n=== $Message ===" } else { Write-Host "`n=== $Message ===" -ForegroundColor Blue } }

$script:FailureCount = 0
$script:WarningCount = 0

function Test-Command {
    param([string]$Command, [string]$Description)
    if ($Verbose) { Write-Info "Running: $Command" }
    
    try {
        $result = Invoke-Expression $Command 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success $Description
            return $true
        } else {
            Write-Error "$Description (Exit code: $LASTEXITCODE)"
            if ($Verbose -and $result) { Write-Host $result }
            $script:FailureCount++
            return $false
        }
    } catch {
        Write-Error "$Description (Exception: $($_.Exception.Message))"
        $script:FailureCount++
        return $false
    }
}

function Install-Tools {
    Write-Section "Installing/Updating Tools"
    
    $tools = @(
        @{Name = "goimports"; Package = "golang.org/x/tools/cmd/goimports@latest"}
        @{Name = "staticcheck"; Package = "honnef.co/go/tools/cmd/staticcheck@latest"}
        @{Name = "golangci-lint"; Package = "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}
        @{Name = "gocyclo"; Package = "github.com/fzipp/gocyclo/cmd/gocyclo@latest"}
        @{Name = "gocognit"; Package = "github.com/uudashr/gocognit/cmd/gocognit@latest"}
        @{Name = "ineffassign"; Package = "github.com/gordonklaus/ineffassign@latest"}
        @{Name = "misspell"; Package = "github.com/client9/misspell/cmd/misspell@latest"}
        @{Name = "gosec"; Package = "github.com/securego/gosec/v2/cmd/gosec@latest"}
    )
    
    foreach ($tool in $tools) {
        if ($Verbose) { Write-Info "Installing $($tool.Name)..." }
        go install $tool.Package *>$null
    }
    Write-Success "All tools installed/updated"
}

Write-Host "Sleep No More - Local Quality Checks" -ForegroundColor Magenta
Write-Host "====================================" -ForegroundColor Magenta

if (-not $Quick) { Install-Tools }

# 1. Formatting
Write-Section "Code Formatting"

$unformatted = gofmt -l .
if ($unformatted) {
    Write-Error "Files need formatting:"
    $unformatted | ForEach-Object { Write-Host "  $_" }
    Write-Info "Fix with: gofmt -w ."
    $script:FailureCount++
} else {
    Write-Success "All files properly formatted"
}

if (-not $Quick) {
    go install golang.org/x/tools/cmd/goimports@latest *>$null
    $importIssues = goimports -l .
    if ($importIssues) {
        Write-Error "Files have import issues:"
        $importIssues | ForEach-Object { Write-Host "  $_" }
        Write-Info "Fix with: goimports -w ."
        $script:FailureCount++
    } else {
        Write-Success "All imports properly formatted"
    }
}

# 2. Go Vet
Write-Section "Go Vet"
Test-Command "go vet ./..." "go vet analysis" | Out-Null

# 3. Tests
Write-Section "Tests"
if ($Quick) {
    $args = @('test', '-short', './...')
    $result = & go $args 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Quick tests"
    } else {
        Write-Error "Quick tests (Exit code: $LASTEXITCODE)"
        $script:FailureCount++
    }
} else {
    $args = @('test', '-v', '-short', '-coverprofile=coverage.txt', '-covermode=atomic', './...')
    $result = & go $args 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Full test suite with coverage"
    } else {
        Write-Error "Full test suite with coverage (Exit code: $LASTEXITCODE)"
        $script:FailureCount++
    }
    
    if (Test-Path "coverage.txt") {
        Write-Section "Test Coverage"
        $coverage = & go tool cover -func coverage.txt | Select-String "total:"
        if ($coverage) {
            $coveragePercent = [regex]::Match($coverage, '(\d+\.\d+)%').Groups[1].Value
            $coverageNum = [double]$coveragePercent
            
            if ($coverageNum -ge 80) {
                Write-Success "Excellent coverage: $coveragePercent% (≥80%)"
            } elseif ($coverageNum -ge 60) {
                Write-Warning "Good coverage: $coveragePercent% (≥60%, but could be improved)"
                $script:WarningCount++
            } elseif ($coverageNum -ge 30) {
                Write-Warning "Acceptable coverage: $coveragePercent% (≥30%, consider improving)"
                $script:WarningCount++
            } else {
                Write-Error "Poor coverage: $coveragePercent% (needs improvement)"
                $script:FailureCount++
            }
        }
    }
}

# 4. Build Test
Write-Section "Build Test"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
Test-Command "go build -ldflags='-H windowsgui -s -w' -o sleepnomore_amd64.exe main.go" "Windows AMD64 build" | Out-Null
if (Test-Path "sleepnomore_amd64.exe") { Remove-Item "sleepnomore_amd64.exe" -Force }

if (-not $Quick) {
    $env:GOARCH = "arm64"
    Test-Command "go build -ldflags='-H windowsgui -s -w' -o sleepnomore_arm64.exe main.go" "Windows ARM64 build" | Out-Null
    if (Test-Path "sleepnomore_arm64.exe") { Remove-Item "sleepnomore_arm64.exe" -Force }
}

if (-not $Quick) {
    # 5. Static Analysis
    Write-Section "Static Analysis"
    Test-Command "staticcheck ./..." "staticcheck analysis" | Out-Null
    
    # 6. Linting
    Write-Section "Linting"
    Test-Command "golangci-lint run --timeout=5m --config=.golangci.yml" "golangci-lint analysis" | Out-Null
    
    # 7. Security
    Write-Section "Security Analysis"
    Test-Command "gosec -fmt text ./..." "gosec security scan" | Out-Null
    
    # 8. Code Complexity
    Write-Section "Code Complexity"
    
    $complexFunctions = gocyclo -over 15 . 2>$null
    if ($complexFunctions) {
        Write-Warning "Functions with high cyclomatic complexity (>15):"
        $complexFunctions | ForEach-Object { Write-Host "  $_" }
        $script:WarningCount++
    } else {
        Write-Success "All functions have acceptable cyclomatic complexity (≤15)"
    }
    
    $cognitiveFunctions = gocognit -over 20 . 2>$null
    if ($cognitiveFunctions) {
        Write-Warning "Functions with high cognitive complexity (>20):"
        $cognitiveFunctions | ForEach-Object { Write-Host "  $_" }
        $script:WarningCount++
    } else {
        Write-Success "All functions have acceptable cognitive complexity (≤20)"
    }
    
    $ineffective = ineffassign ./... 2>$null
    if ($ineffective) {
        Write-Error "Ineffective assignments found:"
        $ineffective | ForEach-Object { Write-Host "  $_" }
        $script:FailureCount++
    } else {
        Write-Success "No ineffective assignments found"
    }
    
    $misspelled = misspell . 2>$null
    if ($misspelled) {
        Write-Warning "Spelling errors found:"
        $misspelled | ForEach-Object { Write-Host "  $_" }
        $script:WarningCount++
    } else {
        Write-Success "No spelling errors found"
    }
    
    # 9. Dependencies
    Write-Section "Dependencies"
    Test-Command "go mod verify" "module verification" | Out-Null
    Test-Command "go mod tidy" "dependency cleanup" | Out-Null
}

# Final Results
Write-Host "`n====================================" -ForegroundColor Magenta
Write-Section "Final Results"

if ($script:FailureCount -eq 0) {
    Write-Success "ALL CHECKS PASSED!"
    if ($script:WarningCount -gt 0) {
        Write-Warning "$script:WarningCount warning(s) found - consider addressing them"
    }
    Write-Info "Your code is ready for GitHub Actions"
    exit 0
} else {
    Write-Error " $script:FailureCount check(s) failed"
    if ($script:WarningCount -gt 0) {
        Write-Warning "$script:WarningCount warning(s) found"
    }
    Write-Info "Please fix the issues above before pushing"
    exit 1
}