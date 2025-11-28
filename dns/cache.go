package dns

import (
	"sync"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/sync/singleflight"
)

type CacheItem struct {
	value      *dns.Msg
	expireTime time.Time
}

type StealthDNSCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheItem
	group singleflight.Group
	ttl   time.Duration
}

func NewStealthDNSCache(ttl time.Duration) *StealthDNSCache {
	return &StealthDNSCache{
		cache: make(map[string]*CacheItem),
		ttl:   ttl,
	}
}

func (pc *StealthDNSCache) GetCache(key string) (*CacheItem, bool) {
	item, found := pc.get(key)
	if !found {
		return nil, false
	}
	return item, true
}

func (pc *StealthDNSCache) SetCache(key string, value *dns.Msg) {
	pc.set(key, value, pc.ttl)
}

func (pc *StealthDNSCache) SetCacheWithTTL(key string, value *dns.Msg, ttl time.Duration) {
	pc.set(key, value, ttl)
}

func (pc *StealthDNSCache) Delete(key string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.cache, key)
}

func (pc *StealthDNSCache) get(key string) (*CacheItem, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	item, exists := pc.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expireTime) {
		return nil, false
	}

	return item, true
}

func (pc *StealthDNSCache) set(key string, value *dns.Msg, ttl time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.cache[key] = &CacheItem{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
}

func (pc *StealthDNSCache) CleanupExpired() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for key, item := range pc.cache {
		if now.After(item.expireTime) {
			delete(pc.cache, key)
		}
	}
}
