package main

import (
	"log"
	"os"

	config "github.com/jiraconnector/internal/configReader"
)

func main() {
	// Open config file
	cfgPath := "../../configs/config.yml" // dev
	configFile, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}

	// Load date from config file
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	// Setting logger
	log.Printf("jira url: %s", cfg.JiraCfg.Url)

}
