package del

import (
	"context"
	"fmt"
	"tiktok_e-commerce/product/infra/elastic/client"
	"tiktok_e-commerce/product/infra/kafka/model"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func DeleteProductToElasticSearch(ctx context.Context, product *model.DeleteProductSendToKafka) (err error) {

	// 调用esapi.BulkRequest将数据导入到product索引库中
	response, _ := esapi.DeleteRequest{
		Index:      "product",
		DocumentID: fmt.Sprintf("%d", product.ID),
	}.Do(ctx, client.ElasticClient)
	print(response.StatusCode)
	return nil
}
