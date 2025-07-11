package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	RecentKeyCacheConf = &bigcache.Config{
		Shards:             1024,
		MaxEntriesInWindow: 5000,
		LifeWindow:         5 * time.Second,
		CleanWindow:        5 * time.Second,
		Verbose:            true,
	}
)

func NewRecentKeyCache() *bigcache.BigCache {
	cache, _ := bigcache.New(context.Background(), *RecentKeyCacheConf)
	return cache
}
