package elastic

import (
	"tiktok_e-commerce/product/infra/elastic/check"
	"tiktok_e-commerce/product/infra/elastic/client"
)

func InitClient() {
	client.InitClient()
	check.ProduceIndicesInit()
}
