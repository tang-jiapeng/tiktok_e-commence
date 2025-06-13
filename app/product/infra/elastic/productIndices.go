package elastic

import (
	"context"
	"encoding/json"
	"strings"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	"tiktok_e-commerce/product/biz/vo"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var mapping = `{
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "analyzer": "ik_smart"
      },
      "description": {
        "type": "text"
      }
    }
  }
}
`

func ProduceIndicesInit() {
	// 构建请求
	productIndicesExist, err := esapi.IndicesExistsRequest{
		Index: []string{"product"},
	}.Do(context.Background(), &ElasticClient)
	if err != nil {
		klog.Error(err)
		return
	}
	// 如果product不存在，就创建这个索引库
	if productIndicesExist.StatusCode != 200 {
		create, err := esapi.IndicesCreateRequest{
			Index: "product",
			Body:  strings.NewReader(mapping),
		}.Do(context.Background(), &ElasticClient)
		if err != nil {
			klog.Info(err)
		}
		if create.StatusCode != 200 {
			klog.Error("创建商品索引失败")
			return
		}

		//将数据导入到product索引库中
		// 1 从数据库中获取数据
		var products []model.Product
		result := mysql.DB.Table("tb_product").Select("*").Find(&products)
		if result.Error != nil {
			klog.Error(result.Error)
			return
		}
		// 2 遍历数据，将数据转换为json格式
		for i := range products {
			pro := products[i]
			dataVo := vo.ProductSearchDataVo{
				Name:        pro.Name,
				Description: pro.Description,
			}
			jsonData, _ := json.Marshal(dataVo)
			// 3 调用esapi.BulkRequest将数据导入到product索引库中
			_, _ = esapi.IndexRequest{
				Index:   "product",
				Body:    strings.NewReader(string(jsonData)),
				Refresh: "true",
			}.Do(context.Background(), &ElasticClient)
		}
	}
}
