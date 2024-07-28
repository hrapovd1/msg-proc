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
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	logger := log.New(os.Stdout, "APP\t", log.Ldate|log.Ltime)
	// Чтение переменных среды и установка конфигурации приложения
	appConf, err := config.NewAppConf()
	if err != nil {
		logger.Fatalf("Error when init config: %v\n", err)
	}

	if BuildVersion == "" {
		BuildVersion = "N/A"
	}
	if BuildDate == "" {
		BuildDate = "N/A"
	}
	if BuildCommit == "" {
		BuildCommit = "N/A"
	}

	logger.Printf("\tBuild version: %s\n", BuildVersion)
	logger.Printf("\tBuild date: %s\n", BuildDate)
	logger.Printf("\tBuild commit: %s\n", BuildCommit)
	logger.Println("Server start on ", appConf.ServerAddress)

	handlerMessages := handlers.NewHandler(*appConf, logger)
	handlerStorage := handlerMessages.Storage.(types.Storager)
	defer func() {
		if err := handlerStorage.Close(); err != nil {
			logger.Print(err)
		}
	}()

	handlerMetrics := handlerMessages.Metrics

	handlerBusConsumer := handlerMessages.MessageBus.(types.MsgConsumer)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	// Запускаю потребителя сообщений в фоне
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer logger.Println("consumer are stopped")
		handlerBusConsumer.Consume(ctx, handlerMessages.Storage)
	}()

	// Запускаю синхронизацию метрик из базы
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer logger.Println("DB sync are stopped")
		handlerMetrics.SyncWithDB(ctx)
	}()

	router := chi.NewRouter()
	router.Use(handlerMessages.GzipMiddle)
	router.Use(handlers.Authenticator)
	router.Get("/ping", handlerMessages.PingDB)
	router.Get("/metrics", handlerMessages.Metric)
	router.Post("/api", handlerMessages.SaveHandler)
	router.Post("/*", handlers.NotImplementedHandler)

	server := http.Server{
		Addr:    appConf.ServerAddress,
		Handler: router,
	}

	wg.Add(1)
	go func(c context.Context, w *sync.WaitGroup, s *http.Server, l *log.Logger) {
		defer w.Done()
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
