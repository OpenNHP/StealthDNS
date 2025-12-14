import { useState, useEffect } from 'react'
import { ServiceStatus } from '../App'
import { useLanguage } from '../contexts/LanguageContext'
import '../styles/StatusPanel.css'

interface SystemDNSInfo {
  dnsServers: string[]
  listenPort: number
  isProxyActive: boolean
}

interface StatusPanelProps {
  status: ServiceStatus
  autoRestart: boolean
  loading: boolean
  onStart: () => void
  onStop: () => void
  onRestart: () => void
  onAutoRestartToggle: () => void
}

export function StatusPanel({ 
  status, 
  autoRestart, 
  loading, 
  onStart, 
  onStop, 
  onRestart, 
  onAutoRestartToggle 
}: StatusPanelProps) {
  const { t } = useLanguage()
  const [dnsInfo, setDnsInfo] = useState<SystemDNSInfo>({
    dnsServers: [],
    listenPort: 53,
    isProxyActive: false
  })

  // Fetch system DNS information
  useEffect(() => {
    const fetchDNSInfo = async () => {
      try {
        const info = await window.go.main.App.GetSystemDNS()
        setDnsInfo(info)
      } catch (err) {
        console.error('Failed to get DNS info:', err)
      }
    }

    fetchDNSInfo()
    // Refresh every 1 second
    const interval = setInterval(fetchDNSInfo, 1000)
    return () => clearInterval(interval)
  }, [status.running])

  // Get localized message based on status code
  const getStatusMessage = (messageCode: string): string => {
    const statusMessages: Record<string, string> = {
      'status_running': t.status.messageRunning,
      'status_stopped': t.status.messageStopped,
      'status_started': t.status.messageStarted,
      'status_starting': t.status.messageStarting,
      'status_stopping': t.status.messageStopping,
    }
    return statusMessages[messageCode] || messageCode
  }

  return (
    <div className="status-panel">
      {/* Status indicator */}
      <div className="status-indicator">
        <div className={`status-orb ${status.running ? 'running' : 'stopped'}`}>
          <div className="orb-inner"></div>
          <div className="orb-pulse"></div>
        </div>
        <div className="status-info">
          <h2 className="status-title">
            {status.running ? t.status.serviceRunning : t.status.serviceStopped}
          </h2>
          <p className="status-message">{getStatusMessage(status.message)}</p>
          {status.restartCount > 0 && (
            <p className="restart-count">{t.status.restartCount}: {status.restartCount}</p>
          )}
        </div>
      </div>

      {/* Control buttons */}
      <div className="control-buttons">
        {!status.running ? (
          <>
            <button 
              className="btn btn-primary btn-large"
              onClick={onStart}
              disabled={loading || status.message === 'status_starting'}
            >
              <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                <path d="M8 5v14l11-7z"/>
              </svg>
              {(loading || status.message === 'status_starting') ? t.status.starting : t.status.startService}
            </button>
            <span className="admin-hint">{t.status.adminRequired}</span>
          </>
        ) : (
          <>
            <button 
              className="btn btn-danger"
              onClick={onStop}
              disabled={loading}
            >
              <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                <path d="M6 6h12v12H6z"/>
              </svg>
              {loading ? t.status.stopping : t.status.stopService}
            </button>
            <button 
              className="btn btn-secondary"
              onClick={onRestart}
              disabled={loading}
            >
              <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
              </svg>
              {t.status.restartService}
            </button>
          </>
        )}
      </div>

      {/* Auto restart settings */}
      <div className="settings-section">
        <h3 className="section-title">{t.tabs.settings}</h3>
        <div className="setting-item">
          <div className="setting-info">
            <span className="setting-label">{t.status.autoRestart}</span>
            <span className="setting-desc">{t.status.autoRestartDesc}</span>
          </div>
          <label className="toggle-switch">
            <input 
              type="checkbox" 
              checked={autoRestart}
              onChange={onAutoRestartToggle}
            />
            <span className="toggle-slider"></span>
          </label>
        </div>
      </div>

      {/* Status information cards */}
      <div className="info-cards info-cards-row">
        <div className="info-card info-card-dns">
          <div className="info-icon">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
              <path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM19 18H6c-2.21 0-4-1.79-4-4s1.79-4 4-4h.71C7.37 7.69 9.48 6 12 6c3.04 0 5.5 2.46 5.5 5.5v.5H19c1.66 0 3 1.34 3 3s-1.34 3-3 3z"/>
            </svg>
          </div>
          <div className="info-content">
            <span className="info-label">{t.status.currentDNS}</span>
            <div className="dns-list">
              {dnsInfo.dnsServers.length > 0 ? (
                dnsInfo.dnsServers.slice(0, 2).map((dns, index) => (
                  <span key={index} className="dns-item">{dns}</span>
                ))
              ) : (
                <span className="dns-item">{t.status.loading}</span>
              )}
              {dnsInfo.dnsServers.length > 2 && (
                <span className="dns-item dns-more">+{dnsInfo.dnsServers.length - 2}</span>
              )}
            </div>
            <span className="info-sub">{t.status.listenPort}: {dnsInfo.listenPort}</span>
          </div>
        </div>
        <div className="info-card info-card-nhp">
          <div className="info-icon">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
              <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"/>
            </svg>
          </div>
          <div className="info-content">
            <span className="info-label">{t.status.nhpDomain}</span>
            <span className="info-value">*.nhp</span>
          </div>
        </div>
      </div>
    </div>
  )
}
