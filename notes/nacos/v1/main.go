package main

import (
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	configClient, err := NewConfigClient()
	if err != nil {
		log.Fatalf("new config client error: %v", err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "dev-config.yaml",
		Group:  "USER",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("content:%s \n", content)
}

func NewServerConfig() []constant.ServerConfig {
	return []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}
}

func NewConfigClient() (config_client.IConfigClient, error) {
	sc := NewServerConfig()
	cc := constant.ClientConfig{
		TimeoutMs:            5000,
		CacheDir:             "cache",
		UpdateThreadNum:      0,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: false,
		LogDir:               "log",
		LogLevel:             "debug",
		NamespaceId:          "ab7e7974-9f47-452f-87e7-0849cb43c2f0",
	}
	return clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
}
