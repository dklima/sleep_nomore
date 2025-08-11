# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2025-08-11

### Fixed

- Fixed critical bug where system would not return to sleep mode after selecting "Disable" from the menu
- Improved `allowSleep()` function to properly clear execution state using dual Windows API calls
- Enhanced error handling with more descriptive error messages for sleep state operations

### Technical Improvements

- Added proper state reset using `SetThreadExecutionState(0)` after clearing continuous state
- Fixed Go version compatibility in GitHub Actions workflows (1.24 → 1.21)
- Added `.gitattributes` to enforce consistent line endings across development environments
- Normalized line endings in CI/CD workflows to prevent formatting issues

## [1.0.0] - 2025-08-09

### Added

- Initial release of Sleep No More
- System tray icon with visual status indicators (red when disabled, green when enabled)
- Single-click toggle functionality to enable/disable awake mode
- Right-click context menu with Enable/Disable and Quit options
- Windows API integration using `SetThreadExecutionState` to prevent sleep
- Support for both x86_64 and ARM64 architectures
- Embedded icon resources (no external files required)
- Automatic cleanup on application exit
- Signal handling for graceful shutdown (Ctrl+C, SIGTERM)
- Fallback icon generation if embedded resources fail to load
- English language interface
- MIT License

### Technical Details

- Written in Go 1.21
- Uses `github.com/getlantern/systray` for system tray functionality
- Single binary distribution with no external dependencies
- Optimized binary size with `-s -w` linker flags
- Windows GUI subsystem to prevent console window

### Platform Support

- Windows 10 (version 1809 or later)
- Windows 11
- Windows Server 2019 or later

[Unreleased]: https://github.com/yourusername/sleep-no-more/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/yourusername/sleep-no-more/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/yourusername/sleep-no-more/releases/tag/v1.0.0
