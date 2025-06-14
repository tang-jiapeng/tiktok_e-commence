package service

import (
	"context"
	"io"
	"strings"

	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/vo"
	"tiktok_e-commerce/product/infra/elastic"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/common/json"
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
			MutiMatch: vo.ProductSearchMultiMatchQuery{
				Query:  req.Query,
				Fields: []string{"name", "description"},
			},
		},
	}
	jsonData, _ := json.Marshal(queryBody)
	//发往elastic
	//TODO 将关键词发往elastic，检索数据
	search, _ := esapi.SearchRequest{
		Index: []string{"product"},
		Body:  strings.NewReader(string(jsonData)),
	}.Do(context.Background(), &elastic.ElasticClient)
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
	var products = []*product.Product{}
	for i := range productHitsList {
		productData := productHitsList[i].Source
		pro := product.Product{
			Name:        productData.Name,
			Description: productData.Description,
		}
		products = append(products, &pro)
	}
	//TODO 将返回的数据返回到前端
	resp = &product.SearchProductsResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Results:    products,
	}
	return resp, nil
}
