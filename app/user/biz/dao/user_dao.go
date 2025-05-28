package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"tiktok_e-commerce/app/hertz/biz/dal"
	"tiktok_e-commerce/app/hertz/biz/dal/mysql"
	cacheRedis "tiktok_e-commerce/app/hertz/biz/dal/redis"
	"tiktok_e-commerce/app/user/biz/model"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserDAO struct {
	db    *gorm.DB
	cache *redis.Client
}

var (
	instanceDAO *UserDAO
	onceDAO     sync.Once
)

func GetUserDAO() *UserDAO {
	return NewUserDAO()
}

func NewUserDAO() *UserDAO {
	onceDAO.Do(func() {
		dal.Init()
		db := mysql.DB
		cache := cacheRedis.RedisClient
		err := db.AutoMigrate(&model.User{})
		if err != nil {
			klog.Error("failed to auto migrate user table: %v")
			return
		}
		instanceDAO = &UserDAO{
			db:    db,
			cache: cache,
		}
	})
	return instanceDAO
}

func (dao *UserDAO) DB() *gorm.DB {
	return dao.db
}

func (dao *UserDAO) Cache() *redis.Client {
	return dao.cache
}

func (dao *UserDAO) Insert(user *model.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDAO) FindOne(ctx context.Context, userId int64) (*model.User, error) {
	q := &model.User{
		UserId: userId,
	}
	cacheKey := q.GetCacheKey()

	// 尝试从 Redis 缓存中获取用户信息
	val, err := dao.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var redisResult model.User
		err := json.Unmarshal([]byte(val), &redisResult)
		if err != nil {
			klog.Error("用户信息反序列化失败，UserId:", userId, " 错误信息:", err)
			return nil, fmt.Errorf("用户信息反序列化失败: %v", err)
		}
		klog.Info("从 Redis 缓存中成功获取用户信息，UserId:", userId)
		return &redisResult, nil
	}
	if err != redis.Nil {
		klog.Error("Redis 查询失败，UserId:", userId, " 错误信息:", err)
		return nil, err
	}

	// 如果缓存中没有，从数据库中查询
	var user model.User
	err = dao.db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Info("用户不存在，UserId:", userId)
			return nil, nil // 用户不存在，返回 nil, nil
		}
		klog.Error("数据库查询失败，UserId:", userId, " 错误信息:", err)
		return nil, err // 其他错误，返回具体错误信息
	}

	// 将查询到的用户信息存入 Redis 缓存
	userJSON, err := json.Marshal(user)
	if err != nil {
		klog.Error("用户信息序列化失败，UserId:", userId, " 错误信息:", err)
		return nil, fmt.Errorf("用户信息序列化失败: %v", err)
	}

	err = dao.cache.Set(ctx, cacheKey, userJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，UserId:", userId, " 错误信息:", err)
		return nil, fmt.Errorf("设置 Redis 缓存失败: %v", err)
	}

	return &user, nil
}

func (dao *UserDAO) Update(user *model.User) error {
	err := dao.db.Save(user).Error
	if err != nil {
		return err
	}
	cacheKey := user.GetCacheKey()

	userJson, err := json.Marshal(user)
	if err != nil {
		klog.Error("用户信息序列化失败，UserId:", user.UserId, " 错误信息:", err)
		return fmt.Errorf("用户信息序列化失败: %v", err)
	}

	err = dao.cache.Set(context.Background(), cacheKey, userJson, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，UserId:", user.UserId, " 错误信息:", err)
		return fmt.Errorf("设置 Redis 缓存失败: %v", err)
	}
	return nil
}

func (dao *UserDAO) Delete(userId int64) error {
	user := &model.User{
		UserId: userId,
	}

	// 使用 GORM 的软删除
	err := dao.db.Delete(user).Error
	if err != nil {
		return err
	}

	err = dao.cache.Del(context.Background(), user.GetCacheKey()).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

func (dao *UserDAO) FindByEmail(email string) (*model.User, error) {
	// 尝试从 Redis 缓存中获取用户信息
	cacheKey := fmt.Sprintf("user:email:%s", email)
	val, err := dao.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cachedUser model.User
		err := json.Unmarshal([]byte(val), &cachedUser)
		if err != nil {
			klog.Error("用户信息反序列化失败，Email:", email, " 错误信息:", err)
			return nil, fmt.Errorf("用户信息反序列化失败: %v", err)
		}
		return &cachedUser, nil
	}
	if err != redis.Nil {
		klog.Error("Redis 查询失败:", err)
		return nil, err
	}

	// 如果缓存中没有，从数据库中查询
	var user model.User
	err = dao.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在，返回 nil, nil
		}
		klog.Error("数据库查询失败:", err)
		return nil, err // 其他错误，返回具体错误信息
	}

	// 将查询到的用户信息存入 Redis 缓存
	userJSON, err := json.Marshal(user)
	if err != nil {
		klog.Error("用户信息序列化失败，Email:", email, " 错误信息:", err)
		return nil, fmt.Errorf("用户信息序列化失败: %v", err)
	}

	err = dao.cache.Set(context.Background(), cacheKey, userJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败:", err)
	}
	return &user, nil
}

func (dao *UserDAO) FindByUsername(username string) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:username:%s", username)
	val, err := dao.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cachedUser model.User
		err := json.Unmarshal([]byte(val), &cachedUser)
		if err != nil {
			klog.Error("用户信息反序列化失败，Username:", username, " 错误信息:", err)
			return nil, fmt.Errorf("用户信息反序列化失败: %v", err)
		}
		return &cachedUser, nil
	}
	if err != redis.Nil {
		klog.Error("Redis 查询失败:", err)
		return nil, err
	}

	var user model.User
	err = dao.db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在，返回 nil, nil
		}
		klog.Error("数据库查询失败:", err)
		return nil, err // 其他错误，返回具体错误信息
	}

	err = dao.cache.Set(context.Background(), cacheKey, user, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败:", err)
	}

	return &user, nil
}

// GetUserPermissions 根据用户ID查询用户权限
func (dao *UserDAO) GetUserPermissions(ctx context.Context, userId int64) (int32, error) {
	user, err := dao.FindOne(ctx, userId)
	if err != nil {
		return -1, err
	}
	if user == nil {
		return -1, fmt.Errorf("用户不存在")
	}

	return user.UserPermissions, nil
}
