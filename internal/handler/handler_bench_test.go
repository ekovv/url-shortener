package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/config"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

func BenchmarkHandler_GetLongURL(b *testing.B) {
	s := 1000
	cfg := config.Config{}
	st, err := storage.New(cfg)
	if err != nil {
		return
	}

	us, err := service.NewService(st, cfg)
	if err != nil {
		return
	}
	ss, err := service.NewSessionService(st)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
		return
	}
	handler := NewHandler(&us, &ss, cfg)

	router := gin.Default()
	router.Use(handler.UpdateAndGetShort)
	req := httptest.NewRequest("POST", "/", nil)

	b.ResetTimer()
	b.Run("UpdateAndGetShort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("%d", s)))
			b.StartTimer()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			b.StopTimer()
			log.Println()
			s++
		}
	})
	uindex := 0
	router = gin.Default()
	router.GET("/:id", handler.GetLongURL)
	b.Run("GetLongURL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req = httptest.NewRequest("GET", "/Xz",
				nil)
			w := httptest.NewRecorder()

			b.StartTimer()
			router.ServeHTTP(w, req)
			b.StopTimer()
			uindex++
		}
	})
}
