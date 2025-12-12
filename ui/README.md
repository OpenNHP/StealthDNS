# StealthDNS UI

Cross-platform StealthDNS graphical interface management tool, built with Go + Wails + React.

## Features

- 🖥️ **Cross-platform Support**: Windows / macOS / Linux
- 🔄 **Process Management**: Start, stop, restart DNS proxy service through UI
- 🛡️ **Automatic Crash Recovery**: Automatically restart when DNS service crashes
- ⚙️ **Configuration Management**: Visual configuration of NHP server address, public key and client private key
- 📥 **System Tray**: Minimize to system tray for background operation

## Prerequisites

### Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Install Node.js

Requires Node.js 16+ version, recommend using [nvm](https://github.com/nvm-sh/nvm) to manage Node versions.

### Platform-specific Dependencies

#### Windows
- WebView2 Runtime (usually pre-installed on Windows 10/11)

#### macOS
- Xcode Command Line Tools

#### Linux
- GTK3, WebKit2GTK

Ubuntu/Debian:
```bash
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev
```

Fedora:
```bash
sudo dnf install gtk3-devel webkit2gtk3-devel
```

Arch:
```bash
sudo pacman -S gtk3 webkit2gtk
```

## Build

### Using Makefile (Recommended)

```bash
# Build UI only
make ui

# Build complete package (DNS service + UI)
make full

# Development mode (hot reload)
make ui-dev
```

### Build directly with Wails CLI

```bash
cd ui

# Development mode
wails dev

# Production build
wails build
```

### Windows Batch File

```batch
:: Build UI only
build.bat ui

:: Build complete package
build.bat full
```

## Directory Structure

```
ui/
├── app.go              # Backend application logic (process management, configuration management)
├── main.go             # Wails application entry point
├── tray.go             # System tray functionality
├── wails.json          # Wails configuration file
├── go.mod              # Go module definition
├── build/
│   └── appicon.svg     # Application icon
└── frontend/           # Frontend code
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    └── src/
        ├── App.tsx     # Main application component
        ├── main.tsx    # React entry point
        ├── components/ # UI components
        │   ├── TitleBar.tsx
        │   ├── StatusPanel.tsx
        │   ├── ConfigPanel.tsx
        │   └── ServerPanel.tsx
        └── styles/     # CSS styles
```

## Usage

1. **Start Application**: Run the `stealthdns-ui` executable
2. **Start DNS Service**: Click "Start Service" button (requires admin privileges)
3. **Configure Client**: Set private key and other parameters in "Client Configuration" tab
4. **Configure Server**: Add/edit NHP servers in "Server Configuration" tab
5. **Minimize to Tray**: Click the minimize button in the title bar

## Notes

- StealthDNS DNS service requires admin/root privileges to listen on port 53
- UI program and `stealth-dns` executable need to be in the same directory
- Configuration files are located in the `etc/` directory

## Development

### Hot Reload Development

```bash
cd ui
wails dev
```

Frontend changes will hot reload automatically, backend changes require recompilation.

### Generate Bindings

When Go backend API changes, regenerate TypeScript bindings:

```bash
cd ui
wails generate module
```

