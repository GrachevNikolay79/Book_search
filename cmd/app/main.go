package main

import (
	"book_search/internal/app"
	"book_search/internal/config"
	"flag"
	"log"
	"os"
)

func main() {
	cfg, err := configure()
	if err != nil {
		log.Println(err)
		return
	}
	app := app.NewApp(cfg)
	app.InitDatabase()
	app.Run()
	app.ShutDown()
}

func configure() (*config.Config, error) {
	configFileName := ""
	flag.StringVar(&configFileName, "c", "config.yaml", "Specify config filename")
	flag.Parse()

	_, err := os.Stat(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File dos not exist, create sample cfg file")
			config.SaveSampleConfig(configFileName)
			return nil, err
		} else {
			log.Println(err)
			return nil, err
		}
	}

	return config.GetConfig(configFileName), nil
}
