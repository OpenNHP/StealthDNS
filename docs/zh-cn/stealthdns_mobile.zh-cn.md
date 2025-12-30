# StealthDNS Mobile Browser

StealthDNS 移动端浏览器应用，使用原生开发（Android: Kotlin, iOS: Swift），结合 Go 语言实现 NHP 协议核心逻辑。

## 功能特性

- 🛡️ **NHP 协议支持**: 自动识别 `.nhp` 域名并执行 NHP 敲门操作
- 🌐 **原生 WebView**: Android 使用 WebView，iOS 使用 WKWebView
- 🔐 **零信任访问**: 通过 NHP 协议实现安全的网络访问
- 📱 **原生性能**: 使用 Kotlin 和 Swift 原生开发，性能优异
- 🔧 **Go 核心库**: DNS 解析和 NHP 敲门逻辑使用 Go 实现，通过 gomobile 编译

## 项目结构

```
mobile/
├── golib/                    # Go NHP 核心库（gomobile 编译）
│   ├── nhpcore.go            # NHP 核心逻辑
│   └── go.mod                # Go 模块配置
├── android/                  # Android 原生项目 (Kotlin)
│   ├── app/
│   │   ├── src/main/
│   │   │   ├── java/.../     # Kotlin 源代码
│   │   │   │   ├── MainActivity.kt
│   │   │   │   └── StealthDNSApp.kt
│   │   │   ├── res/          # 资源文件
│   │   │   └── AndroidManifest.xml
│   │   ├── libs/             # 存放 nhpcore.aar
│   │   └── build.gradle.kts
│   └── build.gradle.kts
├── ios/                      # iOS 原生项目 (Swift)
│   ├── StealthDNS/
│   │   ├── AppDelegate.swift
│   │   ├── BrowserViewController.swift
│   │   ├── NhpcoreBridge.swift
│   │   ├── Info.plist
│   │   └── Assets.xcassets/
│   └── StealthDNS.xcodeproj/
├── assets/                   # 共享资源
│   └── resources.json        # NHP 资源配置
├── scripts/
│   └── build.sh              # 构建脚本
├── Makefile                  # Make 构建命令
└── README.md
```

## 环境要求

### Go 库编译
- Go 1.21+
- gomobile (`go install golang.org/x/mobile/cmd/gomobile@latest`)

### Android 构建
- **JDK 17** (推荐) 或 JDK 11+ (必须)
- Android Studio 或 Android SDK
- Android SDK (API 21+)
- Android NDK

> ⚠️ **重要**: Android Gradle Plugin 8.x 需要 **Java 11 或更高版本**。
> 
> 安装 JDK 17 (推荐):
> ```bash
> # macOS (Homebrew)
> brew install openjdk@17
> 
> # 设置 JAVA_HOME
> export JAVA_HOME=$(/usr/libexec/java_home -v 17)
> ```

### iOS 构建 (仅 macOS)
- Xcode 14+
- iOS 14.0+ 目标

## 快速开始

### 1. 安装 gomobile

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
gomobile init
```

### 2. 构建 Go 库

```bash
cd mobile

# 构建 Android 库
make golib-android

# 构建 iOS 库 (仅 macOS)
make golib-ios
```

### 3. 构建应用

```bash
# Android APK
make android

# iOS (使用 Xcode 打开后构建)
make ios
```

或使用脚本：

```bash
./scripts/build.sh --all --release
```

## 构建命令

| 命令 | 说明 |
|------|------|
| `make golib` | 构建所有平台的 Go 库 |
| `make golib-android` | 构建 Android AAR |
| `make golib-ios` | 构建 iOS XCFramework |
| `make android` | 构建 Android APK (release) |
| `make android-debug` | 构建 Android APK (debug) |
| `make android-bundle` | 构建 Android App Bundle |
| `make ios` | 构建 iOS 归档 |
| `make clean` | 清理构建产物 |

## 输出文件

- **Android APK**: `build/android/StealthDNS-Browser.apk`
- **Android AAB**: `build/android/StealthDNS-Browser.aab`
- **iOS Archive**: `build/ios/StealthDNS.xcarchive`
- **Go Android库**: `build/golib/nhpcore.aar`
- **Go iOS库**: `build/golib/Nhpcore.xcframework`

## 架构说明

### Go NHP 核心库 (golib/)

使用 Go 语言实现的核心功能，通过 gomobile 编译为移动平台库：

```go
// 主要 API
nhpcore.Initialize(workDir, logLevel, upstreamDNS)  // 初始化
nhpcore.AddResource(...)                             // 添加资源配置
nhpcore.IsNHPDomain(domain)                          // 检查是否为 NHP 域名
nhpcore.Knock(resourceId)                            // 执行 NHP 敲门
nhpcore.GetKnockResultJSON(resourceId)               // 获取敲门结果
nhpcore.Cleanup()                                    // 清理资源
```

### Android (Kotlin)

- `MainActivity.kt`: 主界面，包含 WebView 和导航栏
- 调用 `nhpcore.aar` 中的 Go 函数处理 NHP 域名
- 使用 Android WebView 渲染网页

### iOS (Swift)

- `BrowserViewController.swift`: 主界面，包含 WKWebView
- `NhpcoreBridge.swift`: Go 库的 Swift 桥接
- 使用 WKWebView 渲染网页

## NHP 工作流程

1. 用户输入 URL（如 `https://demo.nhp`）
2. 检测到 `.nhp` 后缀，识别为 NHP 域名
3. 提取资源 ID（`demo`）
4. 调用 Go 库执行 NHP 敲门：
   - 向 NHP 服务器发送认证请求
   - 服务器验证后返回实际资源 IP
5. 使用返回的 IP 替换域名，加载实际页面
6. 结果被缓存，有效期内无需重复敲门

## 配置资源

编辑 `assets/resources.json`：

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

## 签名发布

### Android

1. 生成签名密钥：
```bash
keytool -genkey -v -keystore stealthdns.jks -keyalg RSA -keysize 2048 -validity 10000 -alias stealthdns
```

2. 设置环境变量：
```bash
export ANDROID_KEYSTORE_PATH=/path/to/stealthdns.jks
export ANDROID_KEYSTORE_PASSWORD=your_password
export ANDROID_KEY_ALIAS=stealthdns
export ANDROID_KEY_PASSWORD=your_key_password
```

3. 构建发布版：
```bash
make android
```

### iOS

1. 在 Xcode 中打开 `ios/StealthDNS.xcodeproj`
2. 配置开发者证书和 Provisioning Profile
3. Product → Archive → Distribute App

## 故障排除

### gomobile 错误

```bash
# 重新初始化
gomobile init

# 检查 NDK
echo $ANDROID_NDK_HOME
```

### Android 构建错误

```bash
# 清理并重建
cd android
./gradlew clean
cd ..
make android
```

### iOS 构建错误

```bash
# 重新生成 Xcode 配置
rm -rf ios/build
make ios
```

## 许可证

Apache 2.0 License
