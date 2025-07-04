package key

type KeyConf struct {
	ServiceName string `json:"service_name"` // 微服务名称，用于区分不同微服务
	Threshold   int64  `json:"threshold"`    //要判断为热key的key的阈值
	Interval    int64  `json:"interval"`     //判断为热key的key的判断时间,单位秒
	Duration    int64  `json:"duration"`     //是hotkey的话在本地的缓存时间，单位毫秒,为0的话为永久不过期
}

// NewKeyConf1 返回一个KeyConf，在缓存中存储3000毫秒，每10秒检测一次，超过100次判定为热键
func NewKeyConf1(serviceName string) KeyConf {
	return KeyConf{
		ServiceName: serviceName,
		Threshold:   1,
		Interval:    2,
		Duration:    3000,
	}
}
