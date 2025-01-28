package main

import (
	"log"
	"os"

	"github.com/jiraconnector/cmd/app"
	config "github.com/jiraconnector/internal/configReader"
)

func main() {
	// Open config file
	cfgPath := "../../configs/config.yml" // dev
	configFile, err := os.Open(cfgPath)
	if err != nil {
		log.Println("error open config")
		panic(err)
	}
	log.Println("open config")

	// Load date from config file
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Println("error load config")
		panic(err)
	}
	log.Println("load config")

	// Setting logger
	log.Printf("set logger")

	// Create connector app
	a, err := app.NewApp(cfg)
	if err != nil {
		log.Println("error create app")
		panic(err)
	}
	log.Println("created app")

	if err := a.Run(); err != nil {
		log.Println("error run app")
		panic(err)
	}
	defer a.Close()
}
