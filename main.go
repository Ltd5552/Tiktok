package main

import (
	"Tiktok/config"
	"Tiktok/controller"
	"Tiktok/model"
	"Tiktok/pkg/log"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func init() {
	config.InitViper()
	model.InitDB()
}

func main() {
	defer func() {
		err := log.Sync()
		if err != nil {
			log.Error("log sync err", zap.Error(err))
		}
	}()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.ServerSetting.Port),
		Handler:        controller.InitRouter(),
		ReadTimeout:    config.ServerSetting.ReadTimeout,
		WriteTimeout:   config.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Error("ListenAndServe err", zap.Error(err))
	}
}
