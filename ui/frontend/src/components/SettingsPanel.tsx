import { useLanguage } from '../contexts/LanguageContext'
import { languages, Language } from '../i18n'
import '../styles/SettingsPanel.css'

export function SettingsPanel() {
  const { language, setLanguage, t } = useLanguage()

  const handleLanguageChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setLanguage(e.target.value as Language)
  }

  return (
    <div className="settings-panel">
      <div className="panel-header">
        <h2>{t.settings.title}</h2>
      </div>

      <div className="settings-section">
        <div className="setting-card">
          <div className="setting-icon">
            <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor">
              <path d="M12.87 15.07l-2.54-2.51.03-.03c1.74-1.94 2.98-4.17 3.71-6.53H17V4h-7V2H8v2H1v1.99h11.17C11.5 7.92 10.44 9.75 9 11.35 8.07 10.32 7.3 9.19 6.69 8h-2c.73 1.63 1.73 3.17 2.98 4.56l-5.09 5.02L4 19l5-5 3.11 3.11.76-2.04zM18.5 10h-2L12 22h2l1.12-3h4.75L21 22h2l-4.5-12zm-2.62 7l1.62-4.33L19.12 17h-3.24z"/>
            </svg>
          </div>
          <div className="setting-content">
            <div className="setting-info">
              <span className="setting-label">{t.settings.language}</span>
              <span className="setting-desc">{t.settings.languageDesc}</span>
            </div>
            <select 
              className="form-select language-select"
              value={language}
              onChange={handleLanguageChange}
            >
              {languages.map(lang => (
                <option key={lang.code} value={lang.code}>
                  {lang.nativeName}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      <div className="language-preview">
        <h3>语言预览 / Language Preview</h3>
        <div className="preview-grid">
          {languages.map(lang => (
            <div 
              key={lang.code} 
              className={`preview-card ${language === lang.code ? 'active' : ''}`}
              onClick={() => setLanguage(lang.code)}
            >
              <span className="preview-native">{lang.nativeName}</span>
              <span className="preview-english">{lang.name}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}


