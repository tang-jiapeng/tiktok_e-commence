package conf

import (
	"fmt"
	"os"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/config-nacos/nacos"
	"github.com/kr/pretty"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env           string
	Kitex         Kitex         `yaml:"kitex"`
	MySQL         MySQL         `yaml:"mysql"`
	Redis         Redis         `yaml:"redis"`
	Registry      Registry      `yaml:"registry"`
	Elasticsearch ElasticSearch `yaml:"elasticsearch"`
	Kafka         Kafka         `yaml:"kafka"`
	XxlJob        XxlJob        `yaml:"xxl_job"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Kitex struct {
	Service       string `yaml:"service"`
	Address       string `yaml:"address"`
	MetricsPort   string `yaml:"metrics_port"`
	LogLevel      string `yaml:"log_level"`
	LogFileName   string `yaml:"log_file_name"`
	LogMaxSize    int    `yaml:"log_max_size"`
	LogMaxBackups int    `yaml:"log_max_backups"`
	LogMaxAge     int    `yaml:"log_max_age"`
}

type Registry struct {
	RegistryAddress string `yaml:"registry_address"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
}

type Kafka struct {
	BizKafka BizKafka `yaml:"biz_kafka"`
}

type BizKafka struct {
	BootstrapServers string `yaml:"bootstrap_servers"`
	ProductTopicId   string `yaml:"product_topic_id"`
}

type ElasticSearch struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type XxlJob struct {
	XxlJobAddress string `yaml:"xxl_job_address"`
	ExecutorIp    string `yaml:"executor_ip"`
	AccessToken   string `yaml:"access_token"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	client, err := nacos.NewClient(nacos.Options{
		Address:     os.Getenv("NACOS_ADDR"),
		NamespaceID: "e45ccc29-3e7d-4275-917b-febc49052d58",
		Group:       "DEFAULT_GROUP",
		Username:    "nacos",
		Password:    os.Getenv("NACOS_PASSWORD"),
		Port:        8848,
	})
	if err != nil {
		panic(err)
	}
	param := vo.ConfigParam{
		DataId: "product_conf.yaml",
		Group:  "DEFAULT_GROUP",
		Type:   "yaml",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Printf("Config changed - namespace: %s, group: %s, data-id: %s\n", namespace, group, dataId)

			// 解析 YAML 配置
			var config interface{}
			err := yaml.Unmarshal([]byte(data), &config)
			if err != nil {
				fmt.Printf("Error parsing YAML: %v\n", err)
				return
			}

			// 输出解析结果
			fmt.Printf("Parsed YAML: %v\n", config)
		},
	}

	client.RegisterConfigCallback(param, func(data string, parser nacos.ConfigParser) {
		// 处理配置数据的逻辑
		if conf == nil {
			conf = new(Config)
		}
		err := yaml.Unmarshal([]byte(data), &conf)
		if err != nil {
			klog.Error("Error parsing YAML: %v\n", err)
			return
		}
		_, err = pretty.Printf("%+v\n", conf)
		if err != nil {
			klog.Error("pretty print error - %v", err)
		}
	}, 5000)
	conf.Env = GetEnv()
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}

func LogLevel() klog.Level {
	level := GetConf().Kitex.LogLevel
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
