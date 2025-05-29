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
	"tiktok_e-commerce/app/product/biz/model"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CategoryDAO struct {
	db    *gorm.DB
	cache *redis.Client
}

var (
	instanceCategoryDAO *CategoryDAO
	onceCategoryDAO     sync.Once
)

func GetCategoryDAO() *CategoryDAO {
	return NewCategoryDAO()
}

func NewCategoryDAO() *CategoryDAO {
	onceCategoryDAO.Do(func() {
		dal.Init()
		db := mysql.DB
		cache := cacheRedis.RedisClient
		err := db.AutoMigrate(&model.Category{})
		if err != nil {
			klog.Error("failed to auto migrate category table: %v", err)
			return
		}
		instanceCategoryDAO = &CategoryDAO{
			db:    db,
			cache: cache,
		}
	})
	return instanceCategoryDAO
}

func (dao *CategoryDAO) DB() *gorm.DB {
	return dao.db
}

func (dao *CategoryDAO) Cache() *redis.Client {
	return dao.cache
}

func (dao *CategoryDAO) GetCategoryByName(ctx context.Context, categoryName string) (*model.Category, error) {
	if categoryName == "" {
		klog.Error("分类名称为空！")
		return nil, errors.New("分类名称为空！")
	}

	cacheKey := (&model.Category{Name: categoryName}).GetCacheKey()
	val, err := dao.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var redisResult model.Category
		if err := json.Unmarshal([]byte(val), &redisResult); err != nil {
			klog.Error("分类信息反序列化失败，categoryName:", categoryName, "错误信息:", err)
			return nil, fmt.Errorf("分类信息反序列化失败: %v", err)
		}
		klog.Info("从 Redis 缓存中成功获取分类信息，categoryName:", categoryName)
		return &redisResult, nil
	}

	if err != redis.Nil {
		klog.Error("Redis 查询失败，categoryName:", categoryName, "错误信息:", err)
		return nil, err
	}

	var category model.Category
	err = dao.db.WithContext(ctx).Where("name = ?", categoryName).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Info("分类不存在, categoryName:", categoryName)
			return nil, nil
		}
		klog.Error("数据库查询失败, categoryName:", categoryName, "错误信息:", err)
		return nil, err
	}

	categoryJSON, err := json.Marshal(category)
	if err != nil {
		klog.Error("分类信息序列化失败，categoryName:", categoryName, "错误信息:", err)
		return nil, fmt.Errorf("分类信息序列化失败: %v", err)
	}

	err = dao.cache.Set(ctx, cacheKey, categoryJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，categoryName:", categoryName, "错误信息:", err)
		return nil, fmt.Errorf("设置 Redis 缓存失败: %v", err)
	}
	return &category, nil
}

func (dao *CategoryDAO) GetOrCreateCategoryByName(ctx context.Context, categoryName string) (*model.Category, error) {
	if categoryName == "" {
		klog.Error("分类名称为空！")
		return nil, errors.New("分类名称为空！")
	}

	category, err := dao.GetCategoryByName(ctx, categoryName)
	if err == nil && category != nil {
		return category, nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var newCategory model.Category
	err = dao.db.WithContext(ctx).
		Where("name = ?", categoryName).
		FirstOrCreate(&newCategory, model.Category{Name: categoryName}).Error
	if err != nil {
		klog.Error("数据库查询失败, categoryName:", categoryName, "错误信息:", err)
		return nil, err
	}

	cacheKey := newCategory.GetCacheKey()
	categoryJSON, err := json.Marshal(newCategory)
	if err != nil {
		klog.Error("分类信息序列化失败，categoryName:", categoryName, "错误信息:", err)
		return nil, fmt.Errorf("分类信息序列化失败: %v", err)
	}
	err = dao.cache.Set(ctx, cacheKey, categoryJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，categoryName:", categoryName, "错误信息:", err)
		return nil, fmt.Errorf("设置 Redis 缓存失败: %v", err)
	}
	return &newCategory, nil
}

func (dao *CategoryDAO) DelUnusedCategory(ctx context.Context, category *model.Category) error {
	if category == nil {
		klog.Error("分类信息为空！")
		return errors.New("分类信息为空！")
	}

	var count int64
	err := dao.db.WithContext(ctx).Table("product_category").
		Where("product_category.category_id = ?", category.ID).Count(&count).Error
	if err != nil {
		klog.Error("数据库查询失败, categoryName:", category.Name, "错误信息:", err)
		return err
	}
	if count == 0 {
		err = dao.db.WithContext(ctx).Delete(category).Error
		if err != nil {
			klog.Error("数据库删除分类失败, categoryName:", category.Name, "错误信息:", err)
			return err
		}

		cacheKey := category.GetCacheKey()
		err = dao.cache.Del(ctx, cacheKey).Err()
		if err != nil {
			klog.Error("删除 Redis 缓存失败, categoryName:", category.Name, "错误信息:", err)
			return err
		}
	}
	return nil
}
