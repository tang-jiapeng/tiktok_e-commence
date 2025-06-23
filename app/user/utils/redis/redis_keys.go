package redis

import "fmt"

const (
	UserInfoCacheKey      = "user:info:%d"
	UserAddressesCacheKey = "user:addresses:%d"
)

func GetUserKey(userId int32) string {
	return fmt.Sprintf(UserInfoCacheKey, userId)
}

func GetUserAddressesKey(userId int32) string {
	return fmt.Sprintf(UserAddressesCacheKey, userId)
}
