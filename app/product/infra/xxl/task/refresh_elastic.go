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
	"tiktok_e-commerce/product/infra/elastic/client"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/xxl-job/xxl-job-executor-go"
)

func RefreshElastic(ctx context.Context, param *xxl.RunReq) string {
	index := param.BroadcastIndex
	total := param.BroadcastTotal
	klog.CtxInfof(ctx, "刷新Elastic开始 CheckAccountTask start")
	err := refresh(ctx, index, total)
	if err != nil {
		klog.Errorf("刷新Elastic失败 CheckAccountTask failed, err: %v", err)
		return err.Error()
	}
	return "refresh elastic success"
}

func refresh(ctx context.Context, index, total int64) (err error) {
	// 从数据库获取数据
	allProduct, err := model.SelectProductAll(mysql.DB, ctx, index, total)
	if err != nil {
		klog.Errorf("从数据库获取数据失败 CheckAccountTask failed, err: %v", err)
		return err
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
		klog.Errorf("序列化es查询参数失败 CheckAccountTask failed, err: %v", err)
		return err
	}
	searchIdResponse, err := esapi.SearchRequest{
		Index: []string{"product"},
		Body:  strings.NewReader(string(searchIdBytes)),
	}.Do(ctx, client.ElasticClient)
	// 解析数据
	searchIdResponseBytes, err := io.ReadAll(searchIdResponse.Body)
	if err != nil {
		klog.Errorf("解析es查询结果失败 CheckAccountTask failed, err: %v", err)
		return err
	}
	elasticSearchVo := vo.ProductSearchAllDataVo{}
	err = sonic.Unmarshal(searchIdResponseBytes, &elasticSearchVo)
	if err != nil {
		klog.Errorf("反序列化es查询结果失败 CheckAccountTask failed, err: %v", err)
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
			klog.Errorf("序列化es更新参数失败 CheckAccountTask failed, err: %v", err)
			return err
		}
		docBytes, err := sonic.Marshal(doc)
		if err != nil {
			klog.Errorf("序列化es更新参数失败 CheckAccountTask failed, err: %v", err)
			return err
		}
		bulkBody = append(bulkBody, updateBytes...)
		bulkBody = append(bulkBody, docBytes...)
	}
	// 更新ElasticSearch
	bulkResponse, err := esapi.BulkRequest{
		Index: "product",
		Body:  bytes.NewBuffer(bulkBody),
	}.Do(ctx, client.ElasticClient)
	if err != nil {
		klog.Errorf("批量刷新es失败 CheckAccountTask failed, err: %v", err)
		return err
	}
	fmt.Printf("%v", bulkResponse)
	if bulkResponse.StatusCode != 200 {
		klog.Errorf("批量刷新es失败 CheckAccountTask failed, err: %v", err)
		return
	}
	klog.CtxInfof(ctx, "刷新Elastic成功,index:%d,total:%d", index, total)
	return
}
