package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/app/domains/mocks"
)

type serviceMock func(c *mocks.Usecase)

func TestHandler_UpdateAndRetShort(t *testing.T) {

	tests := []struct {
		name        string
		body        string
		serviceMock serviceMock
		wantCode    int
	}{
		{
			name: "OK1",
			body: "https://ya.ru",
			serviceMock: func(c *mocks.Usecase) {
				c.Mock.On("RetShort", "https://ya.ru").Return("adh35Kof", nil).Times(1)
			},
			wantCode: http.StatusCreated,
		},
	}

	g := gin.Default()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewUsecase(t)
			h := NewHandler(service)
			tt.serviceMock(service)

			path := "/t"
			g.POST(path, h.UpdateAndRetShort)

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
			serviceMock: func(c *mocks.Usecase) {
				c.Mock.On("RetLong", "fs23oyrh").Return("https://ya.ru", nil).Times(1)
			},
			wantCode: http.StatusTemporaryRedirect,
		},
	}

	g := gin.Default()

	service := mocks.NewUsecase(t)
	h := NewHandler(service)

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
