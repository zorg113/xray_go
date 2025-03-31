package internalhttp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/app"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage/memory"
)

type handler func(http.ResponseWriter, *http.Request)

type want struct {
	code        int
	handler     handler
	httpMetod   string
	url         string
	response    string
	contentType string
	bodyContent string
}

func Test_ServerHTTP(t *testing.T) {
	log, err := logger.New("INFO", "")
	require.NoError(t, err)
	stor := memorystorage.New()
	app := app.New(log, stor)
	srv := NewServer("localhost", "8080", *log, app)
	tests := []struct {
		name   string
		expect want
	}{
		{
			name: "hello world",
			expect: want{
				code:        http.StatusOK,
				handler:     srv.helloWorld,
				httpMetod:   "GET",
				url:         "/hello",
				response:    "Hello World!",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name: "create event",
			expect: want{
				code:        http.StatusOK,
				handler:     srv.createEvent,
				httpMetod:   "POST",
				url:         "/create",
				contentType: "application/json",
				bodyContent: `{"title": "Test event", "startDate": "2025-07-06T12:34:33.000000001Z", "endDate": "2025-07-06T13:34:33.000000001Z"}`, //nolint:lll
				response:    "Event created",
			},
		},
		{
			name: "update event",
			expect: want{
				code:        http.StatusOK,
				handler:     srv.updateEvent,
				httpMetod:   "POST",
				url:         "/update",
				contentType: "application/json",
				bodyContent: `{"title": "Test event", "startDate": "2025-07-06T12:34:33.000000001Z", "endDate": "2025-07-06T15:34:33.000000001Z"}`, //nolint:lll
				response:    "Event updated",
			},
		},
		{
			name: "get events",
			expect: want{
				code:        http.StatusOK,
				handler:     srv.getEvents,
				httpMetod:   "GET",
				url:         "/get",
				contentType: "application/json",
				bodyContent: `{"title": "Test event", "startDate": "2021-07-06T12:34:33.000000001Z", "endDate": "2026-07-06T12:34:33.000000001Z"}`,                                                                           //nolint:lll
				response:    "[{\"ID\":0,\"Title\":\"Test event\",\"startDate\":\"2025-07-06T12:34:33.000000001Z\",\"endDate\":\"2025-07-06T15:34:33.000000001Z\",\"Description\":\"\",\"OwnerID\":\"\",\"RemindIn\":\"\"}]", //nolint:lll
			},
		},
		{
			name: "delete event",
			expect: want{
				code:        http.StatusOK,
				handler:     srv.deleteEvent,
				httpMetod:   "POST",
				url:         "/get",
				contentType: "application/json",
				bodyContent: `{"title": "Test event", "startDate": "2025-07-06T12:34:33.000000001Z", "endDate": "2025-07-06T15:34:33.000000001Z"}`, //nolint:lll
				response:    "Event deleted",
			},
		},
	}
	rxctx := chi.NewRouteContext()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.expect.httpMetod, test.expect.url, nil)
			request.Header.Set("Content-type", test.expect.contentType)
			request.ContentLength = int64(len(test.expect.bodyContent))
			if request.ContentLength == 0 {
				request.Body = io.NopCloser(nil)
			} else {
				request.Body = io.NopCloser(strings.NewReader(test.expect.bodyContent))
			}
			w := httptest.NewRecorder()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rxctx))
			test.expect.handler(w, request)
			res := w.Result()
			assert.Equal(t, test.expect.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			if res.StatusCode != http.StatusTemporaryRedirect {
				assert.Equal(t, test.expect.response, string(resBody))
			}
		})
	}
}
