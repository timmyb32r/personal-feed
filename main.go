package main

import (
	"context"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"os"
	"personal-feed/pkg/config"
	_ "personal-feed/pkg/crawlers/registry"
	_ "personal-feed/pkg/repo/registry"
	"personal-feed/pkg/server"
	"time"
)

func main() {
	var logger = logrus.New()

	parser := argparse.NewParser("personal-feed", "daemon, collect useful info from the internet and to structure it")
	configPath := parser.String("c", "config", &argparse.Options{Required: true, Help: "path to config file"})
	isOnce := parser.Flag("o", "once", &argparse.Options{Required: false, Help: "if specified, exits after one cycle"})
	err := parser.Parse(os.Args)
	if err != nil {
		logger.Errorf("unable to parse arguments: %s" + err.Error())
		os.Exit(1)
	}

	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.Fatalf("unable to open file: %s", *configPath)
	}
	currConfig, err := config.Load(configFile)
	if err != nil {
		logger.Errorf("unable to load config: %s" + err.Error())
		os.Exit(1)
	}

	currServer := server.NewServer(currConfig, logger)
	defer currServer.Close()

	for {
		err := currServer.HandleAllSources(context.TODO())
		if err != nil {
			logger.Errorf("server returned error: %s" + err.Error())
		}
		if isOnce != nil && *isOnce {
			break
		}
		time.Sleep(time.Minute)
	}
}
