# StealthDNS Mobile Browser

A mobile browser application for StealthDNS, built with native development (Android: Kotlin, iOS: Swift), integrating Go language for NHP protocol core logic.

## Features

- 🛡️ **NHP Protocol Support**: Automatically identifies `.nhp` domains and performs NHP knock operations
- 🌐 **Native WebView**: Android uses WebView, iOS uses WKWebView
- 🔐 **Zero Trust Access**: Secure network access through NHP protocol
- 📱 **Native Performance**: Built with Kotlin and Swift for excellent performance
- 🔧 **Go Core Library**: DNS resolution and NHP knock logic implemented in Go, compiled via gomobile

## Project Structure

```
mobile/
├── golib/                    # Go NHP core library (gomobile compiled)
│   ├── nhpcore.go            # NHP core logic
│   └── go.mod                # Go module configuration
├── android/                  # Android native project (Kotlin)
│   ├── app/
│   │   ├── src/main/
│   │   │   ├── java/.../     # Kotlin source code
│   │   │   │   ├── MainActivity.kt
│   │   │   │   └── StealthDNSApp.kt
│   │   │   ├── res/          # Resource files
│   │   │   └── AndroidManifest.xml
│   │   ├── libs/             # Place nhpcore.aar here
│   │   └── build.gradle.kts
│   └── build.gradle.kts
├── ios/                      # iOS native project (Swift)
│   ├── StealthDNS/
│   │   ├── AppDelegate.swift
│   │   ├── BrowserViewController.swift
│   │   ├── NhpcoreBridge.swift
│   │   ├── Info.plist
│   │   └── Assets.xcassets/
│   └── StealthDNS.xcodeproj/
├── assets/                   # Shared assets
│   └── resources.json        # NHP resource configuration
├── scripts/
│   └── build.sh              # Build script
├── Makefile                  # Make build commands
└── README.md
```

## Requirements

### Go Library Compilation
- Go 1.21+
- gomobile (`go install golang.org/x/mobile/cmd/gomobile@latest`)

### Android Build
- **JDK 17** (recommended) or JDK 11+ (required)
- Android Studio or Android SDK
- Android SDK (API 21+)
- Android NDK

> ⚠️ **Important**: Android Gradle Plugin 8.x requires **Java 11 or higher**.
> 
> Install JDK 17 (recommended):
> ```bash
> # macOS (Homebrew)
> brew install openjdk@17
> 
> # Set JAVA_HOME
> export JAVA_HOME=$(/usr/libexec/java_home -v 17)
> ```

### iOS Build (macOS only)
- Xcode 14+
- iOS 14.0+ target

## Quick Start

### 1. Install gomobile

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
gomobile init
```

### 2. Build Go Library

```bash
cd mobile

# Build Android library
make golib-android

# Build iOS library (macOS only)
make golib-ios
```

### 3. Build Application

```bash
# Android APK
make android

# iOS (build after opening in Xcode)
make ios
```

Or use the script:

```bash
./scripts/build.sh --all --release
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make golib` | Build Go library for all platforms |
| `make golib-android` | Build Android AAR |
| `make golib-ios` | Build iOS XCFramework |
| `make android` | Build Android APK (release) |
| `make android-debug` | Build Android APK (debug) |
| `make android-bundle` | Build Android App Bundle |
| `make ios` | Build iOS archive |
| `make clean` | Clean build artifacts |

## Output Files

- **Android APK**: `build/android/StealthDNS-Browser.apk`
- **Android AAB**: `build/android/StealthDNS-Browser.aab`
- **iOS Archive**: `build/ios/StealthDNS.xcarchive`
- **Go Android Library**: `build/golib/nhpcore.aar`
- **Go iOS Library**: `build/golib/Nhpcore.xcframework`

## Architecture

### Go NHP Core Library (golib/)

Core functionality implemented in Go, compiled to mobile platform libraries via gomobile:

```go
// Main APIs
nhpcore.Initialize(workDir, logLevel, upstreamDNS)  // Initialize
nhpcore.AddResource(...)                             // Add resource configuration
nhpcore.IsNHPDomain(domain)                          // Check if NHP domain
nhpcore.Knock(resourceId)                            // Perform NHP knock
nhpcore.GetKnockResultJSON(resourceId)               // Get knock result
nhpcore.Cleanup()                                    // Cleanup resources
```

### Android (Kotlin)

- `MainActivity.kt`: Main interface with WebView and navigation bar
- Calls Go functions in `nhpcore.aar` to handle NHP domains
- Uses Android WebView for rendering web pages

### iOS (Swift)

- `BrowserViewController.swift`: Main interface with WKWebView
- `NhpcoreBridge.swift`: Swift bridge for Go library
- Uses WKWebView for rendering web pages

## NHP Workflow

1. User enters URL (e.g., `https://demo.nhp`)
2. Detects `.nhp` suffix, identifies as NHP domain
3. Extracts resource ID (`demo`)
4. Calls Go library to perform NHP knock:
   - Sends authentication request to NHP server
   - Server validates and returns actual resource IP
5. Replaces domain with returned IP, loads actual page
6. Result is cached, no need to re-knock within validity period

## Resource Configuration

Edit `assets/resources.json`:

```json
[
  {
    "authServiceId": "your-auth-service",
    "resourceId": "demo",
    "serverIp": "",
    "serverHostname": "nhp.example.com",
    "serverPort": 62206
  }
]
```

## Release Signing

### Android

1. Generate signing key:
```bash
keytool -genkey -v -keystore stealthdns.jks -keyalg RSA -keysize 2048 -validity 10000 -alias stealthdns
```

2. Set environment variables:
```bash
export ANDROID_KEYSTORE_PATH=/path/to/stealthdns.jks
export ANDROID_KEYSTORE_PASSWORD=your_password
export ANDROID_KEY_ALIAS=stealthdns
export ANDROID_KEY_PASSWORD=your_key_password
```

3. Build release version:
```bash
make android
```

### iOS

1. Open `ios/StealthDNS.xcodeproj` in Xcode
2. Configure developer certificate and Provisioning Profile
3. Product → Archive → Distribute App

## Troubleshooting

### gomobile Errors

```bash
# Reinitialize
gomobile init

# Check NDK
echo $ANDROID_NDK_HOME
```

### Android Build Errors

```bash
# Clean and rebuild
cd android
./gradlew clean
cd ..
make android
```

### iOS Build Errors

```bash
# Regenerate Xcode configuration
rm -rf ios/build
make ios
```

## License

Apache 2.0 License

