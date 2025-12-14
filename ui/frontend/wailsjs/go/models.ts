export namespace main {
	
	export class ClientConfig {
	    privateKeyBase64: string;
	    defaultCipherScheme: number;
	    userId: string;
	    organizationId: string;
	    logLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new ClientConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.privateKeyBase64 = source["privateKeyBase64"];
	        this.defaultCipherScheme = source["defaultCipherScheme"];
	        this.userId = source["userId"];
	        this.organizationId = source["organizationId"];
	        this.logLevel = source["logLevel"];
	    }
	}
	export class LogEntry {
	    timestamp: string;
	    level: string;
	    message: string;
	    raw: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.level = source["level"];
	        this.message = source["message"];
	        this.raw = source["raw"];
	    }
	}
	export class LogFile {
	    name: string;
	    size: number;
	    modified: string;
	
	    static createFrom(source: any = {}) {
	        return new LogFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.modified = source["modified"];
	    }
	}
	export class ServerConfig {
	    hostname: string;
	    ip: string;
	    port: number;
	    pubKeyBase64: string;
	    expireTime: number;
	
	    static createFrom(source: any = {}) {
	        return new ServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.pubKeyBase64 = source["pubKeyBase64"];
	        this.expireTime = source["expireTime"];
	    }
	}
	export class ServiceStatus {
	    running: boolean;
	    restartCount: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.restartCount = source["restartCount"];
	        this.message = source["message"];
	    }
	}
	export class SystemDNSInfo {
	    dnsServers: string[];
	    listenPort: number;
	    isProxyActive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SystemDNSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dnsServers = source["dnsServers"];
	        this.listenPort = source["listenPort"];
	        this.isProxyActive = source["isProxyActive"];
	    }
	}
	export class TrayManager {
	
	
	    static createFrom(source: any = {}) {
	        return new TrayManager(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class VersionInfo {
	    version: string;
	    buildNumber: string;
	    commitId: string;
	    buildTime: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.buildNumber = source["buildNumber"];
	        this.commitId = source["commitId"];
	        this.buildTime = source["buildTime"];
	    }
	}

}

