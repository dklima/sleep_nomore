# Sleep No More

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-Windows-0078D4.svg)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/badge/release-v1.0.0-brightgreen.svg)](https://github.com/dklima/sleep_nomore/releases)

[![Tests](https://github.com/dklima/sleep_nomore/actions/workflows/test.yml/badge.svg)](https://github.com/dklima/sleep_nomore/actions/workflows/test.yml)
[![Release](https://github.com/dklima/sleep_nomore/actions/workflows/release.yml/badge.svg)](https://github.com/dklima/sleep_nomore/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dklima/sleep_nomore)](https://goreportcard.com/report/github.com/dklima/sleep_nomore)
[![codecov](https://codecov.io/gh/dklima/sleep_nomore/branch/main/graph/badge.svg)](https://codecov.io/gh/dklima/sleep_nomore)

A Windows application that keeps your computer awake by preventing it from entering sleep mode.

## Features

- ✅ System tray icon
- ✅ Red icon when disabled
- ✅ Green icon when enabled
- ✅ Single click to toggle awake mode
- ✅ Right-click context menu
- ✅ Single binary with no external dependencies
- ✅ Support for x86_64 and ARM64 architectures

## Usage

1. Run `sleepnomore_amd64.exe` (for Intel/AMD processors) or `sleepnomore_arm64.exe` (for ARM processors)
2. The icon will appear in the system tray (near the clock)
3. Click the icon to toggle awake mode on/off
4. Right-click to access the context menu

## Building from Source

### Prerequisites

- Go 1.21 or higher

### Windows

```batch
build.bat
```

### Linux/Mac (cross-compile for Windows)

```bash
chmod +x build.sh
./build.sh
```

## Output Files

- `sleepnomore_amd64.exe` - Windows x86_64 (Intel/AMD)
- `sleepnomore_arm64.exe` - Windows ARM64

## How It Works

The application uses the Windows API `SetThreadExecutionState` to prevent the system from entering sleep mode, keeping both the system and display active when the mode is enabled.

## License

MIT License - feel free to use this project for any purpose.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
