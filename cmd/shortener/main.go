package main

import (
	"fmt"
	"log"
	"url-shortener/config"
	"url-shortener/internal/handler"
	myLog "url-shortener/internal/logger"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"

	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	conf := config.New()
	stM, err := storage.New(conf)
	if err != nil {
		log.Fatalf("Error creating storage: %s", err)
		return
	}
	sr, err := service.NewService(stM, conf)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
		return
	}

	ss, err := service.NewSessionService(stM)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
		return
	}
	h := handler.NewHandler(&sr, &ss, conf)
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	myLog.Sugar = *logger.Sugar()
	h.Start()
	stM.Close()
}
