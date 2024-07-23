package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/hrapovd1/msg-proc/internal/config"
	"github.com/hrapovd1/msg-proc/internal/handlers"
	"github.com/hrapovd1/msg-proc/internal/types"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	logger := log.New(os.Stdout, "APP\t", log.Ldate|log.Ltime)
	// Чтение флагов и установка конфигурации сервера
	serverConf, err := config.NewAppConf(config.GetAppFlags())
	if err != nil {
		logger.Fatalln(err)
	}

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	logger.Printf("\tBuild version: %s\n", buildVersion)
	logger.Printf("\tBuild date: %s\n", buildDate)
	logger.Printf("\tBuild commit: %s\n", buildCommit)
	logger.Println("Server start on ", serverConf.ServerAddress)

	handlerMetrics := handlers.NewHandler(*serverConf, logger)
	handlerStorage := handlerMetrics.Storage.(types.Storager)
	defer func() {
		if err := handlerStorage.Close(); err != nil {
			logger.Print(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	router := chi.NewRouter()
	router.Use(handlerMetrics.GzipMiddle)
	router.Get("/ping", handlerMetrics.PingDB)
	router.Post("/value/", handlerMetrics.SaveHandler)
	router.Post("/*", handlers.NotImplementedHandler)

	server := http.Server{
		Addr:    serverConf.ServerAddress,
		Handler: router,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(c context.Context, w *sync.WaitGroup, s *http.Server, l *log.Logger) {
		defer wg.Done()
		<-c.Done()
		l.Println("got signal to stop")
		if err := s.Shutdown(c); err != nil {
			l.Println(err)
		}

	}(ctx, &wg, &server, logger)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal(err)
	}

	wg.Wait()
	logger.Println("app stopped gracefully")
}
