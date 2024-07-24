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
	"github.com/hrapovd1/msg-proc/internal/msgbus"
	"github.com/hrapovd1/msg-proc/internal/types"
	"github.com/hrapovd1/msg-proc/internal/usecase"
)

// Handler тип обработчиков API
// содержит конфигурацию и хранилище
type Handler struct {
	Storage    types.Repository
	MessageBus types.BusMessenger
	Config     config.Config
	logger     *log.Logger
}

// NewHandler возвращает обработчик API
func NewHandler(conf config.Config, logger *log.Logger) *Handler {
	h := &Handler{Config: conf, logger: logger}
	// db storage init
	db, err := dbstorage.NewDBStorage(
		conf.DatabaseDSN,
		logger,
	)
	if err != nil {
		logger.Fatal(err)
	}
	h.Storage = db
	// message bus init
	h.MessageBus = msgbus.NewKfkBus(conf, logger)

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

	// Write new metrics value
	err = usecase.WriteJSONMessage(
		ctx,
		data,
		h.Storage,
	)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

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

// NotImplementedHandler обработчик для ответа на не реализованные url
func NotImplementedHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
	_, err := rw.Write([]byte("It's not implemented yet."))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotImplemented)
		return
	}
}
