package routers

import (
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, headers map[string]string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	for headerKey, headerVal := range headers {
		req.Header.Add(headerKey, headerVal)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	type want struct {
		statusCode int
		headers    map[string]string
	}
	tests := []struct {
		urlToCall string
		name      string
		method    string
		headers   map[string]string
		want      want
	}{
		{
			name:      "should return 200 when add gauge",
			urlToCall: "/update/gauge/NewMetric/0.23",
			method:    http.MethodPost,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusOK,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},
		{
			name:      "should return 200 when add counter",
			urlToCall: "/update/counter/NewCounterMetric/5",
			method:    http.MethodPost,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusOK,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},
		{
			name:      "should return 400 when metric value NaN",
			urlToCall: "/update/counter/NewCounterMetric/stringHere",
			method:    http.MethodPost,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				headers:    map[string]string{},
			},
		},
		{
			name:      "should return 405 when caling not POST method",
			urlToCall: "/update/counter/NewCounterMetric/1",
			method:    http.MethodGet,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				headers:    map[string]string{},
			},
		},
		{
			name:      "should return 501 when unknown type supplied",
			urlToCall: "/update/unknown/NewCounterMetric/1",
			method:    http.MethodPost,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: http.StatusNotImplemented,
				headers:    map[string]string{},
			},
		},
		//{
		//	name:      "should return 400 when no Content-Type header supplied",
		//	urlToCall: "/update/counter/NewCounterMetric/1",
		//	method:    http.MethodPost,
		//	headers:   map[string]string{},
		//	want: want{
		//		statusCode: 400,
		//		headers: map[string]string{
		//			"Content-Type": "text/plain",
		//		},
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.ServerConfig{
				Address:       "localhost:8080",
				StoreInterval: 300 * time.Second,
				Restore:       false,
			}
			store := repository.NewMapRepo()
			r := NewRouter(&cfg, store)
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, _ := testRequest(t, ts, tt.method, tt.urlToCall, tt.headers)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			for k, v := range tt.want.headers {
				assert.Equal(t, v, resp.Header.Get(k))
			}

		})
	}
}
