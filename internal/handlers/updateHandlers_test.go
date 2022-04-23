package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleMetricPositive(t *testing.T) {
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
				statusCode: 200,
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
				statusCode: 200,
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
			},
		},
		{
			name:      "should return 400 when calling wrong url within /update",
			urlToCall: "/update/wrong/NewMetric/0.23",
			method:    http.MethodPost,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: 400,
				headers:    map[string]string{},
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
				statusCode: 400,
				headers:    map[string]string{},
			},
		},
		{
			name:      "should return 400 when caling not POST method",
			urlToCall: "/update/counter/NewCounterMetric/1",
			method:    http.MethodGet,
			headers: map[string]string{
				"Content-Type": "text/plain",
			},
			want: want{
				statusCode: 400,
				headers:    map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.urlToCall, nil)
			for headerKey, headerVal := range tt.headers {
				request.Header.Add(headerKey, headerVal)
			}
			responseRecorder := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(HandleMetric)

			handlerFunc.ServeHTTP(responseRecorder, request)
			result := responseRecorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			for k, v := range tt.want.headers {
				assert.Equal(t, v, result.Header.Get(k))
			}

		})
	}
}
