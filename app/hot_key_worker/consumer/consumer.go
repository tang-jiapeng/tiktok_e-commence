package consumer

import (
	"context"
	"hot_key/model/key"
	cache "hot_key/model/key_cache"
	"hot_key/model/util"
	"hot_key/redis"
	"hot_key/tool"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
)

var hotKeyCache *bigcache.BigCache

func init() {
	hotKeyCache = cache.NewRecentKeyCache()
}

func Consume() {
	for {
		hotKeyModel := util.BlQueue.Take()
		if hotKeyModel.Remove {
			err := removeKey(hotKeyModel)
			if err != nil {
				klog.CtxErrorf(context.Background(), "remove key error: %s", err)
			}
		} else {
			err := newKey(hotKeyModel)
			if err != nil {
				klog.CtxErrorf(context.Background(), "new key error: %s", err)
			}
		}
	}
}

func removeKey(hotKeyModel key.HotKeyModel) (err error) {
	buildKey := BuildKey(hotKeyModel)
	//从热键缓存中删除
	err = hotKeyCache.Delete(buildKey)
	if err != nil {
		return err
	}
	//通知所有的client集群删除
	err = redis.PublishClientChannel(hotKeyModel)
	if err != nil {
		return err
	}
	return nil
}

func newKey(hotKeyModel key.HotKeyModel) (err error) {
	buildKey := BuildKey(hotKeyModel)
	// 判断key是否刚热过
	_, err = hotKeyCache.Get(buildKey)
	if err == nil {
		return nil
	}
	slidingWindow := getWindow(hotKeyModel, buildKey)
	hot := slidingWindow.AddCount(hotKeyModel.Count.GetCount())
	//不热放进所有键的缓存
	if !hot {
		marshal, _ := sonic.Marshal(slidingWindow)
		cache.GetAllKeyCache(hotKeyModel.ServiceName).Set(buildKey, marshal)
		return nil
	}
	//热放进热键的缓存
	hotKeyCache.Set(buildKey, []byte("hot"))
	hotKeyModel.CreatedAt = time.Now().UnixMilli()
	err = redis.PublishClientChannel(hotKeyModel)
	if err != nil {
		return err
	}
	return nil
}

func BuildKey(hotKeyModel key.HotKeyModel) string {
	return hotKeyModel.ServiceName + "+" + hotKeyModel.Key
}

func getWindow(hotKeyModel key.HotKeyModel, key string) *tool.SlidingWindow {
	bigCache := cache.GetAllKeyCache(hotKeyModel.ServiceName)
	window, emptyError := bigCache.Get(key)
	if emptyError != nil {
		slideWindow := tool.NewSlidingWindow(hotKeyModel.Interval, hotKeyModel.Threshold)
		windowJson, _ := sonic.Marshal(slideWindow)
		bigCache.Set(key, windowJson)
		return slideWindow
	}
	var slideWindow tool.SlidingWindow
	sonic.Unmarshal(window, &slideWindow)
	return &slideWindow
}
