package main

import (
	"context"
	"flag"
	"github.com/CharLemAznable/violet"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
)

func main() {
	initConfigFile()
	initEtcdClient()
	dataPlane := startDataPlane()
	cleanFn := startDaemon(dataPlane)
	defer cleanFn()

	ctrlServer := &http.Server{Addr: ":22920",
		Handler: violet.NewCtrlPlane(dataPlane)}
	_ = ctrlServer.ListenAndServe()
}

var configFile string

func initConfigFile() {
	flag.StringVar(&configFile, "configFile",
		"config.toml", "config file path")
	flag.Parse()
}

var (
	client    *etcd.Client
	configKey string
)

func initEtcdClient() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Default().Fatalf(
			"Read etcd config file: [%s], error: %v\n", configFile, err)
		return
	}
	config, err := LoadConfig(string(data))
	if err != nil {
		log.Default().Fatalf(
			"Load etcd config data error: %v, config data: \n%s\n", err, string(data))
		return
	}
	log.Default().Printf("Load etcd config data: \n%s\n", string(data))

	etcdConfig, err := ParseEtcdClientConfig(config.EtcdClient)
	if err != nil {
		log.Default().Fatalf(
			"Parse etcd config error: %v\n", err)
		return
	}
	clt, err := etcd.New(*etcdConfig)
	if err != nil {
		log.Default().Fatalf(
			"New etcd client error: %v\n", err)
		return
	}
	client = clt

	configKey = config.EtcdConfigKey
	if configKey == "" {
		configKey = "violet.default"
	}
}

func startDataPlane() violet.DataPlane {
	dataPlane := violet.NewDataPlane(initVioletConfig())
	go func() {
		dataServer := &http.Server{Addr: ":22915", Handler: dataPlane}
		_ = dataServer.ListenAndServe()
	}()
	return dataPlane
}

func initVioletConfig() *violet.Config {
	response, err := client.Get(context.Background(), configKey)
	if err != nil {
		log.Default().Printf(
			"Get config by key: [%s], error: %v", configKey, err)
		return &violet.Config{}
	}
	if len(response.Kvs) != 1 {
		log.Default().Printf(
			"Get config by key: [%s], failed with response:\n%s", configKey, response)
		return &violet.Config{}
	}
	return loadVioletConfig(response.Kvs[0].Value)
}

func loadVioletConfig(data []byte) *violet.Config {
	config, err := violet.LoadConfig(string(data))
	if err != nil {
		log.Default().Printf(
			"Load config data error: %v, config data: \n%s\n=====", err, string(data))
		return &violet.Config{}
	}
	log.Default().Printf("Load config data: \n%s\n=====", string(data))
	return config
}

func startDaemon(dataPlane violet.DataPlane) func() {
	watchChan := client.Watch(context.Background(), configKey)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case response := <-watchChan:
				for _, event := range response.Events {
					dataPlane.SetConfig(loadVioletConfig(event.Kv.Value))
				}
			case <-quit:
				return
			}
		}
	}()
	return func() {
		close(quit)
	}
}
