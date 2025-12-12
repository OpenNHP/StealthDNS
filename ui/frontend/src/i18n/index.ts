// 多语言支持

export type Language = 'zh-CN' | 'en-US' | 'ja-JP' | 'ko-KR' | 'de-DE' | 'fr-FR'

export interface LanguageOption {
  code: Language
  name: string
  nativeName: string
}

export const languages: LanguageOption[] = [
  { code: 'zh-CN', name: 'Chinese (Simplified)', nativeName: '简体中文' },
  { code: 'en-US', name: 'English', nativeName: 'English' },
  { code: 'ja-JP', name: 'Japanese', nativeName: '日本語' },
  { code: 'ko-KR', name: 'Korean', nativeName: '한국어' },
  { code: 'de-DE', name: 'German', nativeName: 'Deutsch' },
  { code: 'fr-FR', name: 'French', nativeName: 'Français' },
]

export interface Translations {
  // 通用
  app: {
    title: string
    version: string
  }
  // 标签页
  tabs: {
    status: string
    clientConfig: string
    serverConfig: string
    logs: string
    settings: string
  }
  // 状态面板
  status: {
    running: string
    stopped: string
    serviceRunning: string
    serviceStopped: string
    startService: string
    stopService: string
    restartService: string
    starting: string
    stopping: string
    autoRestart: string
    autoRestartDesc: string
    restartCount: string
    listenPort: string
    nhpDomain: string
    // 状态消息
    messageRunning: string
    messageStopped: string
    messageStarted: string
    messageStarting: string
    messageStopping: string
    // DNS 信息
    currentDNS: string
    loading: string
  }
  // 客户端配置
  clientConfig: {
    title: string
    description: string
    privateKey: string
    privateKeyHint: string
    privateKeyPlaceholder: string
    userId: string
    organizationId: string
    cipherScheme: string
    logLevel: string
    save: string
    reset: string
    saving: string
    unsavedChanges: string
  }
  // 服务器配置
  serverConfig: {
    title: string
    description: string
    addServer: string
    hostname: string
    hostnameHint: string
    ipAddress: string
    ipAddressHint: string
    port: string
    expireTime: string
    publicKey: string
    publicKeyHint: string
    publicKeyPlaceholder: string
    deleteServer: string
    newServer: string
    save: string
    reset: string
    saving: string
    unsavedChanges: string
  }
  // 设置
  settings: {
    title: string
    language: string
    languageDesc: string
    theme: string
    themeDesc: string
    darkTheme: string
    lightTheme: string
  }
  // 加密方案
  cipherSchemes: {
    gmsm: string
    curve25519: string
  }
  // 日志级别
  logLevels: {
    silent: string
    error: string
    info: string
    audit: string
    debug: string
    trace: string
  }
  // 通知消息
  messages: {
    serviceStarted: string
    serviceStopped: string
    serviceRestarted: string
    startFailed: string
    stopFailed: string
    restartFailed: string
    configSaved: string
    saveFailed: string
    serverConfigSaved: string
  }
  // 标题栏
  titleBar: {
    minimize: string
    maximize: string
    restore: string
    quit: string
  }
  // 日志
  logs: {
    title: string
    description: string
    noFiles: string
    noLogs: string
    loading: string
    refresh: string
    clear: string
    autoScroll: string
    autoRefresh: string
    searchPlaceholder: string
    allLevels: string
    entries: string
    filtered: string
    live: string
  }
}

// 简体中文
const zhCN: Translations = {
  app: {
    title: 'StealthDNS',
    version: '版本 1.0.0',
  },
  tabs: {
    status: '状态',
    clientConfig: '客户端配置',
    serverConfig: '服务器配置',
    logs: '日志',
    settings: '设置',
  },
  status: {
    running: '运行中',
    stopped: '已停止',
    serviceRunning: 'DNS 代理运行中',
    serviceStopped: 'DNS 代理已停止',
    startService: '启动服务',
    stopService: '停止服务',
    restartService: '重启服务',
    starting: '启动中...',
    stopping: '停止中...',
    autoRestart: '崩溃自动重启',
    autoRestartDesc: '当 DNS 服务异常退出时自动重启',
    restartCount: '自动重启次数',
    listenPort: '监听端口',
    nhpDomain: 'NHP 域名',
    messageRunning: 'DNS 代理服务正在运行',
    messageStopped: 'DNS 代理服务已停止',
    messageStarted: 'DNS 代理服务已成功启动',
    messageStarting: '正在启动 DNS 代理服务...',
    messageStopping: '正在停止 DNS 代理服务...',
    currentDNS: '当前 DNS 地址',
    loading: '加载中...',
  },
  clientConfig: {
    title: '客户端配置',
    description: '配置客户端私钥和连接参数',
    privateKey: '私钥 (Base64)',
    privateKeyHint: '客户端的私钥，用于与 NHP 服务器通信',
    privateKeyPlaceholder: '输入 Base64 编码的私钥',
    userId: '用户 ID',
    organizationId: '组织 ID',
    cipherScheme: '加密方案',
    logLevel: '日志级别',
    save: '保存配置',
    reset: '重置',
    saving: '保存中...',
    unsavedChanges: '有未保存的更改',
  },
  serverConfig: {
    title: 'NHP 服务器配置',
    description: '配置 NHP 服务器地址和公钥',
    addServer: '添加服务器',
    hostname: '主机名',
    hostnameHint: '域名形式，优先于 IP 地址',
    ipAddress: 'IP 地址',
    ipAddressHint: '当主机名为空时使用',
    port: '端口',
    expireTime: '过期时间',
    publicKey: '服务器公钥 (Base64)',
    publicKeyHint: 'NHP 服务器的公钥',
    publicKeyPlaceholder: '输入 Base64 编码的公钥',
    deleteServer: '删除服务器',
    newServer: '新服务器',
    save: '保存配置',
    reset: '重置',
    saving: '保存中...',
    unsavedChanges: '有未保存的更改',
  },
  settings: {
    title: '应用设置',
    language: '界面语言',
    languageDesc: '选择应用程序的显示语言',
    theme: '主题',
    themeDesc: '选择界面主题风格',
    darkTheme: '深色主题',
    lightTheme: '浅色主题',
  },
  cipherSchemes: {
    gmsm: 'GMSM (国密)',
    curve25519: 'Curve25519',
  },
  logLevels: {
    silent: '静默 (Silent)',
    error: '错误 (Error)',
    info: '信息 (Info)',
    audit: '审计 (Audit)',
    debug: '调试 (Debug)',
    trace: '追踪 (Trace)',
  },
  messages: {
    serviceStarted: 'DNS 服务已启动',
    serviceStopped: 'DNS 服务已停止',
    serviceRestarted: 'DNS 服务已重启',
    startFailed: '启动失败',
    stopFailed: '停止失败',
    restartFailed: '重启失败',
    configSaved: '客户端配置已保存',
    saveFailed: '保存失败',
    serverConfigSaved: '服务器配置已保存',
  },
  titleBar: {
    minimize: '最小化到托盘',
    maximize: '最大化',
    restore: '还原',
    quit: '退出',
  },
  logs: {
    title: '运行日志',
    description: '查看 DNS 代理服务的实时运行日志',
    noFiles: '没有日志文件',
    noLogs: '暂无日志',
    loading: '加载中...',
    refresh: '刷新',
    clear: '清空',
    autoScroll: '自动滚动',
    autoRefresh: '实时刷新',
    searchPlaceholder: '搜索日志...',
    allLevels: '所有级别',
    entries: '条目',
    filtered: '已过滤',
    live: '实时',
  },
}

// English
const enUS: Translations = {
  app: {
    title: 'StealthDNS',
    version: 'Version 1.0.0',
  },
  tabs: {
    status: 'Status',
    clientConfig: 'Client Config',
    serverConfig: 'Server Config',
    logs: 'Logs',
    settings: 'Settings',
  },
  status: {
    running: 'Running',
    stopped: 'Stopped',
    serviceRunning: 'DNS Proxy Running',
    serviceStopped: 'DNS Proxy Stopped',
    startService: 'Start Service',
    stopService: 'Stop Service',
    restartService: 'Restart Service',
    starting: 'Starting...',
    stopping: 'Stopping...',
    autoRestart: 'Auto Restart on Crash',
    autoRestartDesc: 'Automatically restart DNS service on abnormal exit',
    restartCount: 'Restart Count',
    listenPort: 'Listen Port',
    nhpDomain: 'NHP Domain',
    messageRunning: 'DNS proxy service is running',
    messageStopped: 'DNS proxy service is stopped',
    messageStarted: 'DNS proxy service started successfully',
    messageStarting: 'Starting DNS proxy service...',
    messageStopping: 'Stopping DNS proxy service...',
    currentDNS: 'Current DNS Address',
    loading: 'Loading...',
  },
  clientConfig: {
    title: 'Client Configuration',
    description: 'Configure client private key and connection parameters',
    privateKey: 'Private Key (Base64)',
    privateKeyHint: 'Client private key for NHP server communication',
    privateKeyPlaceholder: 'Enter Base64 encoded private key',
    userId: 'User ID',
    organizationId: 'Organization ID',
    cipherScheme: 'Cipher Scheme',
    logLevel: 'Log Level',
    save: 'Save Config',
    reset: 'Reset',
    saving: 'Saving...',
    unsavedChanges: 'Unsaved changes',
  },
  serverConfig: {
    title: 'NHP Server Configuration',
    description: 'Configure NHP server address and public key',
    addServer: 'Add Server',
    hostname: 'Hostname',
    hostnameHint: 'Domain name, takes priority over IP',
    ipAddress: 'IP Address',
    ipAddressHint: 'Used when hostname is empty',
    port: 'Port',
    expireTime: 'Expire Time',
    publicKey: 'Server Public Key (Base64)',
    publicKeyHint: 'NHP server public key',
    publicKeyPlaceholder: 'Enter Base64 encoded public key',
    deleteServer: 'Delete Server',
    newServer: 'New Server',
    save: 'Save Config',
    reset: 'Reset',
    saving: 'Saving...',
    unsavedChanges: 'Unsaved changes',
  },
  settings: {
    title: 'Application Settings',
    language: 'Interface Language',
    languageDesc: 'Select application display language',
    theme: 'Theme',
    themeDesc: 'Select interface theme style',
    darkTheme: 'Dark Theme',
    lightTheme: 'Light Theme',
  },
  cipherSchemes: {
    gmsm: 'GMSM (Chinese Standard)',
    curve25519: 'Curve25519',
  },
  logLevels: {
    silent: 'Silent',
    error: 'Error',
    info: 'Info',
    audit: 'Audit',
    debug: 'Debug',
    trace: 'Trace',
  },
  messages: {
    serviceStarted: 'DNS service started',
    serviceStopped: 'DNS service stopped',
    serviceRestarted: 'DNS service restarted',
    startFailed: 'Start failed',
    stopFailed: 'Stop failed',
    restartFailed: 'Restart failed',
    configSaved: 'Client configuration saved',
    saveFailed: 'Save failed',
    serverConfigSaved: 'Server configuration saved',
  },
  titleBar: {
    minimize: 'Minimize to Tray',
    maximize: 'Maximize',
    restore: 'Restore',
    quit: 'Quit',
  },
  logs: {
    title: 'Runtime Logs',
    description: 'View real-time logs of the DNS proxy service',
    noFiles: 'No log files',
    noLogs: 'No logs available',
    loading: 'Loading...',
    refresh: 'Refresh',
    clear: 'Clear',
    autoScroll: 'Auto Scroll',
    autoRefresh: 'Auto Refresh',
    searchPlaceholder: 'Search logs...',
    allLevels: 'All Levels',
    entries: 'Entries',
    filtered: 'Filtered',
    live: 'Live',
  },
}

// 日本語
const jaJP: Translations = {
  app: {
    title: 'StealthDNS',
    version: 'バージョン 1.0.0',
  },
  tabs: {
    status: 'ステータス',
    clientConfig: 'クライアント設定',
    serverConfig: 'サーバー設定',
    logs: 'ログ',
    settings: '設定',
  },
  status: {
    running: '実行中',
    stopped: '停止中',
    serviceRunning: 'DNS プロキシ実行中',
    serviceStopped: 'DNS プロキシ停止中',
    startService: 'サービス開始',
    stopService: 'サービス停止',
    restartService: 'サービス再起動',
    starting: '起動中...',
    stopping: '停止中...',
    autoRestart: 'クラッシュ時自動再起動',
    autoRestartDesc: 'DNS サービスが異常終了した場合に自動的に再起動',
    restartCount: '再起動回数',
    listenPort: 'リッスンポート',
    nhpDomain: 'NHP ドメイン',
    messageRunning: 'DNS プロキシサービスが実行中です',
    messageStopped: 'DNS プロキシサービスが停止しています',
    messageStarted: 'DNS プロキシサービスが正常に開始されました',
    messageStarting: 'DNS プロキシサービスを開始しています...',
    messageStopping: 'DNS プロキシサービスを停止しています...',
    currentDNS: '現在の DNS アドレス',
    loading: '読み込み中...',
  },
  clientConfig: {
    title: 'クライアント設定',
    description: 'クライアントの秘密鍵と接続パラメータを設定',
    privateKey: '秘密鍵 (Base64)',
    privateKeyHint: 'NHP サーバーとの通信に使用するクライアントの秘密鍵',
    privateKeyPlaceholder: 'Base64 エンコードされた秘密鍵を入力',
    userId: 'ユーザー ID',
    organizationId: '組織 ID',
    cipherScheme: '暗号方式',
    logLevel: 'ログレベル',
    save: '設定を保存',
    reset: 'リセット',
    saving: '保存中...',
    unsavedChanges: '未保存の変更があります',
  },
  serverConfig: {
    title: 'NHP サーバー設定',
    description: 'NHP サーバーのアドレスと公開鍵を設定',
    addServer: 'サーバーを追加',
    hostname: 'ホスト名',
    hostnameHint: 'ドメイン名形式、IP アドレスより優先',
    ipAddress: 'IP アドレス',
    ipAddressHint: 'ホスト名が空の場合に使用',
    port: 'ポート',
    expireTime: '有効期限',
    publicKey: 'サーバー公開鍵 (Base64)',
    publicKeyHint: 'NHP サーバーの公開鍵',
    publicKeyPlaceholder: 'Base64 エンコードされた公開鍵を入力',
    deleteServer: 'サーバーを削除',
    newServer: '新規サーバー',
    save: '設定を保存',
    reset: 'リセット',
    saving: '保存中...',
    unsavedChanges: '未保存の変更があります',
  },
  settings: {
    title: 'アプリケーション設定',
    language: 'インターフェース言語',
    languageDesc: 'アプリケーションの表示言語を選択',
    theme: 'テーマ',
    themeDesc: 'インターフェーステーマスタイルを選択',
    darkTheme: 'ダークテーマ',
    lightTheme: 'ライトテーマ',
  },
  cipherSchemes: {
    gmsm: 'GMSM (中国国家標準)',
    curve25519: 'Curve25519',
  },
  logLevels: {
    silent: 'サイレント',
    error: 'エラー',
    info: '情報',
    audit: '監査',
    debug: 'デバッグ',
    trace: 'トレース',
  },
  messages: {
    serviceStarted: 'DNS サービスが開始されました',
    serviceStopped: 'DNS サービスが停止されました',
    serviceRestarted: 'DNS サービスが再起動されました',
    startFailed: '開始に失敗しました',
    stopFailed: '停止に失敗しました',
    restartFailed: '再起動に失敗しました',
    configSaved: 'クライアント設定が保存されました',
    saveFailed: '保存に失敗しました',
    serverConfigSaved: 'サーバー設定が保存されました',
  },
  titleBar: {
    minimize: 'トレイに最小化',
    maximize: '最大化',
    restore: '元に戻す',
    quit: '終了',
  },
  logs: {
    title: '実行ログ',
    description: 'DNSプロキシサービスのリアルタイムログを表示',
    noFiles: 'ログファイルがありません',
    noLogs: 'ログがありません',
    loading: '読み込み中...',
    refresh: '更新',
    clear: 'クリア',
    autoScroll: '自動スクロール',
    autoRefresh: '自動更新',
    searchPlaceholder: 'ログを検索...',
    allLevels: 'すべてのレベル',
    entries: 'エントリ',
    filtered: 'フィルター済み',
    live: 'ライブ',
  },
}

// 한국어
const koKR: Translations = {
  app: {
    title: 'StealthDNS',
    version: '버전 1.0.0',
  },
  tabs: {
    status: '상태',
    clientConfig: '클라이언트 설정',
    serverConfig: '서버 설정',
    logs: '로그',
    settings: '설정',
  },
  status: {
    running: '실행 중',
    stopped: '중지됨',
    serviceRunning: 'DNS 프록시 실행 중',
    serviceStopped: 'DNS 프록시 중지됨',
    startService: '서비스 시작',
    stopService: '서비스 중지',
    restartService: '서비스 재시작',
    starting: '시작 중...',
    stopping: '중지 중...',
    autoRestart: '충돌 시 자동 재시작',
    autoRestartDesc: 'DNS 서비스가 비정상 종료 시 자동으로 재시작',
    restartCount: '재시작 횟수',
    listenPort: '수신 포트',
    nhpDomain: 'NHP 도메인',
    messageRunning: 'DNS 프록시 서비스가 실행 중입니다',
    messageStopped: 'DNS 프록시 서비스가 중지되었습니다',
    messageStarted: 'DNS 프록시 서비스가 성공적으로 시작되었습니다',
    messageStarting: 'DNS 프록시 서비스를 시작하는 중...',
    messageStopping: 'DNS 프록시 서비스를 중지하는 중...',
    currentDNS: '현재 DNS 주소',
    loading: '로딩 중...',
  },
  clientConfig: {
    title: '클라이언트 설정',
    description: '클라이언트 개인 키 및 연결 매개변수 설정',
    privateKey: '개인 키 (Base64)',
    privateKeyHint: 'NHP 서버 통신에 사용되는 클라이언트 개인 키',
    privateKeyPlaceholder: 'Base64로 인코딩된 개인 키 입력',
    userId: '사용자 ID',
    organizationId: '조직 ID',
    cipherScheme: '암호화 방식',
    logLevel: '로그 레벨',
    save: '설정 저장',
    reset: '재설정',
    saving: '저장 중...',
    unsavedChanges: '저장되지 않은 변경 사항이 있습니다',
  },
  serverConfig: {
    title: 'NHP 서버 설정',
    description: 'NHP 서버 주소 및 공개 키 설정',
    addServer: '서버 추가',
    hostname: '호스트명',
    hostnameHint: '도메인 형식, IP 주소보다 우선',
    ipAddress: 'IP 주소',
    ipAddressHint: '호스트명이 비어 있을 때 사용',
    port: '포트',
    expireTime: '만료 시간',
    publicKey: '서버 공개 키 (Base64)',
    publicKeyHint: 'NHP 서버의 공개 키',
    publicKeyPlaceholder: 'Base64로 인코딩된 공개 키 입력',
    deleteServer: '서버 삭제',
    newServer: '새 서버',
    save: '설정 저장',
    reset: '재설정',
    saving: '저장 중...',
    unsavedChanges: '저장되지 않은 변경 사항이 있습니다',
  },
  settings: {
    title: '애플리케이션 설정',
    language: '인터페이스 언어',
    languageDesc: '애플리케이션 표시 언어 선택',
    theme: '테마',
    themeDesc: '인터페이스 테마 스타일 선택',
    darkTheme: '다크 테마',
    lightTheme: '라이트 테마',
  },
  cipherSchemes: {
    gmsm: 'GMSM (중국 표준)',
    curve25519: 'Curve25519',
  },
  logLevels: {
    silent: '무음',
    error: '오류',
    info: '정보',
    audit: '감사',
    debug: '디버그',
    trace: '추적',
  },
  messages: {
    serviceStarted: 'DNS 서비스가 시작되었습니다',
    serviceStopped: 'DNS 서비스가 중지되었습니다',
    serviceRestarted: 'DNS 서비스가 재시작되었습니다',
    startFailed: '시작 실패',
    stopFailed: '중지 실패',
    restartFailed: '재시작 실패',
    configSaved: '클라이언트 설정이 저장되었습니다',
    saveFailed: '저장 실패',
    serverConfigSaved: '서버 설정이 저장되었습니다',
  },
  titleBar: {
    minimize: '트레이로 최소화',
    maximize: '최대화',
    restore: '복원',
    quit: '종료',
  },
  logs: {
    title: '실행 로그',
    description: 'DNS 프록시 서비스의 실시간 로그 보기',
    noFiles: '로그 파일이 없습니다',
    noLogs: '로그가 없습니다',
    loading: '로딩 중...',
    refresh: '새로 고침',
    clear: '지우기',
    autoScroll: '자동 스크롤',
    autoRefresh: '자동 새로 고침',
    searchPlaceholder: '로그 검색...',
    allLevels: '모든 레벨',
    entries: '항목',
    filtered: '필터됨',
    live: '실시간',
  },
}

// Deutsch
const deDE: Translations = {
  app: {
    title: 'StealthDNS',
    version: 'Version 1.0.0',
  },
  tabs: {
    status: 'Status',
    clientConfig: 'Client-Konfiguration',
    serverConfig: 'Server-Konfiguration',
    logs: 'Logs',
    settings: 'Einstellungen',
  },
  status: {
    running: 'Läuft',
    stopped: 'Gestoppt',
    serviceRunning: 'DNS-Proxy läuft',
    serviceStopped: 'DNS-Proxy gestoppt',
    startService: 'Dienst starten',
    stopService: 'Dienst stoppen',
    restartService: 'Dienst neustarten',
    starting: 'Startet...',
    stopping: 'Stoppt...',
    autoRestart: 'Automatischer Neustart bei Absturz',
    autoRestartDesc: 'DNS-Dienst bei abnormalem Beenden automatisch neustarten',
    restartCount: 'Neustarts',
    listenPort: 'Port',
    nhpDomain: 'NHP-Domain',
    messageRunning: 'DNS-Proxy-Dienst läuft',
    messageStopped: 'DNS-Proxy-Dienst ist gestoppt',
    messageStarted: 'DNS-Proxy-Dienst erfolgreich gestartet',
    messageStarting: 'DNS-Proxy-Dienst wird gestartet...',
    messageStopping: 'DNS-Proxy-Dienst wird gestoppt...',
    currentDNS: 'Aktuelle DNS-Adresse',
    loading: 'Laden...',
  },
  clientConfig: {
    title: 'Client-Konfiguration',
    description: 'Privaten Schlüssel und Verbindungsparameter konfigurieren',
    privateKey: 'Privater Schlüssel (Base64)',
    privateKeyHint: 'Client-Schlüssel für NHP-Server-Kommunikation',
    privateKeyPlaceholder: 'Base64-kodierten privaten Schlüssel eingeben',
    userId: 'Benutzer-ID',
    organizationId: 'Organisations-ID',
    cipherScheme: 'Verschlüsselungsschema',
    logLevel: 'Log-Level',
    save: 'Speichern',
    reset: 'Zurücksetzen',
    saving: 'Speichert...',
    unsavedChanges: 'Ungespeicherte Änderungen',
  },
  serverConfig: {
    title: 'NHP Server-Konfiguration',
    description: 'NHP-Serveradresse und öffentlichen Schlüssel konfigurieren',
    addServer: 'Server hinzufügen',
    hostname: 'Hostname',
    hostnameHint: 'Domainname, hat Vorrang vor IP',
    ipAddress: 'IP-Adresse',
    ipAddressHint: 'Verwendet wenn Hostname leer',
    port: 'Port',
    expireTime: 'Ablaufzeit',
    publicKey: 'Server-Öffentlicher Schlüssel (Base64)',
    publicKeyHint: 'Öffentlicher Schlüssel des NHP-Servers',
    publicKeyPlaceholder: 'Base64-kodierten öffentlichen Schlüssel eingeben',
    deleteServer: 'Server löschen',
    newServer: 'Neuer Server',
    save: 'Speichern',
    reset: 'Zurücksetzen',
    saving: 'Speichert...',
    unsavedChanges: 'Ungespeicherte Änderungen',
  },
  settings: {
    title: 'Anwendungseinstellungen',
    language: 'Sprache',
    languageDesc: 'Anzeigesprache der Anwendung auswählen',
    theme: 'Design',
    themeDesc: 'Interface-Design auswählen',
    darkTheme: 'Dunkles Design',
    lightTheme: 'Helles Design',
  },
  cipherSchemes: {
    gmsm: 'GMSM (Chinesischer Standard)',
    curve25519: 'Curve25519',
  },
  logLevels: {
    silent: 'Stumm',
    error: 'Fehler',
    info: 'Info',
    audit: 'Audit',
    debug: 'Debug',
    trace: 'Trace',
  },
  messages: {
    serviceStarted: 'DNS-Dienst gestartet',
    serviceStopped: 'DNS-Dienst gestoppt',
    serviceRestarted: 'DNS-Dienst neugestartet',
    startFailed: 'Start fehlgeschlagen',
    stopFailed: 'Stopp fehlgeschlagen',
    restartFailed: 'Neustart fehlgeschlagen',
    configSaved: 'Client-Konfiguration gespeichert',
    saveFailed: 'Speichern fehlgeschlagen',
    serverConfigSaved: 'Server-Konfiguration gespeichert',
  },
  titleBar: {
    minimize: 'In Taskleiste minimieren',
    maximize: 'Maximieren',
    restore: 'Wiederherstellen',
    quit: 'Beenden',
  },
  logs: {
    title: 'Laufzeitprotokolle',
    description: 'Echtzeit-Protokolle des DNS-Proxy-Dienstes anzeigen',
    noFiles: 'Keine Protokolldateien',
    noLogs: 'Keine Protokolle verfügbar',
    loading: 'Laden...',
    refresh: 'Aktualisieren',
    clear: 'Löschen',
    autoScroll: 'Automatisch scrollen',
    autoRefresh: 'Automatisch aktualisieren',
    searchPlaceholder: 'Protokolle durchsuchen...',
    allLevels: 'Alle Ebenen',
    entries: 'Einträge',
    filtered: 'Gefiltert',
    live: 'Live',
  },
}

// Français
const frFR: Translations = {
  app: {
    title: 'StealthDNS',
    version: 'Version 1.0.0',
  },
  tabs: {
    status: 'État',
    clientConfig: 'Config Client',
    serverConfig: 'Config Serveur',
    logs: 'Journaux',
    settings: 'Paramètres',
  },
  status: {
    running: 'En cours',
    stopped: 'Arrêté',
    serviceRunning: 'Proxy DNS en cours',
    serviceStopped: 'Proxy DNS arrêté',
    startService: 'Démarrer le service',
    stopService: 'Arrêter le service',
    restartService: 'Redémarrer le service',
    starting: 'Démarrage...',
    stopping: 'Arrêt...',
    autoRestart: 'Redémarrage auto en cas de crash',
    autoRestartDesc: 'Redémarrer automatiquement le service DNS en cas de sortie anormale',
    restartCount: 'Nombre de redémarrages',
    listenPort: 'Port d\'écoute',
    nhpDomain: 'Domaine NHP',
    messageRunning: 'Le service proxy DNS est en cours d\'exécution',
    messageStopped: 'Le service proxy DNS est arrêté',
    messageStarted: 'Le service proxy DNS a démarré avec succès',
    messageStarting: 'Démarrage du service proxy DNS...',
    messageStopping: 'Arrêt du service proxy DNS...',
    currentDNS: 'Adresse DNS actuelle',
    loading: 'Chargement...',
  },
  clientConfig: {
    title: 'Configuration Client',
    description: 'Configurer la clé privée et les paramètres de connexion',
    privateKey: 'Clé Privée (Base64)',
    privateKeyHint: 'Clé privée du client pour la communication avec le serveur NHP',
    privateKeyPlaceholder: 'Entrer la clé privée encodée en Base64',
    userId: 'ID Utilisateur',
    organizationId: 'ID Organisation',
    cipherScheme: 'Schéma de chiffrement',
    logLevel: 'Niveau de log',
    save: 'Enregistrer',
    reset: 'Réinitialiser',
    saving: 'Enregistrement...',
    unsavedChanges: 'Modifications non enregistrées',
  },
  serverConfig: {
    title: 'Configuration Serveur NHP',
    description: 'Configurer l\'adresse du serveur NHP et la clé publique',
    addServer: 'Ajouter un serveur',
    hostname: 'Nom d\'hôte',
    hostnameHint: 'Nom de domaine, prioritaire sur l\'IP',
    ipAddress: 'Adresse IP',
    ipAddressHint: 'Utilisé quand le nom d\'hôte est vide',
    port: 'Port',
    expireTime: 'Date d\'expiration',
    publicKey: 'Clé Publique du Serveur (Base64)',
    publicKeyHint: 'Clé publique du serveur NHP',
    publicKeyPlaceholder: 'Entrer la clé publique encodée en Base64',
    deleteServer: 'Supprimer le serveur',
    newServer: 'Nouveau serveur',
    save: 'Enregistrer',
    reset: 'Réinitialiser',
    saving: 'Enregistrement...',
    unsavedChanges: 'Modifications non enregistrées',
  },
  settings: {
    title: 'Paramètres de l\'application',
    language: 'Langue de l\'interface',
    languageDesc: 'Sélectionner la langue d\'affichage',
    theme: 'Thème',
    themeDesc: 'Sélectionner le style du thème',
    darkTheme: 'Thème sombre',
    lightTheme: 'Thème clair',
  },
  cipherSchemes: {
    gmsm: 'GMSM (Standard chinois)',
    curve25519: 'Curve25519',
  },
  logLevels: {
    silent: 'Silencieux',
    error: 'Erreur',
    info: 'Info',
    audit: 'Audit',
    debug: 'Debug',
    trace: 'Trace',
  },
  messages: {
    serviceStarted: 'Service DNS démarré',
    serviceStopped: 'Service DNS arrêté',
    serviceRestarted: 'Service DNS redémarré',
    startFailed: 'Échec du démarrage',
    stopFailed: 'Échec de l\'arrêt',
    restartFailed: 'Échec du redémarrage',
    configSaved: 'Configuration client enregistrée',
    saveFailed: 'Échec de l\'enregistrement',
    serverConfigSaved: 'Configuration serveur enregistrée',
  },
  titleBar: {
    minimize: 'Réduire dans la barre des tâches',
    maximize: 'Agrandir',
    restore: 'Restaurer',
    quit: 'Quitter',
  },
  logs: {
    title: 'Journaux d\'exécution',
    description: 'Afficher les journaux en temps réel du service proxy DNS',
    noFiles: 'Pas de fichiers journaux',
    noLogs: 'Aucun journal disponible',
    loading: 'Chargement...',
    refresh: 'Actualiser',
    clear: 'Effacer',
    autoScroll: 'Défilement auto',
    autoRefresh: 'Actualisation auto',
    searchPlaceholder: 'Rechercher dans les journaux...',
    allLevels: 'Tous les niveaux',
    entries: 'Entrées',
    filtered: 'Filtré',
    live: 'En direct',
  },
}

// 翻译映射
const translations: Record<Language, Translations> = {
  'zh-CN': zhCN,
  'en-US': enUS,
  'ja-JP': jaJP,
  'ko-KR': koKR,
  'de-DE': deDE,
  'fr-FR': frFR,
}

// 获取翻译
export function getTranslations(lang: Language): Translations {
  return translations[lang] || translations['en-US']
}

// 获取浏览器默认语言
export function getBrowserLanguage(): Language {
  const browserLang = navigator.language
  
  // 精确匹配
  if (browserLang in translations) {
    return browserLang as Language
  }
  
  // 语言代码匹配 (例如 zh -> zh-CN)
  const langCode = browserLang.split('-')[0]
  const match = Object.keys(translations).find(key => key.startsWith(langCode))
  if (match) {
    return match as Language
  }
  
  return 'en-US'
}

// 从本地存储获取语言设置
export function getStoredLanguage(): Language {
  const stored = localStorage.getItem('stealthdns-language')
  if (stored && stored in translations) {
    return stored as Language
  }
  return getBrowserLanguage()
}

// 保存语言设置到本地存储
export function setStoredLanguage(lang: Language): void {
  localStorage.setItem('stealthdns-language', lang)
}


