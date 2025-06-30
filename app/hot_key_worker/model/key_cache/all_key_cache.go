package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	AllKeyCacheConf = &bigcache.Config{
		Shards:             1024,
		LifeWindow:         1 * time.Minute,
		CleanWindow:        1 * time.Minute,
		Verbose:            true,
		MaxEntriesInWindow: 50000,
	}
)

func NewAllKeyCache() *bigcache.BigCache {
	cache, _ := bigcache.New(context.Background(), *AllKeyCacheConf)
	return cache
}
