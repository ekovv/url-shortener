package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Handler) AcceptEncoding() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncoding := c.GetHeader("accept-encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			c.Header("Content-Encoding", "gzip")
		}
	}
}

func (s *Handler) Decompressed() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			defer reader.Close()

			buf := new(strings.Builder)
			_, err = io.Copy(buf, reader)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Request.Body = io.NopCloser(strings.NewReader(buf.String()))
			c.Request.Header.Set("Content-Length", string(rune(len(buf.String()))))
		}
		c.Next()
	}
}
