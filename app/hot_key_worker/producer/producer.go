package producer

import (
	"context"
	"hot_key/constant"
	"hot_key/model/key"
	"hot_key/model/util"
	"hot_key/redis"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Listener() {
	register()
	pubsub := redis.Rdb.Subscribe(context.Background(), constant.WorkerChannelId.String())
	ch := pubsub.Channel()

	for msg := range ch {
		klog.CtxInfof(context.Background(), "receive message from channel %s, payload: %s", msg.Channel, msg.Payload)
		var hotKeyModel key.HotKeyModel
		sonic.Unmarshal([]byte(msg.Payload), &hotKeyModel)
		klog.CtxInfof(context.Background(), "unmarshalled hotkey model: %v", hotKeyModel)
		produce(hotKeyModel)
	}
	defer pubsub.Close()
}

// register 将生成的worker channel id 注册到redis中
func register() {
	redis.Rdb.LPush(context.Background(), constant.WorkerChannelIdList, constant.WorkerChannelId.String())
	klog.CtxInfof(context.Background(), "Worker %s is listening...", constant.WorkerChannelId.String())
}

// Checkout 将worker channel id 从redis中移除
func Checkout() {
	redis.Rdb.LRem(context.Background(), constant.WorkerChannelIdList, 0, constant.WorkerChannelId.String())
	klog.CtxInfof(context.Background(), "Worker %s is offline", constant.WorkerChannelId.String())
}

func produce(hotKeyModel key.HotKeyModel) {
	util.BlQueue.Put(hotKeyModel)
}
