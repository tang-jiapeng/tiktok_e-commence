package elastic

import (
	"tiktok_e-commerce/product/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

var (
	ElasticClient elasticsearch7.Client
)

func InitClient() {
	elasticsearchConf := conf.GetConf().Elasticsearch
	var client *elasticsearch7.Client
	client, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{"http://" + elasticsearchConf.Host + ":" + elasticsearchConf.Port},
		Username:  elasticsearchConf.Username,
		Password:  elasticsearchConf.Password,
	})
	ElasticClient = *client
	if err != nil {
		klog.Error(err)
		return
	}
	ProduceIndicesInit()
}
