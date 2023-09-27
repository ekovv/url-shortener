package handler

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"url-shortener/internal/app/domains"
)

type Handler struct {
	service domains.UseCase
}

func NewHandler(service domains.UseCase) *Handler {
	return &Handler{service: service}
}

func (s *Handler) UpdateAndGetShort(c *gin.Context) {
	acceptEncoding := c.GetHeader("accept-encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		c.Header("Content-Encoding", "gzip")
	}
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return

	}
	str := string(body)
	fmt.Println(str)
	short, err := s.service.GetShort(str)
	fmt.Println(short)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.String(http.StatusCreated, short)
}

func (s *Handler) GetLongURL(c *gin.Context) {
	acceptEncoding := c.GetHeader("accept-encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		c.Header("Content-Encoding", "gzip")
	}
	id := c.Param("id")
	long, err := s.service.GetLong(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", long)

}

func (s *Handler) GetShortByJSON(c *gin.Context) {
	acceptEncoding := c.GetHeader("accept-encoding")
	if strings.Contains(acceptEncoding, "gzip") {
		c.Header("Content-Encoding", "gzip")
	}
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("JSON NOT GOOD")
		c.Status(http.StatusBadRequest)
		return
	}
	type uriJSON struct {
		URI string `json:"url,omitempty"`
		Res string `json:"result"`
	}

	var js uriJSON
	err = json.Unmarshal(b, &js)
	if err != nil {
		fmt.Println("JSON NOT GOOD")
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	short, err := s.service.GetShort(js.URI)
	fmt.Println(short)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	js.Res = short
	js.URI = ""
	bytes, err := json.MarshalIndent(js, "", "    ")
	if err != nil {
		fmt.Println("JSON NOT GOOD")
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
	c.Header("Content-Type", "application/json")
	c.Writer.Write(bytes)

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
