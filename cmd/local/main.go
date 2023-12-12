package main

import (
	"flag"
	"github.com/CharLemAznable/violet"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	initConfigFile()
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

func startDataPlane() violet.DataPlane {
	dataPlane := violet.NewDataPlane(loadConfig())
	go func() {
		dataServer := &http.Server{Addr: ":22915", Handler: dataPlane}
		_ = dataServer.ListenAndServe()
	}()
	return dataPlane
}

func loadConfig() *violet.Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Default().Printf(
			"Read config file: [%s], error: %v", configFile, err)
		return &violet.Config{}
	}
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
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if updateLastModified() {
					dataPlane.SetConfig(loadConfig())
				}
			case <-quit:
				return
			}
		}
	}()
	return func() {
		ticker.Stop()
		close(quit)
	}
}

var lastModifiedMux sync.Mutex
var lastModified = time.Now()

func updateLastModified() bool {
	lastModifiedMux.Lock()
	defer lastModifiedMux.Unlock()
	fileInfo, err := os.Stat(configFile)
	if err != nil {
		return false
	}
	modTime := fileInfo.ModTime()
	if modTime.After(lastModified) {
		lastModified = modTime
		return true
	}
	return false
}
