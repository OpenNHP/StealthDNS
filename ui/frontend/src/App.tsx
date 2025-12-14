import { useState, useEffect, useCallback } from 'react'
import { StatusPanel } from './components/StatusPanel'
import { ConfigPanel } from './components/ConfigPanel'
import { ServerPanel } from './components/ServerPanel'
import { SettingsPanel } from './components/SettingsPanel'
import { LogPanel } from './components/LogPanel'
import { TitleBar } from './components/TitleBar'
import { LanguageProvider, useLanguage } from './contexts/LanguageContext'
import './styles/App.css'

// Log related types
interface LogEntry {
  timestamp: string
  level: string
  message: string
  raw: string
}

interface LogFile {
  name: string
  size: number
  modified: string
}

interface WatchLogResult {
  entries: LogEntry[]
  size: number
}

// System DNS information type
interface SystemDNSInfo {
  dnsServers: string[]
  listenPort: number
  isProxyActive: boolean
}

interface VersionInfo {
  version: string
  buildNumber: string
  commitId: string
  buildTime: string
}

// Wails runtime types
declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetStatus: () => Promise<ServiceStatus>
          StartDNS: () => Promise<void>
          StopDNS: () => Promise<void>
          RestartDNS: () => Promise<void>
          GetAutoRestart: () => Promise<boolean>
          SetAutoRestart: (enabled: boolean) => Promise<void>
          GetClientConfig: () => Promise<ClientConfig>
          SaveClientConfig: (config: ClientConfig) => Promise<void>
          GetServerConfig: () => Promise<ServerConfig[]>
          SaveServerConfig: (servers: ServerConfig[]) => Promise<void>
          MinimizeToTray: () => Promise<void>
          Quit: () => Promise<void>
          // Log related methods
          GetLogFiles: () => Promise<LogFile[]>
          GetLogContent: (filename: string, lines: number) => Promise<LogEntry[]>
          WatchLogFile: (filename: string, lastSize: number) => Promise<WatchLogResult>
          ClearLogFile: (filename: string) => Promise<void>
          // System DNS information
          GetSystemDNS: () => Promise<SystemDNSInfo>
          GetVersion: () => Promise<string>
          GetVersionInfo: () => Promise<VersionInfo>
        }
      }
    }
    runtime: {
      EventsOn: (eventName: string, callback: (...data: unknown[]) => void) => () => void
      EventsOff: (eventName: string) => void
      WindowMinimise: () => void
      WindowMaximise: () => void
      WindowUnmaximise: () => void
      WindowHide: () => void
      Quit: () => void
    }
  }
}

export interface ServiceStatus {
  running: boolean
  restartCount: number
  message: string
}

export interface ClientConfig {
  privateKeyBase64: string
  defaultCipherScheme: number
  userId: string
  organizationId: string
  logLevel: number
}

export interface ServerConfig {
  hostname: string
  ip: string
  port: number
  pubKeyBase64: string
  expireTime: number
}

type TabType = 'status' | 'config' | 'server' | 'logs' | 'settings'

function AppContent() {
  const { t } = useLanguage()
  const [activeTab, setActiveTab] = useState<TabType>('status')
  const [status, setStatus] = useState<ServiceStatus>({
    running: false,
    restartCount: 0,
    message: ''
  })
  const [autoRestart, setAutoRestart] = useState(true)
  const [clientConfig, setClientConfig] = useState<ClientConfig | null>(null)
  const [serverConfig, setServerConfig] = useState<ServerConfig[]>([])
  const [loading, setLoading] = useState(false)
  const [notification, setNotification] = useState<{ type: 'success' | 'error', message: string } | null>(null)

  // Show notification
  const showNotification = useCallback((type: 'success' | 'error', message: string) => {
    setNotification({ type, message })
    setTimeout(() => setNotification(null), 3000)
  }, [])

  // Load status
  const loadStatus = useCallback(async () => {
    try {
      const s = await window.go.main.App.GetStatus()
      setStatus(s)
      const ar = await window.go.main.App.GetAutoRestart()
      setAutoRestart(ar)
    } catch (err) {
      console.error('Load status failed:', err)
    }
  }, [])

  // Load configuration
  const loadConfig = useCallback(async () => {
    try {
      const cc = await window.go.main.App.GetClientConfig()
      setClientConfig(cc)
      const sc = await window.go.main.App.GetServerConfig()
      setServerConfig(sc)
    } catch (err) {
      console.error('Load config failed:', err)
    }
  }, [])

  useEffect(() => {
    loadStatus()
    loadConfig()

    // Listen for status change events
    const unsubscribe = window.runtime.EventsOn('dns:status', (data: unknown) => {
      setStatus(data as ServiceStatus)
    })

    // Periodically refresh status to ensure UI is in sync
    // GetStatus() now automatically checks actual process state and syncs internal state
    const refreshStatus = async () => {
      try {
        const status = await window.go.main.App.GetStatus()
        setStatus(status)
      } catch (err) {
        console.error('Failed to refresh status:', err)
      }
    }

    // Refresh immediately once
    refreshStatus()
    // Refresh status every 1 second to ensure UI stays in sync
    const statusRefreshInterval = setInterval(refreshStatus, 1000)

    return () => {
      unsubscribe()
      clearInterval(statusRefreshInterval)
    }
  }, [loadStatus, loadConfig])

  // Clear loading state when service status changes to running or stopped
  useEffect(() => {
    if (status.running || status.message === 'status_stopped') {
      setLoading(false)
    }
  }, [status.running, status.message])

  // Start service
  const handleStart = async () => {
    setLoading(true)
    try {
      await window.go.main.App.StartDNS()
      // Don't set loading = false immediately
      // The status will be updated to "starting" and then "running" when the process is detected
      // The button will show "starting..." based on status.message === 'status_starting'
      // Only set loading = false if there's an error
    } catch (err) {
      setLoading(false)
      showNotification('error', `${t.messages.startFailed}: ${err}`)
    }
    // Note: loading will be set to false when status changes to "running" or "stopped"
  }

  // Stop service
  const handleStop = async () => {
    setLoading(true)
    try {
      await window.go.main.App.StopDNS()
      showNotification('success', t.messages.serviceStopped)
    } catch (err) {
      showNotification('error', `${t.messages.stopFailed}: ${err}`)
    } finally {
      setLoading(false)
    }
  }

  // Restart service
  const handleRestart = async () => {
    setLoading(true)
    try {
      await window.go.main.App.RestartDNS()
      showNotification('success', t.messages.serviceRestarted)
    } catch (err) {
      showNotification('error', `${t.messages.restartFailed}: ${err}`)
    } finally {
      setLoading(false)
    }
  }

  // Toggle auto restart
  const handleAutoRestartToggle = async () => {
    try {
      await window.go.main.App.SetAutoRestart(!autoRestart)
      setAutoRestart(!autoRestart)
    } catch (err) {
      showNotification('error', `${t.messages.saveFailed}: ${err}`)
    }
  }

  // Save client configuration
  const handleSaveClientConfig = async (config: ClientConfig) => {
    setLoading(true)
    try {
      await window.go.main.App.SaveClientConfig(config)
      setClientConfig(config)
      showNotification('success', t.messages.configSaved)
    } catch (err) {
      showNotification('error', `${t.messages.saveFailed}: ${err}`)
    } finally {
      setLoading(false)
    }
  }

  // Save server configuration
  const handleSaveServerConfig = async (servers: ServerConfig[]) => {
    setLoading(true)
    try {
      await window.go.main.App.SaveServerConfig(servers)
      setServerConfig(servers)
      showNotification('success', t.messages.serverConfigSaved)
    } catch (err) {
      showNotification('error', `${t.messages.saveFailed}: ${err}`)
    } finally {
      setLoading(false)
    }
  }

  // Minimize to tray
  const handleMinimize = async () => {
    try {
      await window.go.main.App.MinimizeToTray()
    } catch (err) {
      console.error('Minimize failed:', err)
    }
  }

  // Quit application
  const handleQuit = async () => {
    try {
      await window.go.main.App.Quit()
    } catch (err) {
      console.error('Quit failed:', err)
    }
  }

  return (
    <div className="app">
      <TitleBar onMinimize={handleMinimize} onQuit={handleQuit} />
      
      {/* Notification */}
      {notification && (
        <div className={`notification ${notification.type}`}>
          {notification.message}
        </div>
      )}

      {/* Tab navigation */}
      <nav className="tabs">
        <button 
          className={`tab ${activeTab === 'status' ? 'active' : ''}`}
          onClick={() => setActiveTab('status')}
        >
          <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
          {t.tabs.status}
        </button>
        <button 
          className={`tab ${activeTab === 'config' ? 'active' : ''}`}
          onClick={() => setActiveTab('config')}
        >
          <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
            <path d="M12.87 15.07l-2.54-2.51.03-.03c1.74-1.94 2.98-4.17 3.71-6.53H17V4h-7V2H8v2H1v1.99h11.17C11.5 7.92 10.44 9.75 9 11.35 8.07 10.32 7.3 9.19 6.69 8h-2c.73 1.63 1.73 3.17 2.98 4.56l-5.09 5.02L4 19l5-5 3.11 3.11.76-2.04zM18.5 10h-2L12 22h2l1.12-3h4.75L21 22h2l-4.5-12zm-2.62 7l1.62-4.33L19.12 17h-3.24z"/>
          </svg>
          {t.tabs.clientConfig}
        </button>
        <button 
          className={`tab ${activeTab === 'server' ? 'active' : ''}`}
          onClick={() => setActiveTab('server')}
        >
          <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
            <path d="M19 15v4H5v-4h14m1-2H4c-.55 0-1 .45-1 1v6c0 .55.45 1 1 1h16c.55 0 1-.45 1-1v-6c0-.55-.45-1-1-1zM7 18.5c-.82 0-1.5-.67-1.5-1.5s.68-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zM19 5v4H5V5h14m1-2H4c-.55 0-1 .45-1 1v6c0 .55.45 1 1 1h16c.55 0 1-.45 1-1V4c0-.55-.45-1-1-1zM7 8.5c-.82 0-1.5-.67-1.5-1.5S6.18 5.5 7 5.5s1.5.68 1.5 1.5S7.83 8.5 7 8.5z"/>
          </svg>
          {t.tabs.serverConfig}
        </button>
        <button 
          className={`tab ${activeTab === 'logs' ? 'active' : ''}`}
          onClick={() => setActiveTab('logs')}
        >
          <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
            <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/>
          </svg>
          {t.tabs.logs}
        </button>
        <button 
          className={`tab ${activeTab === 'settings' ? 'active' : ''}`}
          onClick={() => setActiveTab('settings')}
        >
          <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
            <path d="M19.14 12.94c.04-.31.06-.63.06-.94 0-.31-.02-.63-.06-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.04.31-.06.63-.06.94s.02.63.06.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/>
          </svg>
          {t.tabs.settings}
        </button>
      </nav>

      {/* Content area */}
      <main className="content">
        {activeTab === 'status' && (
          <StatusPanel
            status={status}
            autoRestart={autoRestart}
            loading={loading}
            onStart={handleStart}
            onStop={handleStop}
            onRestart={handleRestart}
            onAutoRestartToggle={handleAutoRestartToggle}
          />
        )}
        {activeTab === 'config' && clientConfig && (
          <ConfigPanel
            config={clientConfig}
            loading={loading}
            onSave={handleSaveClientConfig}
          />
        )}
        {activeTab === 'server' && (
          <ServerPanel
            servers={serverConfig}
            loading={loading}
            onSave={handleSaveServerConfig}
          />
        )}
        {activeTab === 'logs' && (
          <LogPanel />
        )}
        {activeTab === 'settings' && (
          <SettingsPanel />
        )}
      </main>
    </div>
  )
}

function App() {
  return (
    <LanguageProvider>
      <AppContent />
    </LanguageProvider>
  )
}

export default App
