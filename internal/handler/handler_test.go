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
	"url-shortener/internal/domains/mocks"
)

type serviceMock func(c *mocks.UseCase)
type sessionMock func(c *mocks.SessionUseCase)

func TestHandler_UpdateAndGetShort(t *testing.T) {

	tests := []struct {
		name        string
		body        string
		sessionMock sessionMock
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK1",
			body: "https://ya.ru",
			sessionMock: func(c *mocks.SessionUseCase) {
				c.Mock.On("CreateIfNotExists").Return("ahsjufil12-fk", 12)
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", 12, "https://ya.ru").Return("adh35Kof", nil).Times(1)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "BAD",
			body: "12",
			sessionMock: func(c *mocks.SessionUseCase) {
				c.Mock.On("CreateIfNotExists").Return("ahsjufil12-fk", 12)
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", 12, "12").Return("", errors.New("invalid")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			service := mocks.NewUseCase(t)
			sessionMock := mocks.NewSessionUseCase(t)
			h := NewHandler(service, sessionMock, config.Config{})
			tt.serviceMock(service)
			tt.sessionMock(sessionMock)

			path := "/t"
			g.POST(path, h.UpdateAndGetShort)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(tt.body))

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
		sessionMock sessionMock
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name:     "OK1",
			shortURL: "fs23oyrh",
			sessionMock: func(c *mocks.SessionUseCase) {
				c.Mock.On("CreateIfNotExists").Return("ahsjufil12-fk", 12)
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetLong", 12, "fs23oyrh").Return("https://ya.ru", nil).Times(1)
			},
			wantCode: http.StatusTemporaryRedirect,
		},
		{
			name:     "BAD",
			shortURL: "aaa",
			sessionMock: func(c *mocks.SessionUseCase) {
				c.Mock.On("CreateIfNotExists").Return("ahsjufil12-fk", 12)
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetLong", 12, "aaa").Return("", errors.New("invalid id")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}
	g := gin.Default()

	service := mocks.NewUseCase(t)
	sessionMock := mocks.NewSessionUseCase(t)
	h := NewHandler(service, sessionMock, config.Config{})

	path := "/:id"
	g.GET(path, h.GetLongURL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.serviceMock(service)
			tt.sessionMock(sessionMock)

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
		sessionMock sessionMock
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK1",
			json: jso{
				URI: "https://ya.ru",
			},
			sessionMock: func(c *mocks.SessionUseCase) {
				c.Mock.On("CreateIfNotExists").Return("ahsjufil12-fk", 12)
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", 12, "https://ya.ru").Return("af3gyhj2", nil).Times(1)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "BAD",
			json: jso{
				URI: "as",
			},
			sessionMock: func(c *mocks.SessionUseCase) {
				c.Mock.On("CreateIfNotExists").Return("ahsjufil12-fk", 12)
			},
			serviceMock: func(c *mocks.UseCase) {
				c.Mock.On("GetShort", 12, "as").Return("", errors.New("invalid json")).Times(1)
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()

			service := mocks.NewUseCase(t)
			sessionMock := mocks.NewSessionUseCase(t)
			h := NewHandler(service, sessionMock, config.Config{})

			path := "/api/shorten"
			g.POST(path, h.GetShortByJSON)
			tt.serviceMock(service)
			tt.sessionMock(sessionMock)
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
