package main

// TrayManager tray status manager (simplified version, no system tray library dependency)
// Wails' HideWindowOnClose option already implements basic tray functionality
type TrayManager struct {
	app     *App
	running bool
}

// NewTrayManager creates tray manager
func NewTrayManager(app *App, iconData []byte) *TrayManager {
	return &TrayManager{
		app:     app,
		running: false,
	}
}

// Run starts tray manager
func (t *TrayManager) Run() {
	// Wails' built-in HideWindowOnClose option handles window hiding
	// No additional system tray implementation needed here
}

// UpdateStatus updates service status
func (t *TrayManager) UpdateStatus(running bool) {
	t.running = running
}
