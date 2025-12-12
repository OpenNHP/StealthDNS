import { useState, useEffect } from 'react'
import { ServerConfig } from '../App'
import { useLanguage } from '../contexts/LanguageContext'
import '../styles/ServerPanel.css'

interface ServerPanelProps {
  servers: ServerConfig[]
  loading: boolean
  onSave: (servers: ServerConfig[]) => void
}

export function ServerPanel({ servers, loading, onSave }: ServerPanelProps) {
  const { t, language } = useLanguage()
  const [formData, setFormData] = useState<ServerConfig[]>(servers)
  const [hasChanges, setHasChanges] = useState(false)
  const [expandedIndex, setExpandedIndex] = useState<number | null>(0)

  useEffect(() => {
    setFormData(servers)
    setHasChanges(false)
  }, [servers])

  const handleChange = (index: number, field: keyof ServerConfig, value: string | number) => {
    setFormData(prev => {
      const updated = [...prev]
      updated[index] = { ...updated[index], [field]: value }
      return updated
    })
    setHasChanges(true)
  }

  const handleAddServer = () => {
    const newServer: ServerConfig = {
      hostname: '',
      ip: '',
      port: 62206,
      pubKeyBase64: '',
      expireTime: Math.floor(Date.now() / 1000) + 365 * 24 * 60 * 60 // 1年后过期
    }
    setFormData(prev => [...prev, newServer])
    setExpandedIndex(formData.length)
    setHasChanges(true)
  }

  const handleRemoveServer = (index: number) => {
    if (formData.length <= 1) {
      return // 至少保留一个服务器
    }
    setFormData(prev => prev.filter((_, i) => i !== index))
    setHasChanges(true)
    if (expandedIndex === index) {
      setExpandedIndex(null)
    } else if (expandedIndex !== null && expandedIndex > index) {
      setExpandedIndex(expandedIndex - 1)
    }
  }

  const handleSave = () => {
    onSave(formData)
    setHasChanges(false)
  }

  const handleReset = () => {
    setFormData(servers)
    setHasChanges(false)
  }

  const formatExpireTime = (timestamp: number): string => {
    const date = new Date(timestamp * 1000)
    return date.toLocaleDateString(language, {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }

  return (
    <div className="server-panel">
      <div className="panel-header">
        <div className="header-left">
          <h2>{t.serverConfig.title}</h2>
          <p className="panel-desc">{t.serverConfig.description}</p>
        </div>
        <button className="btn btn-add" onClick={handleAddServer}>
          <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
          {t.serverConfig.addServer}
        </button>
      </div>

      <div className="server-list">
        {formData.map((server, index) => (
          <div 
            key={index} 
            className={`server-card ${expandedIndex === index ? 'expanded' : ''}`}
          >
            <div 
              className="server-card-header"
              onClick={() => setExpandedIndex(expandedIndex === index ? null : index)}
            >
              <div className="server-info">
                <div className="server-icon">
                  <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
                    <path d="M19 15v4H5v-4h14m1-2H4c-.55 0-1 .45-1 1v6c0 .55.45 1 1 1h16c.55 0 1-.45 1-1v-6c0-.55-.45-1-1-1zM7 18.5c-.82 0-1.5-.67-1.5-1.5s.68-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z"/>
                  </svg>
                </div>
                <div className="server-summary">
                  <span className="server-name">
                    {server.hostname || server.ip || t.serverConfig.newServer}
                  </span>
                  <span className="server-port">{t.serverConfig.port}: {server.port}</span>
                </div>
              </div>
              <div className="server-actions">
                {formData.length > 1 && (
                  <button 
                    className="btn-icon btn-remove"
                    onClick={(e) => {
                      e.stopPropagation()
                      handleRemoveServer(index)
                    }}
                    title={t.serverConfig.deleteServer}
                  >
                    <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                      <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                    </svg>
                  </button>
                )}
                <svg 
                  className="expand-icon" 
                  viewBox="0 0 24 24" 
                  width="20" 
                  height="20" 
                  fill="currentColor"
                >
                  <path d="M7.41 8.59L12 13.17l4.59-4.58L18 10l-6 6-6-6 1.41-1.41z"/>
                </svg>
              </div>
            </div>
            
            {expandedIndex === index && (
              <div className="server-card-body">
                <div className="form-row">
                  <div className="form-group">
                    <label className="form-label">
                      <span className="label-text">{t.serverConfig.hostname}</span>
                      <span className="label-hint">{t.serverConfig.hostnameHint}</span>
                    </label>
                    <input
                      type="text"
                      className="form-input"
                      value={server.hostname}
                      onChange={e => handleChange(index, 'hostname', e.target.value)}
                      placeholder="nhp.example.com"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">
                      <span className="label-text">{t.serverConfig.ipAddress}</span>
                      <span className="label-hint">{t.serverConfig.ipAddressHint}</span>
                    </label>
                    <input
                      type="text"
                      className="form-input"
                      value={server.ip}
                      onChange={e => handleChange(index, 'ip', e.target.value)}
                      placeholder="192.168.1.1"
                    />
                  </div>
                </div>

                <div className="form-row">
                  <div className="form-group">
                    <label className="form-label">
                      <span className="label-text">{t.serverConfig.port}</span>
                    </label>
                    <input
                      type="number"
                      className="form-input"
                      value={server.port}
                      onChange={e => handleChange(index, 'port', parseInt(e.target.value) || 0)}
                      placeholder="62206"
                      min="1"
                      max="65535"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">
                      <span className="label-text">{t.serverConfig.expireTime}</span>
                      <span className="label-hint">{formatExpireTime(server.expireTime)}</span>
                    </label>
                    <input
                      type="date"
                      className="form-input"
                      value={new Date(server.expireTime * 1000).toISOString().split('T')[0]}
                      onChange={e => {
                        const date = new Date(e.target.value)
                        handleChange(index, 'expireTime', Math.floor(date.getTime() / 1000))
                      }}
                    />
                  </div>
                </div>

                <div className="form-group">
                  <label className="form-label">
                    <span className="label-text">{t.serverConfig.publicKey}</span>
                    <span className="label-hint">{t.serverConfig.publicKeyHint}</span>
                  </label>
                  <input
                    type="text"
                    className="form-input mono"
                    value={server.pubKeyBase64}
                    onChange={e => handleChange(index, 'pubKeyBase64', e.target.value)}
                    placeholder={t.serverConfig.publicKeyPlaceholder}
                  />
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      <div className="form-actions">
        <button 
          className="btn btn-secondary"
          onClick={handleReset}
          disabled={!hasChanges || loading}
        >
          {t.serverConfig.reset}
        </button>
        <button 
          className="btn btn-primary"
          onClick={handleSave}
          disabled={!hasChanges || loading}
        >
          {loading ? t.serverConfig.saving : t.serverConfig.save}
        </button>
      </div>

      {hasChanges && (
        <div className="unsaved-hint">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
          </svg>
          {t.serverConfig.unsavedChanges}
        </div>
      )}
    </div>
  )
}
