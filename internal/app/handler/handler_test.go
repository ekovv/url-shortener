package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/config"
	"url-shortener/internal/app/domains/mocks"
)

type serviceMock func(c *mocks.UseCase)

func TestHandler_UpdateAndGetShort(t *testing.T) {

	tests := []struct {
		name        string
		body        string
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK1",
			body: "https://ya.ru",
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", "https://ya.ru").Return("adh35Kof", nil).Times(1)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "BAD",
			body: "12",
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", "12").Return("", errors.New("invalid")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewUseCase(t)
			h := NewHandler(service, config.Config{})
			tt.serviceMock(service)

			path := "/t"
			g.POST(path, h.UpdateAndGetShort)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(tt.body))
			// создаём новый Recorder

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_GetLongURL(t *testing.T) {
	tests := []struct {
		name        string
		shortURL    string
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name:     "OK1",
			shortURL: "fs23oyrh",
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetLong", "fs23oyrh").Return("https://ya.ru", nil).Times(1)
			},
			wantCode: http.StatusTemporaryRedirect,
		},
		{
			name:     "BAD",
			shortURL: "aaa",
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetLong", "aaa").Return("", errors.New("invalid id")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}
	g := gin.Default()

	service := mocks.NewUseCase(t)
	h := NewHandler(service, config.Config{})

	path := "/:id"
	g.GET(path, h.GetLongURL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.serviceMock(service)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
			// создаём новый Recorder

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandler_GetShortByJSON(t *testing.T) {
	type jso struct {
		URI string `json:"url"`
	}
	tests := []struct {
		name        string
		json        jso
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK1",
			json: jso{
				URI: "https://ya.ru",
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", "https://ya.ru").Return("af3gyhj2", nil).Times(1)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "BAD",
			json: jso{
				URI: "as",
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", "as").Return("", errors.New("invalid json")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()

			service := mocks.NewUseCase(t)
			h := NewHandler(service, config.Config{})

			path := "/api/shorten"
			g.POST(path, h.GetShortByJSON)
			tt.serviceMock(service)
			b, err := json.Marshal(tt.json)
			if err != nil {
				return
			}
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(string(b)))

			g.ServeHTTP(w, request)

			if w.Code != tt.wantCode {
				t.Errorf("got %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}
