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
	dnsConfigWatch      io.Closer
	resourceConfigWatch io.Closer

	errLoadConfig = fmt.Errorf("dns config load error")
)

type Config struct {
	LogLevel int `json:"logLevel"`
}

type Resources struct {
	Resources []*Resource
}

type Resource struct {
	AuthServiceId  string `json:"authServiceId"`
	ResourceId     string `json:"resourceId"`
	ServerHostname string `json:"serverHostname"`
	ServerIp       string `json:"serverIp"`
	ServerPort     int    `json:"serverPort"`
}

func (p *ProxyService) loadDNSConfig() error {
	// config.toml
	fileName := filepath.Join(common.ExeDirPath, "etc", "config.toml")
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

func (p *ProxyService) loadResources() error {
	// resource.toml
	fileName := filepath.Join(common.ExeDirPath, "etc", "resource.toml")
	if err := p.updateResources(fileName); err != nil {
		// ignore error
		_ = err
	}

	resourceConfigWatch = utils.WatchFile(fileName, func() {
		log.Info("resource config: %s has been updated", fileName)
		p.updateResources(fileName)
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
		p.log.SetLogLevel(conf.LogLevel)
		return err
	}

	// update
	if p.config.LogLevel != conf.LogLevel {
		log.Info("set base log level to %d", conf.LogLevel)
		p.log.SetLogLevel(conf.LogLevel)
		p.config.LogLevel = conf.LogLevel
	}
	return err
}

func (p *ProxyService) updateResources(file string) (err error) {
	utils.CatchPanicThenRun(func() {
		err = errLoadConfig
	})

	content, err := os.ReadFile(file)
	if err != nil {
		log.Error("failed to read resource config: %v", err)
	}

	var resources Resources
	if err := toml.Unmarshal(content, &resources); err != nil {
		log.Error("failed to unmarshal resource config: %v", err)
	}
	tempResourceMap := make(map[string]*Resource)
	for _, resource := range resources.Resources {
		tempResourceMap[resource.ResourceId] = resource
	}

	p.resourceMapLock.Lock()
	defer p.resourceMapLock.Unlock()
	p.resourceMap = tempResourceMap

	return err
}

func (p *ProxyService) StopConfigWatch() {
	if dnsConfigWatch != nil {
		dnsConfigWatch.Close()
	}

	if resourceConfigWatch != nil {
		resourceConfigWatch.Close()
	}
}
