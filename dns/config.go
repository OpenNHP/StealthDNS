package dns

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/OpenNHP/opennhp/nhp/log"
	"github.com/OpenNHP/opennhp/nhp/utils"
	"github.com/pelletier/go-toml/v2"

	"github.com/OpenNHP/StealthDNS/common"
)

var (
	dnsConfigWatch io.Closer

	errLoadConfig = fmt.Errorf("dns config load error")
)

type Config struct {
	UpstreamDNS    string `json:"upstreamDNS"`
	SetSystemDNS   bool   `json:"setSystemDNS"`
	LogLevel       int    `json:"logLevel"`
	RemoveLocalDNS bool   `json:"removeLocalDNS"`
}

func (p *ProxyService) loadDNSConfig() error {
	// config.toml
	fileName := filepath.Join(common.ExeDirPath, "etc", "dns.toml")
	if err := p.updateDNSConfig(fileName); err != nil {
		// report base config error
		return err
	}

	dnsConfigWatch = utils.WatchFile(fileName, func() {
		log.Info("base config: %s has been updated", fileName)
		p.updateDNSConfig(fileName)
	})
	return nil
}

func (p *ProxyService) updateDNSConfig(file string) (err error) {
	utils.CatchPanicThenRun(func() {
		err = errLoadConfig
	})

	content, err := os.ReadFile(file)
	if err != nil {
		log.Error("failed to read dns config: %v", err)
	}

	var conf Config
	if err := toml.Unmarshal(content, &conf); err != nil {
		log.Error("failed to unmarshal dns config: %v", err)
	}

	if p.config == nil {
		p.config = &conf
		p.upstreamDNS = conf.UpstreamDNS
		p.forward = len(p.upstreamDNS) > 0
		p.log.SetLogLevel(conf.LogLevel)
		return err
	}

	// update
	if p.config.LogLevel != conf.LogLevel {
		log.Info("set base log level to %d", conf.LogLevel)
		p.log.SetLogLevel(conf.LogLevel)
		p.config.LogLevel = conf.LogLevel
		p.upstreamDNS = conf.UpstreamDNS
		p.forward = len(p.upstreamDNS) > 0
	}
	return err
}

func (p *ProxyService) StopConfigWatch() {
	if dnsConfigWatch != nil {
		dnsConfigWatch.Close()
	}
}
