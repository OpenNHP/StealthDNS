package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.svg
var icon []byte

// Global app reference for menu callbacks
var globalApp *App

func main() {
	// Create application instance
	app := NewApp()
	globalApp = app

	// Create tray manager (for status sync)
	trayManager := NewTrayManager(app, icon)
	app.SetTrayManager(trayManager)

	// Create application menu
	appMenu := createAppMenu()

	// Only use frameless mode on Windows
	isFrameless := runtime.GOOS == "windows"

	// Create Wails application
	err := wails.Run(&options.App{
		Title:     "StealthDNS",
		Width:     880,
		Height:    720,
		MinWidth:  780,
		MinHeight: 640,
		// Frameless window mode (Windows only, remove system window border)
		Frameless: isFrameless,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 18, B: 24, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnDomReady:       app.onDomReady,
		Bind: []interface{}{
			app,
		},
		// Hide window on close button click instead of quitting
		HideWindowOnClose: true,
		// Application menu
		Menu: appMenu,
		// Single instance lock: ensure only one instance runs, show existing window on second launch
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.opennhp.stealthdns",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				// When second instance tries to launch, show existing window
				if globalApp != nil && globalApp.ctx != nil {
					wailsRuntime.WindowShow(globalApp.ctx)
					wailsRuntime.WindowUnminimise(globalApp.ctx)
					wailsRuntime.WindowSetAlwaysOnTop(globalApp.ctx, true)
					wailsRuntime.WindowSetAlwaysOnTop(globalApp.ctx, false)
				}
			},
		},
		// Windows specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		// Mac specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
			},
			About: &mac.AboutInfo{
				Title:   "StealthDNS",
				Message: "DNS Proxy Management Tool\nVersion 1.0.0",
				Icon:    icon,
			},
		},
		// Linux specific options
		Linux: &linux.Options{
			Icon:                icon,
			WindowIsTranslucent: false,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

// createAppMenu creates application menu
func createAppMenu() *menu.Menu {
	appMenu := menu.NewMenu()

	// macOS application menu
	appMenu.Append(menu.AppMenu())

	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Show Window", keys.CmdOrCtrl("w"), func(cd *menu.CallbackData) {
		if globalApp != nil && globalApp.ctx != nil {
			wailsRuntime.WindowShow(globalApp.ctx)
		}
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) {
		if globalApp != nil {
			globalApp.Quit()
		}
	})

	// Service menu
	serviceMenu := appMenu.AddSubmenu("Service")
	serviceMenu.AddText("Start Service", keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {
		if globalApp != nil {
			go globalApp.StartDNS()
		}
	})
	serviceMenu.AddText("Stop Service", keys.CmdOrCtrl("t"), func(cd *menu.CallbackData) {
		if globalApp != nil {
			go globalApp.StopDNS()
		}
	})
	serviceMenu.AddText("Restart Service", keys.CmdOrCtrl("r"), func(cd *menu.CallbackData) {
		if globalApp != nil {
			go globalApp.RestartDNS()
		}
	})

	// Edit menu
	appMenu.Append(menu.EditMenu())

	return appMenu
}
