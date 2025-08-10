# Sleep No More

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