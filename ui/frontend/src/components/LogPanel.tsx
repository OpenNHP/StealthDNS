import { useState, useEffect, useRef, useCallback } from 'react'
import { useLanguage } from '../contexts/LanguageContext'
import '../styles/LogPanel.css'

// 日志类型定义（使用全局类型）
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

export function LogPanel() {
  const { t } = useLanguage()
  const [logFiles, setLogFiles] = useState<LogFile[]>([])
  const [selectedFile, setSelectedFile] = useState<string>('')
  const [logEntries, setLogEntries] = useState<LogEntry[]>([])
  const [loading, setLoading] = useState(false)
  const [autoScroll, setAutoScroll] = useState(true)
  const [autoRefresh, setAutoRefresh] = useState(true)
  const [filter, setFilter] = useState('')
  const [levelFilter, setLevelFilter] = useState<string>('ALL')
  const logContainerRef = useRef<HTMLDivElement>(null)
  const lastSizeRef = useRef<number>(0)

  // 加载日志文件列表
  const loadLogFiles = useCallback(async () => {
    try {
      const files = await window.go.main.App.GetLogFiles()
      setLogFiles(files || [])
      if (files && files.length > 0 && !selectedFile) {
        setSelectedFile(files[0].name)
      }
    } catch (err) {
      console.error('Failed to load log files:', err)
    }
  }, [selectedFile])

  // 加载日志内容
  const loadLogContent = useCallback(async () => {
    if (!selectedFile) return
    
    setLoading(true)
    try {
      const entries = await window.go.main.App.GetLogContent(selectedFile, 500)
      setLogEntries(entries || [])
      lastSizeRef.current = 0 // 重置，下次从头开始监听
      
      // 滚动到底部
      if (autoScroll && logContainerRef.current) {
        setTimeout(() => {
          logContainerRef.current?.scrollTo({
            top: logContainerRef.current.scrollHeight,
            behavior: 'smooth'
          })
        }, 100)
      }
    } catch (err) {
      console.error('Failed to load log content:', err)
    } finally {
      setLoading(false)
    }
  }, [selectedFile, autoScroll])

  // 刷新新日志
  const refreshLogs = useCallback(async () => {
    if (!selectedFile || !autoRefresh) return
    
    try {
      const result = await window.go.main.App.WatchLogFile(selectedFile, lastSizeRef.current)
      if (result && result.entries && result.entries.length > 0) {
        setLogEntries(prev => [...prev, ...result.entries])
        lastSizeRef.current = result.size
        
        if (autoScroll && logContainerRef.current) {
          setTimeout(() => {
            logContainerRef.current?.scrollTo({
              top: logContainerRef.current.scrollHeight,
              behavior: 'smooth'
            })
          }, 100)
        }
      }
    } catch (err) {
      console.error('Failed to refresh logs:', err)
    }
  }, [selectedFile, autoRefresh, autoScroll])

  // 清空日志
  const handleClearLog = async () => {
    if (!selectedFile) return
    
    try {
      await window.go.main.App.ClearLogFile(selectedFile)
      setLogEntries([])
      lastSizeRef.current = 0
    } catch (err) {
      console.error('Failed to clear log:', err)
    }
  }

  // 初始加载
  useEffect(() => {
    loadLogFiles()
  }, [loadLogFiles])

  // 当选择的文件变化时加载内容
  useEffect(() => {
    if (selectedFile) {
      loadLogContent()
    }
  }, [selectedFile, loadLogContent])

  // 自动刷新
  useEffect(() => {
    if (!autoRefresh) return
    
    const interval = setInterval(refreshLogs, 2000)
    return () => clearInterval(interval)
  }, [autoRefresh, refreshLogs])

  // 过滤日志
  const filteredEntries = logEntries.filter(entry => {
    // 级别过滤
    if (levelFilter !== 'ALL' && entry.level !== levelFilter) {
      return false
    }
    // 文本过滤
    if (filter && !entry.raw.toLowerCase().includes(filter.toLowerCase())) {
      return false
    }
    return true
  })

  // 获取日志级别样式
  const getLevelClass = (level: string) => {
    switch (level?.toUpperCase()) {
      case 'ERROR':
      case 'ERR':
        return 'log-error'
      case 'WARN':
      case 'WARNING':
        return 'log-warn'
      case 'DEBUG':
        return 'log-debug'
      case 'TRACE':
        return 'log-trace'
      case 'INFO':
      default:
        return 'log-info'
    }
  }

  // 格式化文件大小
  const formatSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  return (
    <div className="log-panel">
      <div className="panel-header">
        <div className="header-left">
          <h2>{t.logs.title}</h2>
          <p className="panel-desc">{t.logs.description}</p>
        </div>
      </div>

      {/* 工具栏 */}
      <div className="log-toolbar">
        <div className="toolbar-left">
          {/* 文件选择 */}
          <select 
            className="form-select file-select"
            value={selectedFile}
            onChange={(e) => setSelectedFile(e.target.value)}
          >
            {logFiles.length === 0 ? (
              <option value="">{t.logs.noFiles}</option>
            ) : (
              logFiles.map(file => (
                <option key={file.name} value={file.name}>
                  {file.name} ({formatSize(file.size)})
                </option>
              ))
            )}
          </select>

          {/* 级别过滤 */}
          <select 
            className="form-select level-select"
            value={levelFilter}
            onChange={(e) => setLevelFilter(e.target.value)}
          >
            <option value="ALL">{t.logs.allLevels}</option>
            <option value="ERROR">ERROR</option>
            <option value="WARN">WARN</option>
            <option value="INFO">INFO</option>
            <option value="DEBUG">DEBUG</option>
            <option value="TRACE">TRACE</option>
          </select>

          {/* 搜索过滤 */}
          <div className="search-input">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
              <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
            </svg>
            <input
              type="text"
              className="form-input"
              placeholder={t.logs.searchPlaceholder}
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>
        </div>

        <div className="toolbar-right">
          {/* 自动滚动 */}
          <label className="toolbar-toggle">
            <input
              type="checkbox"
              checked={autoScroll}
              onChange={(e) => setAutoScroll(e.target.checked)}
            />
            <span>{t.logs.autoScroll}</span>
          </label>

          {/* 自动刷新 */}
          <label className="toolbar-toggle">
            <input
              type="checkbox"
              checked={autoRefresh}
              onChange={(e) => setAutoRefresh(e.target.checked)}
            />
            <span>{t.logs.autoRefresh}</span>
          </label>

          {/* 刷新按钮 */}
          <button 
            className="btn btn-icon"
            onClick={loadLogContent}
            disabled={loading}
            title={t.logs.refresh}
          >
            <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
              <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
            </svg>
          </button>

          {/* 清空按钮 */}
          <button 
            className="btn btn-icon btn-danger-icon"
            onClick={handleClearLog}
            title={t.logs.clear}
          >
            <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
              <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
            </svg>
          </button>
        </div>
      </div>

      {/* 日志内容 */}
      <div className="log-container" ref={logContainerRef}>
        {loading && logEntries.length === 0 ? (
          <div className="log-loading">
            <div className="loading-spinner"></div>
            <span>{t.logs.loading}</span>
          </div>
        ) : filteredEntries.length === 0 ? (
          <div className="log-empty">
            <svg viewBox="0 0 24 24" width="48" height="48" fill="currentColor">
              <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/>
            </svg>
            <span>{t.logs.noLogs}</span>
          </div>
        ) : (
          <div className="log-entries">
            {filteredEntries.map((entry, index) => (
              <div key={index} className={`log-entry ${getLevelClass(entry.level)}`}>
                {entry.timestamp && (
                  <span className="log-timestamp">{entry.timestamp}</span>
                )}
                {entry.level && (
                  <span className={`log-level ${getLevelClass(entry.level)}`}>
                    {entry.level}
                  </span>
                )}
                <span className="log-message">{entry.message || entry.raw}</span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 状态栏 */}
      <div className="log-statusbar">
        <span>{t.logs.entries}: {filteredEntries.length}</span>
        {filter && <span>{t.logs.filtered}: {logEntries.length - filteredEntries.length}</span>}
        {autoRefresh && <span className="status-live">● {t.logs.live}</span>}
      </div>
    </div>
  )
}

