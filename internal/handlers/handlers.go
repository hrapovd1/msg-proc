// Модуль handlers содержит типы, методы и константы для
// API handlers
package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/hrapovd1/msg-proc/internal/config"
	dbstorage "github.com/hrapovd1/msg-proc/internal/dbstrorage"
	"github.com/hrapovd1/msg-proc/internal/metric"
	"github.com/hrapovd1/msg-proc/internal/msgbus"
	"github.com/hrapovd1/msg-proc/internal/types"
	"github.com/hrapovd1/msg-proc/internal/usecase"
)

// Handler тип обработчиков API
// содержит конфигурацию и хранилище
type Handler struct {
	Storage    types.Repository
	MessageBus types.BusMessenger
	Metrics    metric.Metrics
	Config     config.Config
	logger     *log.Logger
}

// NewHandler возвращает обработчик API
func NewHandler(conf config.Config, logger *log.Logger) *Handler {
	h := &Handler{Config: conf, logger: logger}
	var stor types.Repository
	// db storage init
	stor, err := dbstorage.NewDBStorage(
		conf.DatabaseDSN,
		logger,
	)
	if err != nil {
		logger.Fatal(err)
	}
	h.Storage = stor
	// metrics storage init
	metrics := *metric.NewMetrics(stor.(types.Storager))
	h.Metrics = metrics
	// message bus init
	h.MessageBus = msgbus.NewKfkBus(conf, logger, &metrics)

	return h
}

// UpdateHandler POST обработчик сохранения сообщения в JSON формате
func (h *Handler) SaveHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Println(err)
		}
	}()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var data types.Message
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	usecase.AddMessageID(&data)

	// Write new message in DB
	usecase.WriteJSONMessage(
		ctx,
		data,
		h.Storage,
	)
	// Send message to Bus
	usecase.SendJSONMessage(
		ctx,
		data,
		h.MessageBus,
	)
	// Count input message
	h.Metrics.IncrementInput()

	rw.WriteHeader(http.StatusOK)
}

// PingDB GET обработчик проверки доступности базы
func (h *Handler) PingDB(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	dbstor := h.Storage.(types.Storager)
	if !dbstor.Ping(ctx) {
		http.Error(rw, "DB connect is NOT ok", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)
}

// Metrics GET обработчик API метрик
func (h *Handler) Metric(rw http.ResponseWriter, r *http.Request) {
	data := h.Metrics.GetMetrics()
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NotImplementedHandler обработчик для ответа на не реализованные url
func NotImplementedHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
	_, err := rw.Write([]byte("It's not implemented yet."))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotImplemented)
		return
	}
}
