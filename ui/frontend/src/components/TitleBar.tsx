import { useState, useEffect } from 'react'
import { useLanguage } from '../contexts/LanguageContext'
import '../styles/TitleBar.css'

interface TitleBarProps {
  onMinimize: () => void
  onQuit: () => void
}

// Detect if running on macOS
function isMacOS(): boolean {
  return navigator.platform.toUpperCase().indexOf('MAC') >= 0 || 
         navigator.userAgent.toUpperCase().indexOf('MAC') >= 0
}

export function TitleBar({ onMinimize, onQuit }: TitleBarProps) {
  const { t } = useLanguage()
  const [isMac, setIsMac] = useState(false)
  const [isMaximized, setIsMaximized] = useState(false)
  const [version, setVersion] = useState<string>(t.app.version)

  useEffect(() => {
    setIsMac(isMacOS())
    // Get actual version from backend (includes build number if available)
    const fetchVersion = async () => {
      try {
        if (window.go?.main?.App?.GetVersion) {
          const v = await window.go.main.App.GetVersion()
          // Get the version prefix from i18n (e.g., "版本", "Version", "バージョン")
          const versionPrefix = t.app.version.split(' ')[0]
          // Display: "版本 1.0.0+123" or "Version 1.0.0+123"
          setVersion(`${versionPrefix} ${v}`)
        }
      } catch (err) {
        console.error('Failed to get version:', err)
        // Fallback to i18n version if GetVersion fails
        setVersion(t.app.version)
      }
    }
    fetchVersion()
  }, [t.app.version])

  // Toggle maximize/restore
  const handleToggleMaximize = () => {
    if (isMaximized) {
      window.runtime.WindowUnmaximise()
    } else {
      window.runtime.WindowMaximise()
    }
    setIsMaximized(!isMaximized)
  }

  return (
    <div 
      className={`title-bar ${isMac ? 'title-bar-macos' : ''}`} 
      style={{ '--wails-draggable': 'drag' } as React.CSSProperties}
    >
      <div className="title-bar-logo">
        <svg viewBox="0 0 100 100" width="24" height="24">
          <defs>
            <linearGradient id="shield-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#00d4aa" />
              <stop offset="100%" stopColor="#00a080" />
            </linearGradient>
          </defs>
          <path 
            d="M50 5 L90 20 L90 45 C90 70 70 90 50 95 C30 90 10 70 10 45 L10 20 Z" 
            fill="url(#shield-gradient)" 
            opacity="0.9"
          />
          <path 
            d="M50 15 L80 27 L80 45 C80 65 65 80 50 85 C35 80 20 65 20 45 L20 27 Z" 
            fill="none" 
            stroke="#fff" 
            strokeWidth="2" 
            opacity="0.6"
          />
          <text x="50" y="58" textAnchor="middle" fill="#fff" fontSize="24" fontWeight="bold" fontFamily="monospace">D</text>
        </svg>
        <span className="title-bar-text">{t.app.title}</span>
        <span className="title-bar-version">{version}</span>
      </div>
      
      <div className="title-bar-controls">
        <button 
          className="title-bar-btn minimize"
          onClick={onMinimize}
          title={t.titleBar.minimize}
        >
          <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
            <path d="M19 13H5v-2h14v2z"/>
          </svg>
        </button>
        {!isMac && (
          <button 
            className="title-bar-btn maximize"
            onClick={handleToggleMaximize}
            title={isMaximized ? t.titleBar.restore : t.titleBar.maximize}
          >
            {isMaximized ? (
              <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
                <path d="M4 8h4V4h12v12h-4v4H4V8zm2 2v8h8v-8H6zm10-4v2h2v8h-2v2h4V6h-4z"/>
              </svg>
            ) : (
              <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
                <path d="M4 4h16v16H4V4zm2 2v12h12V6H6z"/>
              </svg>
            )}
          </button>
        )}
        <button 
          className="title-bar-btn close"
          onClick={onQuit}
          title={t.titleBar.quit}
        >
          <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
          </svg>
        </button>
      </div>
    </div>
  )
}
