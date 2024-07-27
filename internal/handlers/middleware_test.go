package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gzipWriter_Write(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		Writer         io.Writer
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Positive",
			fields:  fields{Writer: &bytes.Buffer{}},
			args:    args{b: []byte(`Test string`)},
			want:    11,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := gzipWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				Writer:         tt.fields.Writer,
			}
			got, err := w.Write(tt.args.b)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMetricsHandler_GzipMiddle(t *testing.T) {
	handler := Handler{}
	simpleHandl := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Test is OK."))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	request := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Run("without compress", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.GzipMiddle(http.HandlerFunc(simpleHandl)).ServeHTTP(rec, request)
		result := rec.Result()
		defer assert.Nil(t, result.Body.Close())
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, 11, len(body))
	})
	t.Run("with compress", func(t *testing.T) {
		request.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()
		handler.GzipMiddle(http.HandlerFunc(simpleHandl)).ServeHTTP(rec, request)
		result := rec.Result()
		defer assert.Nil(t, result.Body.Close())
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, 39, len(body))
	})

}
