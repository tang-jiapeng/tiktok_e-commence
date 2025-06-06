package error

import (
	"fmt"
	"tiktok_e-commerce/common/constant"
)

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("错误 %d: %s", e.Code, e.Message)
}

func NewServiceError(code int) error {
	return &ServiceError{
		Code:    code,
		Message: constant.GetMsg(code),
	}
}
