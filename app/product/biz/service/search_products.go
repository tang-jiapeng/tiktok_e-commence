package service

import (
	"context"
	"io"
	"strconv"
	"strings"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/dal/redis"
	"tiktok_e-commerce/product/biz/model"

	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/vo"
	"tiktok_e-commerce/product/infra/elastic"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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
		Query: vo.ProductSearchQuery{
			MultiMatch: vo.ProductSearchMultiMatchQuery{
				Query:  req.Query,
				Fields: []string{"name", "description"},
			},
		},
		Source: vo.ProductSearchSource{
			"id",
		},
	}
	jsonData, _ := json.Marshal(queryBody)
	//发往elastic
	search, _ := esapi.SearchRequest{
		Index: []string{"product"},
		Body:  strings.NewReader(string(jsonData)),
	}.Do(context.Background(), elastic.ElasticClient)
	// 解析数据
	searchData, _ := io.ReadAll(search.Body)
	elasticSearchVo := vo.ProductSearchAllDataVo{}
	err = json.Unmarshal(searchData, &elasticSearchVo)
	if err != nil {
		resp = &product.SearchProductsResp{
			StatusCode: 2013,
			StatusMsg:  constant.GetMsg(2013),
		}
		return
	}
	productHitsList := elasticSearchVo.Hits.Hits
	var searchIds []int64
	for i := range productHitsList {
		sourceData := productHitsList[i].Source
		searchIds = append(searchIds, sourceData.ID)
	}
	var products []*product.Product
	//根据id从缓存或者数据库中获取数据
	//keys是redis的key
	var keys []string
	//根据返回的数据确认是否有缺失数据，有的话把当前的id存进去
	var missingIds []int64
	//将id转换为redis的key
	for i := range searchIds {
		keys = append(keys, "product:"+strconv.FormatInt(searchIds[i], 10))
	}
	//先判断redis是否存在数据，如果存在，则直接返回数据
	values, err := redis.RedisClient.MGet(s.ctx, keys...).Result()
	if err != nil {
		return
	}
	for i, value := range values {
		// 提取ID部分
		idStr := keys[i][len("product:"):]
		if value == nil {
			if id, err := strconv.ParseInt(idStr, 10, 64); err != nil {
				missingIds = append(missingIds, id)
			}
		} else {
			//解析数据
			productData := product.Product{}
			marshal, err := sonic.Marshal(value)
			if err != nil {
				return
			}
			err = json.Unmarshal(marshal, &productData)
			if err != nil {
				return
			}
			products = append(products, &productData)
		}
	}
	//如果不存在，则从数据库中获取数据，并存入redis
	if len(missingIds) > 0 {
		//从数据库中获取数据
		if list, err := model.SelectProductList(mysql.DB, s.ctx, missingIds); err == nil {
			missingProducts := make([]*product.Product, len(list))
			var missingKeysAndValues map[string]product.Product
			for i := range list {
				p := product.Product{
					Id:          list[i].ID,
					Name:        list[i].Name,
					Description: list[i].Description,
					Picture:     list[i].Picture,
					Price:       list[i].Price,
					Stock:       list[i].Stock,
					Sale:        list[i].Sale,
					CategoryId:  list[i].CategoryId,
					BrandId:     list[i].BrandId,
				}
				missingProducts[i] = &p
				missingKeysAndValues["product:"+strconv.FormatInt(list[i].ID, 10)] = p
			}
			products = append(products, missingProducts...)
			//存入redis
			if _, err := redis.RedisClient.MSet(s.ctx, missingKeysAndValues).Result(); err != nil {
				klog.Error("MSet products err:", err)
				return
			}
		} else {
			return
		}
	}
	resp = &product.SearchProductsResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Results:    products,
	}
	return resp, nil
}
