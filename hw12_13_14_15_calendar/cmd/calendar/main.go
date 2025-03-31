package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/app"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/config"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}
	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("can't start logger: %v", err)
	}

	var storage app.Storage

	if config.Storage.Type == "DB" {
		storage, err = sqlstorage.New(context.Background(), config.Storage.Database)
		if err != nil {
			log.Fatalf("can't start logger: %v", err)
		}
		logg.Info("set database storage")
	}

	if config.Storage.Type == "memory" {
		storage = memorystorage.New()
		logg.Info("set memory storage")
	}

	calendar := app.New(logg, storage)
	grpcServer, err := grpc.NewServer(logg, calendar, config.GRPCserver.Host, config.GRPCserver.Port)
	if err != nil {
		logg.Error(" failed to create grpc server: " + err.Error())
		return
	}
	defer grpcServer.Stop()

	server := internalhttp.NewServer(
		config.HTTPserver.Host,
		config.HTTPserver.Port,
		*logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)
		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcServer.Stop(); err != nil {
			logg.Error("failed to gracefully stop grpc server: " + err.Error())
		}

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()
	wg.Wait()
}
