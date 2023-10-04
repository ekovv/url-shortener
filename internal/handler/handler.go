package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener/config"
	"url-shortener/internal/domains"
	myLog "url-shortener/internal/logger"
)

type Handler struct {
	service domains.UseCase
	engine  *gin.Engine
	config  config.Config
}

func NewHandler(service domains.UseCase, conf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		config:  conf,
		engine:  router,
	}
	router.Use(h.AcceptEncoding())
	router.Use(h.Decompressed())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(myLog.HTTPLogger())
	Route(router, h)
	return h
}

func (s *Handler) Start() {
	s.engine.Run(s.config.Host)
}

func (s *Handler) UpdateAndGetShort(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return

	}
	str := string(body)
	short, err := s.service.GetShort(str)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.String(http.StatusCreated, short)
}

func (s *Handler) GetLongURL(c *gin.Context) {
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
	var js uriJSON
	err := c.ShouldBindJSON(&js)
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

func (s *Handler) GetConnection(c *gin.Context) {
	err := s.service.CheckConn()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}

func (s *Handler) GetBatch(c *gin.Context) {
	var jB []jBatch
	var res []jBatch
	err := c.ShouldBindJSON(&jB)
	if err != nil {
		fmt.Println("JSON NOT GOOD")
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	for _, i := range jB {
		short, err := s.service.SaveLog(i.ID, i.Origin)
		if err != nil {
			i.Short = short
			i.Origin = ""
			res = append(res, i)
			c.JSON(http.StatusConflict, res)
			return
		}
		i.Short = short
		i.Origin = ""
		res = append(res, i)
	}
	//short, err := s.service.SaveLog(batch.ID, batch.Origin)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	//	return
	//}
	//bytes, err := json.MarshalIndent(i, "", "    ")
	//if err != nil {
	//	fmt.Println("JSON NOT GOOD")
	//	c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	//	return
	//}
	//res = append(res, string(bytes))
	c.Status(http.StatusCreated)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, res)
	//c.Writer.Write(bytes)
}