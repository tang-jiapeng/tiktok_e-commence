package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"tiktok_e-commerce/common/infra/hot_key_client/constants"
	keyModel "tiktok_e-commerce/common/infra/hot_key_client/model/key"
	"tiktok_e-commerce/common/infra/hot_key_client/processor"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/dal/redis"
	"tiktok_e-commerce/product/biz/model"
	"tiktok_e-commerce/product/infra/elastic/client"
	"time"

	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/vo"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

type SearchProductsService struct {
	ctx context.Context
} // NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

// Run create note info
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	queryBody := vo.ProductSearchQueryBody{
		Query: &vo.ProductSearchQuery{
			Bool: &vo.ProductSearchBoolQuery{
				Should: &[]*vo.ProductSearchQuery{
					{
						MultiMatch: &vo.ProductSearchMultiMatchQuery{
							Query:  req.Query,
							Fields: []string{"name", "description", "category_name"},
						},
					},
					//{
					//	MultiMatch: &vo.ProductSearchMultiMatchQuery{
					//		Query:  req.CategoryName,
					//		Fields: []string{"category_name"},
					//	},
					//},
				},
			},
			//MultiMatch: &vo.ProductSearchMultiMatchQuery{
			//	Query:  req.Query,
			//	Fields: []string{"name", "description"},
			//},
		},
		Source: &vo.ProductSearchSource{
			"id",
		},
	}
	dslBytes, _ := sonic.Marshal(queryBody)
	//将dsl计算hashcode
	hasher := md5.New()
	hasher.Write(dslBytes)
	md5Bytes := hasher.Sum(nil)
	//从redis查找该hashcode对应的缓存数据
	dslKey := "product:dslBytes:" + string(md5Bytes)
	keyConf := keyModel.NewKeyConf1(constants.ProductService)
	hotKeyModel := keyModel.NewHotKeyModelWithConfig(dslKey, &keyConf)
	var dslCache string
	localDslCache := processor.GetValue(*hotKeyModel)
	if localDslCache != nil {
		dslCache = localDslCache.(string)
	}

	//搜索返回的id
	var searchIds []int64

	if dslCache != "" && dslCache != "null" {
		err = sonic.UnmarshalString(dslCache, &searchIds)
		if err != nil {
			klog.CtxErrorf(s.ctx, "dsl缓存反序列化失败, err: %v", err)
			return
		}
	} else {
		//如果缓存数据存在，则直接返回数据
		//发往elastic
		search, err := esapi.SearchRequest{
			Index: []string{"product"},
			Body:  strings.NewReader(string(dslBytes)),
		}.Do(s.ctx, client.ElasticClient)
		if err != nil {
			klog.CtxErrorf(s.ctx, "es搜索失败, err: %v", err)
			return nil, errors.WithStack(err)
		}
		// 解析数据
		searchData, err := io.ReadAll(search.Body)
		if err != nil {
			klog.CtxErrorf(s.ctx, "es搜索结果解析失败, err: %v", err)
			return nil, errors.WithStack(err)
		}
		elasticSearchVo := vo.ProductSearchAllDataVo{}
		err = json.Unmarshal(searchData, &elasticSearchVo)
		if err != nil {
			klog.CtxErrorf(s.ctx, "es搜索结果反序列化失败, err: %v", err)
			return nil, errors.WithStack(err)
		}
		productHitsList := elasticSearchVo.Hits.Hits
		if len(productHitsList) > 0 {
			for i := range productHitsList {
				sourceData := productHitsList[i].Source
				searchIds = append(searchIds, sourceData.ID)
			}

			//将es返回的数据缓存到redis
			dslCacheToRedis, _ := sonic.Marshal(searchIds)
			_, err = redis.RedisClient.Set(context.Background(), dslKey, dslCacheToRedis, 5*time.Minute).Result()
			if err != nil {
				klog.CtxErrorf(s.ctx, "dsl搜索结果缓存到redis失败, err: %v", err)
				return nil, errors.WithStack(err)
			}
			marshalString, err := sonic.Marshal(searchIds)
			if err != nil {
				klog.CtxErrorf(s.ctx, "dsl搜索结果序列化失败, err: %v", err)
				return nil, errors.WithStack(err)
			}
			processor.SmartSet(*hotKeyModel, marshalString)
		}
	}
	var products []*product.Product = make([]*product.Product, 0)
	//根据id从缓存或者数据库中获取数据
	//keys是redis的key
	var keys []string = make([]string, 0, len(searchIds))
	//根据返回的数据确认是否有缺失数据，有的话把当前的id存进去
	var missingIds []int64
	//将id转换为redis的key
	for i := range searchIds {
		keys = append(keys, "product:"+strconv.FormatInt(searchIds[i], 10))
	}
	//先判断redis是否存在数据，如果存在，则直接返回数据
	if len(keys) > 0 {
		//加入hotkey
		var keysKey = "product:keys:" + string(md5Bytes)
		keysConf := keyModel.NewKeyConf1(constants.ProductService)
		keysHotkeyModel := keyModel.NewHotKeyModelWithConfig(keysKey, &keysConf)
		var values []interface{}
		localKeysCache := processor.GetValue(*keysHotkeyModel)
		if localKeysCache != nil {
			values = localKeysCache.([]interface{})
		} else {
			values, err = redis.RedisClient.MGet(context.Background(), keys...).Result()
			if err != nil {
				klog.CtxErrorf(s.ctx, "redis获取dsl搜索结果失败, err: %v", err)
				return nil, errors.WithStack(err)
			}
			processor.SmartSet(*keysHotkeyModel, values)
		}
		for i, value := range values {
			// 提取ID部分
			idStr := keys[i][len("product:"):]
			if value == nil {
				if id, err := strconv.ParseInt(idStr, 10, 64); err != nil {
					missingIds = append(missingIds, id)
				} else {
					klog.CtxErrorf(s.ctx, "productid转换失败, err: %v", err)
					return nil, errors.WithStack(err)
				}
			} else {
				//解析数据
				productData := product.Product{}
				err = sonic.UnmarshalString(value.(string), &productData)
				if err != nil {
					klog.CtxErrorf(s.ctx, "product数据反序列化失败, err: %v", err)
					return nil, errors.WithStack(err)
				}
				products = append(products, &productData)
			}
		}
	}
	//如果不存在，则从数据库中获取数据，并存入redis
	if len(missingIds) > 0 {
		//从数据库中获取数据
		list, err := model.SelectProductList(mysql.DB, context.Background(), missingIds)
		if err != nil {
			klog.CtxErrorf(s.ctx, "从数据库中获取数据失败, err: %v", err)
			return nil, errors.WithStack(err)
		}
		pipeline := redis.RedisClient.Pipeline()

		missingProducts := make([]*product.Product, len(list))
		for i := range list {
			p := product.Product{
				Id:          list[i].ProductId,
				Name:        list[i].ProductName,
				Description: list[i].ProductDescription,
				Picture:     list[i].ProductPicture,
				Price:       list[i].ProductPrice,
				Stock:       list[i].ProductStock,
				Sale:        list[i].ProductSale,
				CategoryId:  list[i].CategoryID,
			}
			missingProducts[i] = &p
			productKey := "product:" + strconv.FormatInt(list[i].ProductId, 10)
			marshalString, err := sonic.MarshalString(p)
			if err != nil {
				klog.CtxErrorf(s.ctx, "product数据序列化失败, err: %v", err)
				return nil, err
			}
			pipeline.Set(context.Background(), productKey, marshalString, 1*time.Hour)
		}
		products = append(products, missingProducts...)
		//存入redis
		_, err = pipeline.Exec(context.Background())
		if err != nil {
			klog.CtxErrorf(s.ctx, "product数据存入redis失败, err: %v", err)
			return nil, err
		}
	}

	resp = &product.SearchProductsResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Results:    products,
	}
	return resp, nil
}
