package key

import (
	"hot_key/model/base"
	"strconv"
)

type HotKeyModel struct {
	ID          string           `json:"id"`
	CreatedAt   int64            `json:"created_at"` //创建时间时间戳，单位毫秒
	Key         string           `json:"key"`
	ServiceName string           `json:"service_name"`
	Count       base.AtomicCount `json:"atomic_count"`
	Remove      bool             `json:"remove"` //是否删除
	KeyRule     `json:"key_rule"`
}

func (h *HotKeyModel) String() string {
	return "HotKeyModel{" +
		"ID: " + h.ID + ", " +
		"CreateAt: " + strconv.Itoa(int(h.CreatedAt)) + ", " +
		"Key: " + h.Key + ", " +
		"ServiceName: " + h.ServiceName + ", " +
		"Count: " + h.Count.String() + ", " +
		"KeyRule: " + h.KeyRule.String() + ", " +
		"Remove: " + strconv.FormatBool(h.Remove) +
		"}"
}
