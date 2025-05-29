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

type ProductDAO struct {
	db    *gorm.DB
	cache *redis.Client
}

var (
	instanceProductDAO *ProductDAO
	onceProductDAO     sync.Once
)

func GetProductDAO() *ProductDAO {
	return NewProductDAO()
}

func NewProductDAO() *ProductDAO {
	onceProductDAO.Do(func() {
		dal.Init()
		db := mysql.DB
		cache := cacheRedis.RedisClient
		err := db.AutoMigrate(&model.Product{})
		if err != nil {
			klog.Error("failed to auto migrate product table: %v", err)
			return
		}
		instanceProductDAO = &ProductDAO{
			db:    db,
			cache: cache,
		}
	})
	return instanceProductDAO
}

func (dao *ProductDAO) DB() *gorm.DB {
	return dao.db
}

func (dao *ProductDAO) Cache() *redis.Client {
	return dao.cache
}

func (dao *ProductDAO) GetProductById(ctx context.Context, id uint32) (*model.Product, error) {
	if id == 0 {
		return nil, errors.New("商品id不能为空！")
	}
	cacheKey := (&model.Product{ID: id}).GetCacheKey()

	// 尝试从 Redis 缓存中获取商品信息
	val, err := dao.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var redisResult model.Product
		if err := json.Unmarshal([]byte(val), &redisResult); err != nil {
			klog.Error("产品信息反序列化失败，ID:", id, "错误信息:", err)
			return nil, fmt.Errorf("产品信息反序列化失败: %v", err)
		}
		klog.Info("从 Redis 缓存中成功获取商品信息，ID:", id)
		return &redisResult, nil
	}
	if err != redis.Nil {
		klog.Error("Redis 查询失败，ID:", id, "错误信息:", err)
		return nil, err
	}

	// 从数据库查询
	var product model.Product
	err = dao.db.WithContext(ctx).Preload("Categories").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Info("商品不存在，ID:", id)
			return nil, nil // 商品不存在，返回 nil, nil
		}
		klog.Error("数据库查询失败，ID:", id, "错误信息:", err)
		return nil, err // 其他错误，返回具体错误信息
	}

	// 存入redis缓存
	productJSON, err := json.Marshal(product)
	if err != nil {
		klog.Error("产品信息序列化失败，ID:", id, "错误信息:", err)
		return nil, fmt.Errorf("产品信息序列化失败: %v", err)
	}
	err = dao.cache.Set(ctx, cacheKey, productJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，ID:", id, "错误信息:", err)
		return nil, fmt.Errorf("设置 Redis 缓存失败: %v", err)
	}
	return &product, nil
}

func (dao *ProductDAO) GetProductsByQuery(ctx context.Context, query string) ([]*model.Product, error) {
	var products []*model.Product
	err := dao.db.WithContext(ctx).
		Where("MATCH(name, description) AGAINST(?)", query).
		Find(&products).Error
	if err != nil {
		klog.Error("数据库查询失败 Query:", query, "错误信息:", err)
		return nil, err
	}
	return products, nil
}

func (dao *ProductDAO) GetProductPage(ctx context.Context, page, pageSize uint32) ([]model.Product, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("页码和每页数量必须大于0！")
	}
	var products []model.Product
	err := dao.db.WithContext(ctx).
		Model(&model.Product{}).
		Order("id").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&products).Error
	if err != nil {
		klog.Error("数据库分页查询失败，页码:", page, "每页数量:", pageSize, "错误信息:", err)
		return nil, err
	}
	return products, nil
}

func (dao *ProductDAO) GetProductPageByCategory(ctx context.Context, page, pageSize uint32, category *model.Category) ([]model.Product, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("页码和每页数量必须大于0！")
	}
	var products []model.Product
	if category == nil || category.Name == "" {
		return dao.GetProductPage(ctx, page, pageSize)
	}
	cacheKey := fmt.Sprintf("product:category:%s:page:%d:pageSize:%d", category.Name, page, pageSize)

	// 尝试从 Redis 缓存中获取分页数据
	val, err := dao.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedProducts []model.Product
		if err := json.Unmarshal([]byte(val), &cachedProducts); err != nil {
			klog.Error("分页产品信息反序列化失败，Category:", category.Name, "错误信息:", err)
			return nil, fmt.Errorf("分页产品信息反序列化失败: %v", err)
		}
		klog.Info("从 Redis 缓存中成功获取分页产品信息，Category:", category.Name, "Page:", page)
		return cachedProducts, nil
	}
	if err != redis.Nil {
		klog.Error("Redis 查询失败，Category:", category.Name, "错误信息:", err)
		return nil, err
	}

	// 从数据库查询
	err = dao.db.WithContext(ctx).
		Joins("JOIN product_category ON product_category.product_id = product.id").
		Where("product_category.category_id = ?", +category.ID).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&products).Error

	if err != nil {
		klog.Error("数据库分页查询失败，Category:", category.Name, "错误信息:", err)
		return nil, err
	}

	// 存入Redis缓存
	productJSON, err := json.Marshal(products)
	if err != nil {
		klog.Error("分页产品信息序列化失败，Category:", category.Name, "错误信息:", err)
		return products, nil // 返回查询结果但不缓存
	}

	err = dao.cache.Set(ctx, cacheKey, productJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，Category:", category.Name, "错误信息:", err)
		return nil, fmt.Errorf("设置 Redis 缓存失败: %v", err)
	}
	return products, nil
}

func (dao *ProductDAO) CreateProduct(product *model.Product) uint32 {
	if product == nil || product.StoreId == 0 || product.Name == "" || product.Description == "" {
		klog.Error("产品信息无效，StoreId:", product.StoreId, " Name:", product.Name, " Description:", product.Description)
		return 0
	}
	err := dao.db.Create(product).Error
	if err != nil {
		klog.Error("数据库创建产品失败，Name:", product.Name, " 错误信息:", err)
		return 0
	}
	cacheKey := product.GetCacheKey()
	productJSON, err := json.Marshal(product)
	if err != nil {
		klog.Error("产品信息序列化失败，ID:", product.ID, " 错误信息:", err)
		return product.ID
	}
	err = dao.cache.Set(context.Background(), cacheKey, productJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，ID:", product.ID, " 错误信息:", err)
	}
	return product.ID
}

func (dao *ProductDAO) UpdateProduct(product *model.Product) error {
	if product == nil || product.StoreId == 0 || product.Name == "" || product.Description == "" {
		return errors.New("产品信息无效")
	}
	err := dao.db.WithContext(context.Background()).Model(&product).
		Association("Categories").
		Replace(product.Categories)
	if err != nil {
		klog.Error("更新产品分类关联失败，ID:", product.ID, " 错误信息:", err)
		return err
	}
	err = dao.db.WithContext(context.Background()).Save(product).Error
	if err != nil {
		klog.Error("数据库更新产品失败，ID:", product.ID, " 错误信息:", err)
		return err
	}
	cacheKey := product.GetCacheKey()
	productJSON, err := json.Marshal(product)
	if err != nil {
		klog.Error("产品信息序列化失败，ID:", product.ID, " 错误信息:", err)
		return err
	}
	err = dao.cache.Set(context.Background(), cacheKey, productJSON, time.Hour).Err()
	if err != nil {
		klog.Error("设置 Redis 缓存失败，ID:", product.ID, " 错误信息:", err)
		return err
	}
	return nil
}

func (dao *ProductDAO) DeleteProduct(product *model.Product) error {
	if product == nil {
		return errors.New("产品ID无效")
	}
	err := dao.db.WithContext(context.Background()).Delete(product).Error
	if err != nil {
		klog.Error("数据库删除产品失败，ID:", product.ID, " 错误信息:", err)
		return err
	}
	err = dao.cache.Del(context.Background(), product.GetCacheKey()).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		klog.Error("删除 Redis 缓存失败，ID:", product.ID, " 错误信息:", err)
		return err
	}
	return nil
}
