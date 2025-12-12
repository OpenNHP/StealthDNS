import { useState, useEffect } from 'react'
import { ClientConfig } from '../App'
import { useLanguage } from '../contexts/LanguageContext'
import '../styles/ConfigPanel.css'

interface ConfigPanelProps {
  config: ClientConfig
  loading: boolean
  onSave: (config: ClientConfig) => void
}

export function ConfigPanel({ config, loading, onSave }: ConfigPanelProps) {
  const { t } = useLanguage()
  const [formData, setFormData] = useState<ClientConfig>(config)
  const [hasChanges, setHasChanges] = useState(false)
  const [showPassword, setShowPassword] = useState(false)

  useEffect(() => {
    setFormData(config)
    setHasChanges(false)
  }, [config])

  const handleChange = (field: keyof ClientConfig, value: string | number) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    setHasChanges(true)
  }

  const handleSave = () => {
    onSave(formData)
    setHasChanges(false)
  }

  const handleReset = () => {
    setFormData(config)
    setHasChanges(false)
  }

  const cipherSchemeLabels: Record<number, string> = {
    0: t.cipherSchemes.gmsm,
    1: t.cipherSchemes.curve25519
  }

  const logLevelLabels: Record<number, string> = {
    0: t.logLevels.silent,
    1: t.logLevels.error,
    2: t.logLevels.info,
    3: t.logLevels.audit,
    4: t.logLevels.debug,
    5: t.logLevels.trace
  }

  return (
    <div className="config-panel">
      <div className="panel-header">
        <h2>{t.clientConfig.title}</h2>
        <p className="panel-desc">{t.clientConfig.description}</p>
      </div>

      <div className="form-group">
        <label className="form-label">
          <span className="label-text">{t.clientConfig.privateKey}</span>
          <span className="label-hint">{t.clientConfig.privateKeyHint}</span>
        </label>
        <div className="input-wrapper">
          <input
            type={showPassword ? "text" : "password"}
            className="form-input mono"
            value={formData.privateKeyBase64}
            onChange={e => handleChange('privateKeyBase64', e.target.value)}
            placeholder={t.clientConfig.privateKeyPlaceholder}
          />
          <button 
            className="input-action"
            onClick={() => setShowPassword(!showPassword)}
            title={showPassword ? "Hide" : "Show"}
          >
            <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
              <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
            </svg>
          </button>
        </div>
      </div>

      <div className="form-row">
        <div className="form-group">
          <label className="form-label">
            <span className="label-text">{t.clientConfig.userId}</span>
          </label>
          <input
            type="text"
            className="form-input"
            value={formData.userId}
            onChange={e => handleChange('userId', e.target.value)}
            placeholder="agent-0"
          />
        </div>
        <div className="form-group">
          <label className="form-label">
            <span className="label-text">{t.clientConfig.organizationId}</span>
          </label>
          <input
            type="text"
            className="form-input"
            value={formData.organizationId}
            onChange={e => handleChange('organizationId', e.target.value)}
            placeholder="opennhp.cn"
          />
        </div>
      </div>

      <div className="form-row">
        <div className="form-group">
          <label className="form-label">
            <span className="label-text">{t.clientConfig.cipherScheme}</span>
          </label>
          <select
            className="form-select"
            value={formData.defaultCipherScheme}
            onChange={e => handleChange('defaultCipherScheme', parseInt(e.target.value))}
          >
            {Object.entries(cipherSchemeLabels).map(([value, label]) => (
              <option key={value} value={value}>{label}</option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label className="form-label">
            <span className="label-text">{t.clientConfig.logLevel}</span>
          </label>
          <select
            className="form-select"
            value={formData.logLevel}
            onChange={e => handleChange('logLevel', parseInt(e.target.value))}
          >
            {Object.entries(logLevelLabels).map(([value, label]) => (
              <option key={value} value={value}>{label}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="form-actions">
        <button 
          className="btn btn-secondary"
          onClick={handleReset}
          disabled={!hasChanges || loading}
        >
          {t.clientConfig.reset}
        </button>
        <button 
          className="btn btn-primary"
          onClick={handleSave}
          disabled={!hasChanges || loading}
        >
          {loading ? t.clientConfig.saving : t.clientConfig.save}
        </button>
      </div>

      {hasChanges && (
        <div className="unsaved-hint">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
          </svg>
          {t.clientConfig.unsavedChanges}
        </div>
      )}
    </div>
  )
}
