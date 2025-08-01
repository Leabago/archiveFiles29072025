package app

import (
	"archiveFiles/config"
	"archiveFiles/internal/httpcontroller"
	httpclient "archiveFiles/internal/repo/httpClient"
	"archiveFiles/internal/usecase/links"
	"archiveFiles/pkg/httpserver"
	"archiveFiles/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) error {
	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		return err
	}

	// repo
	httpClient := httpclient.New(cfg.App.ContentType, cfg.App.MaxBytesResp)

	// usecase
	u := links.New(httpClient, cfg.App.MaxNumLinks, cfg.App.MaxTaskCount)
	// controller
	httpServer := httpserver.New(cfg.HTTP.Port)
	httpcontroller.NewRouter(cfg, httpServer.Server, u)

	// start server
	httpServer.Start(log)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		log.Info(fmt.Sprintf("app - Run - signal: %s", s.String()))
	case err = <-httpServer.Notify():
		log.Error(fmt.Sprintf("app - Run - httpServer.Notify: %s", err))
		return err
	}

	// shutdown
	err = httpServer.Shutdown(log)
	if err != nil {
		log.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err))
		return err
	}

	return nil
}
