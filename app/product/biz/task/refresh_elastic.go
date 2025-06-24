package task

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	"tiktok_e-commerce/product/biz/vo"
	"tiktok_e-commerce/product/infra/elastic"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/xxl-job/xxl-job-executor-go"
)

func RefreshElastic(ctx context.Context, param *xxl.RunReq) string {
	klog.CtxInfof(ctx, "refresh elastic start")
	err := refresh(ctx)
	if err != nil {
		return err.Error()
	}
	return "refresh elastic success"
}

func refresh(ctx context.Context) (err error) {
	// 从数据库获取数据
	allProduct, err := model.SelectProductAll(mysql.DB, ctx)
	if err != nil {
		return
	}
	productMap := map[int64]model.Product{}
	for i := range allProduct {
		productMap[allProduct[i].ID] = allProduct[i]
	}
	// 从es获取文档的id
	queryBody := vo.ProductSearchQueryBody{
		Query: &vo.ProductSearchQuery{
			MultiMatch: &vo.ProductSearchMultiMatchQuery{},
		},
		Source: &vo.ProductSearchSource{
			"id",
		},
	}
	searchIdBytes, err := sonic.Marshal(queryBody)
	if err != nil {
		return
	}
	searchIdResponse, err := esapi.SearchRequest{
		Index: []string{"product"},
		Body:  strings.NewReader(string(searchIdBytes)),
	}.Do(ctx, elastic.ElasticClient)
	// 解析数据
	searchIdResponseBytes, _ := io.ReadAll(searchIdResponse.Body)
	elasticSearchVo := vo.ProductSearchAllDataVo{}
	err = sonic.Unmarshal(searchIdResponseBytes, &elasticSearchVo)
	if err != nil {
		return
	}
	hits := elasticSearchVo.Hits.Hits
	var bulkBody []byte
	for i := range hits {
		source := hits[i]
		p := productMap[source.Source.ID]
		update := vo.ProductBulkUpdate{
			Update: vo.ProductBulkBody{
				DocID: source.ID,
			},
		}
		doc := vo.ProductBulkDoc{
			Doc: vo.ProductSearchDoc{
				Name:        p.Name,
				Description: p.Description,
			},
		}
		updateBytes, err := sonic.Marshal(update)
		if err != nil {
			return err
		}
		docBytes, err := sonic.Marshal(doc)
		if err != nil {
			return err
		}
		bulkBody = append(bulkBody, updateBytes...)
		bulkBody = append(bulkBody, docBytes...)
	}
	// 更新ElasticSearch
	bulkResponse, err := esapi.BulkRequest{
		Index: "product",
		Body:  bytes.NewBuffer(bulkBody),
	}.Do(ctx, elastic.ElasticClient)
	if err != nil {
		return err
	}
	fmt.Printf("%v", bulkResponse)
	if bulkResponse.StatusCode != 200 {
		return
	}
	return
}
