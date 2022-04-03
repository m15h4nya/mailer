package main

import (
	"apitask/config"
	db2 "apitask/db"
	"apitask/logger"
	srv "apitask/server"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func setupSignalHandler() chan os.Signal {

	sigChan := make(chan os.Signal, 4)
	signal.Notify(sigChan, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	return sigChan
}

func main() {
	log := logger.NewLogger()

	cfg, err := config.Load("./config/config.toml")
	if err != nil {
		log.Fatal("Failed to load config file: ", err)
	}

	dbConn := db2.NewConnection(cfg)
	err = dbConn.AutoMigrate(&db2.Mailing{}, &db2.Msg{}, &db2.Client{})
	if err != nil {
		log.Fatal("Failed to migrate tables: ", err)
	}

	ch := make(chan struct{})
	server := srv.NewServer(log, dbConn, cfg, ch)
	go srv.StartMailing(ch, dbConn)
	exit := setupSignalHandler()
	go func() {
		<-exit
		server.Shutdown(context.Background())
	}()

	log.Info("Listening :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Warn(err)
	}
}
