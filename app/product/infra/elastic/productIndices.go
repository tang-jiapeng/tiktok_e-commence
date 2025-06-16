package elastic

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	"tiktok_e-commerce/product/biz/vo"

	"github.com/bytedance/sonic"
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

// func ProduceIndicesInit() {
// 	// 构建请求
// 	productIndicesExist, err := esapi.IndicesExistsRequest{
// 		Index: []string{"product"},
// 	}.Do(context.Background(), &ElasticClient)
// 	if err != nil {
// 		klog.Error(err)
// 		return
// 	}
// 	// 如果product不存在，就创建这个索引库
// 	if productIndicesExist.StatusCode != 200 {
// 		create, err := esapi.IndicesCreateRequest{
// 			Index: "product",
// 			Body:  strings.NewReader(mapping),
// 		}.Do(context.Background(), &ElasticClient)
// 		if err != nil {
// 			klog.Info(err)
// 		}
// 		if create.StatusCode != 200 {
// 			klog.Error("创建商品索引失败")
// 			return
// 		}

// 		//将数据导入到product索引库中
// 		// 1 从数据库中获取数据
// 		var products []model.Product
// 		result := mysql.DB.Table("tb_product").Select("*").Find(&products)
// 		if result.Error != nil {
// 			klog.Error(result.Error)
// 			return
// 		}
// 		// 2 遍历数据，将数据转换为json格式
// 		for i := range products {
// 			pro := products[i]
// 			dataVo := vo.ProductSearchDataVo{
// 				Name:        pro.Name,
// 				Description: pro.Description,
// 			}
// 			jsonData, _ := json.Marshal(dataVo)
// 			// 3 调用esapi.BulkRequest将数据导入到product索引库中
// 			_, _ = esapi.IndexRequest{
// 				Index:   "product",
// 				Body:    strings.NewReader(string(jsonData)),
// 				Refresh: "true",
// 			}.Do(context.Background(), &ElasticClient)
// 		}
// 	}
// }

func ProduceIndicesInit() {
	// 检查 product 索引是否存在
	productIndicesExist, err := esapi.IndicesExistsRequest{
		Index: []string{"product"},
	}.Do(context.Background(), ElasticClient)
	if err != nil {
		klog.Errorf("检查索引存在失败: %v", err)
		return
	}

	// 如果索引不存在，创建索引
	if productIndicesExist.StatusCode != 200 {
		SettingData, err := sonic.Marshal(vo.ProductSearchMappingSetting)
		if err != nil {
			return
		}
		create, err := esapi.IndicesCreateRequest{
			Index: "product",
			Body:  strings.NewReader(string(SettingData)),
		}.Do(context.Background(), ElasticClient)
		if err != nil {
			klog.Errorf("创建索引失败: %v", err)
		}
		if create.StatusCode != 200 {
			body, _ := io.ReadAll(create.Body)
			klog.Errorf("创建商品索引失败，状态码: %d, 响应: %s", create.StatusCode, string(body))
			return
		}
		// klog.Info("成功创建 product 索引")
	}

	// 同步数据库数据到 product 索引
	if err := SyncProductToES(); err != nil {
		klog.Errorf("数据同步失败: %v", err)
	}
}

// SyncProductToES 同步 tb_product 表数据到 Elasticsearch
func SyncProductToES() error {
	// 1. 从数据库中获取数据
	var products []model.Product
	result := mysql.DB.Table("tb_product").Select("*").Find(&products)
	if result.Error != nil {
		klog.Errorf("查询数据库失败: %v", result.Error)
		return result.Error
	}

	// 打印查询到的数据
	// klog.Infof("从数据库查询到 %d 条产品数据", len(products))
	// for i, pro := range products {
	// 		klog.Infof("产品 %d: ID=%d, Name=%s, Description=%s, Picture=%s, Price=%s, Stock=%d, Sale=%f, PublicState=%d, LockStock=%d",
	// 				i+1, pro.ID, pro.Name, pro.Description, pro.Picture, pro.Price, pro.Stock, pro.Sale, pro.PublicState, pro.LockStock)
	// }

	// 2. 遍历数据，同步到 Elasticsearch
	for i := range products {
		pro := products[i]
		dataVo := vo.ProductSearchDataVo{
			Name:        pro.Name,
			Description: pro.Description,
		}
		jsonData, _ := json.Marshal(dataVo)

		// 3. 使用 IndexRequest 插入或更新文档
		_, _ = esapi.IndexRequest{
			Index:   "product",
			Body:    strings.NewReader(string(jsonData)),
			Refresh: "true", // 立即刷新，确保数据可见
		}.Do(context.Background(), ElasticClient)
	}

	return nil
}
