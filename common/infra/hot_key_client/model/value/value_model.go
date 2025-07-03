package value

type ValueModel struct {
	CreatedAt int64 `json:"created_at"` //创建时间
	Duration  int64 `json:"duration"`   //本地缓存时间 单位毫秒
	Value     any   `json:"value"`      //数据值
}

func (model *ValueModel) GetDefaultValue() any {
	if model.Value == nil {
		return nil
	}
	return model.Value
}

// NewValueModel 创建ValueModel,duration本地缓存时间,单位毫秒,value为数据值
func NewValueModel(duration int64, value any) *ValueModel {
	return &ValueModel{
		Duration: duration,
		Value:    value,
	}
}
