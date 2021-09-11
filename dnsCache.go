package main

import (
	"net"
	"sync"
	"time"
)

type item struct {
	value      net.IP
	lastAccess int64
}

type DnsCache struct {
	mutex sync.Mutex
	cache map[string]*item
}

func NewDnsCache(maxTTL int) *DnsCache {
	c := &DnsCache{}
	c.cache = make(map[string]*item)
	go func() {
		for now := range time.Tick(time.Second) {
			c.mutex.Lock()
			for k, v := range c.cache {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(c.cache, k)
				}
			}
			c.mutex.Unlock()
		}
	}()
	return c
}
