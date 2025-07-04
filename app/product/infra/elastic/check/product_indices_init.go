package check

import (
	"context"
	"fmt"
	"strings"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	"tiktok_e-commerce/product/biz/vo"
	"tiktok_e-commerce/product/infra/elastic/client"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func ProduceIndicesInit() {
	ctx := context.Background()
	// 构建请求
	productIndicesExist, err := esapi.IndicesExistsRequest{
		Index: []string{"product"},
	}.Do(ctx, client.ElasticClient)
	if err != nil {
		klog.Error(err)
		return
	}
	// 如果product不存在，就创建这个索引库
	if productIndicesExist.StatusCode != 200 {
		SettingData, err := sonic.Marshal(vo.ProductSearchMappingSetting)
		s := string(SettingData)
		fmt.Printf("%v", s)
		if err != nil {
			return
		}
		create, err := esapi.IndicesCreateRequest{
			Index: "product",
			Body:  strings.NewReader(s),
		}.Do(ctx, client.ElasticClient)
		if err != nil {
			klog.Info(err)
		}
		body := create.Body
		fmt.Printf("%v", body)
		if create.StatusCode != 200 {
			klog.Error("create product indices failed")
			return
		}
		// 将数据导入到product索引库中
		// 从数据库中获取数据

		products, err := model.SelectProductAllWithoutCondition(mysql.DB, ctx)
		if err != nil {
			klog.Errorf("查询数据库失败,err:%v", err)
			return
		}
		// 遍历数据，将数据转换为sonic格式
		for i := range products {
			pro := products[i]
			dataVo := vo.ProductSearchDataVo{
				Name:         pro.ProductName,
				Description:  pro.ProductDescription,
				ID:           pro.ProductId,
				CategoryName: pro.CategoryName,
			}
			sonicData, _ := sonic.Marshal(dataVo)
			// 调用esapi.BulkRequest将数据导入到product索引库中
			_, err = esapi.IndexRequest{
				Index:   "product",
				Body:    strings.NewReader(string(sonicData)),
				Refresh: "true",
			}.Do(ctx, client.ElasticClient)
			if err != nil {
				klog.Errorf("导入数据到索引库product失败,err:%v", err)
				return
			}
		}
	}
}
