package main

import (
	"context"
	"flag"
	"quic_upload/api/pkg/logging"

	"quic_upload/api/internal/application"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized...")
	configFile := flag.String("config", "configs/config.yml", "Path to config file.")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := application.NewApplication(ctx, *configFile, &logger)
	if err != nil {
		logger.Fatal(err)
	}

	if err := app.StartHttp3Server(); err != nil {
		logger.Fatalln(err)
	}

}
