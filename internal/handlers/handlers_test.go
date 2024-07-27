// Модуль handlers содержит типы, методы и константы для
// API handlers
package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hrapovd1/msg-proc/internal/config"
	"github.com/hrapovd1/msg-proc/internal/metric"
	"github.com/hrapovd1/msg-proc/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	conf := config.Config{
		DatabaseDSN:  "",
		KafkaBrokers: "",
		KafkaTopic:   "message",
	}
	t.Run("simple check", func(t *testing.T) {
		hanlr := NewHandler(conf, log.Default())
		assert.IsType(t, &Handler{}, hanlr)
	})
}

func TestMetricsHandler_PingDB(t *testing.T) {
	handler := Handler{
		Storage: make(storgr, 4),
	}
	reqst := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	hndl := http.HandlerFunc(handler.PingDB)
	// qeury server
	hndl.ServeHTTP(rec, reqst)
	result := rec.Result()
	defer result.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestMetricsHandler_Metric(t *testing.T) {
	stor := make(storgr, 4)
	mtrcStor := metric.NewMetrics(stor)
	mtrcStor.Mem[types.InputMetric] = int64(5)
	mtrcStor.Mem[types.ProcessMetric] = int64(2)
	handler := Handler{
		Storage: stor,
		Metrics: *mtrcStor,
	}
	reqst := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	hndl := http.HandlerFunc(handler.Metric)
	// qeury server
	hndl.ServeHTTP(rec, reqst)
	result := rec.Result()
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "application/json", result.Header.Get("Content-Type"))
	assert.Equal(t, []byte(`{"total":5,"processed":2}`), body)
}

func TestMetricsHandler_SaveHandler(t *testing.T) {
	stor := make(storgr, 4)
	handler := Handler{
		Storage:    stor,
		MessageBus: stor,
		Metrics:    *metric.NewMetrics(stor),
		logger:     log.Default(),
	}
	t.Run("GET method", func(t *testing.T) {
		reqst := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		hndl := http.HandlerFunc(handler.SaveHandler)
		// qeury server
		hndl.ServeHTTP(rec, reqst)
		result := rec.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	})
	t.Run("POST method", func(t *testing.T) {
		msg := types.Message{Msg: "test message"}
		data, err := json.Marshal(msg)
		require.NoError(t, err)

		reqst := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(data)))
		reqst.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		hndl := http.HandlerFunc(handler.SaveHandler)
		// qeury server
		hndl.ServeHTTP(rec, reqst)
		result := rec.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, stor[2].(types.Message).Msg, msg.Msg)
	})
}

func TestNotImplementedHandler(t *testing.T) {
	reqst := httptest.NewRequest(http.MethodPost, "/update/any/", nil)
	rec := httptest.NewRecorder()
	hndl := http.HandlerFunc(NotImplementedHandler)
	hndl.ServeHTTP(rec, reqst)

	t.Run("Check not implemented", func(t *testing.T) {
		result := rec.Result()
		defer assert.Nil(t, result.Body.Close())
		_, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotImplemented, result.StatusCode)
	})
}

// Mock тип интерфейса Repository, BusMessenger и Storager
type storgr []any

func (s storgr) Save(ctx context.Context, msgID, message string) {
	s[2] = types.Message{ID: msgID, Msg: message}
}

func (s storgr) Update(ctx context.Context, msgID, status string) {}

func (s storgr) Close() error {
	return nil
}

func (s storgr) Ping(ctx context.Context) bool {
	return false
}

func (s storgr) GetCount(ctx context.Context, processed bool) int64 {
	if processed {
		return s[1].(int64)
	}
	return s[0].(int64)
}

func (s storgr) Write(ctx context.Context, msg types.Message) {
	s[3] = msg
}

func (s storgr) Read(ctx context.Context) types.Message {
	return s[3].(types.Message)
}
