package key

import "fmt"

type KeyRule struct {
	Interval  int64 `json:"interval"`  //触发间隔,单位为秒
	Threshold int64 `json:"threshold"` //触发阈值
	Duration  int64 `json:"duration"`  //在client端的缓存时间，单位为秒
}

func (k *KeyRule) String() string {
	return fmt.Sprintf("interval: %d, threshold: %d, duration: %d", k.Interval, k.Threshold, k.Duration)
}
