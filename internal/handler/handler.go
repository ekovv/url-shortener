package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"url-shortener/config"
	"url-shortener/internal/domains"
	myLog "url-shortener/internal/logger"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
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
	var user string
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.service.SaveAndGetSessionMap(token)
	} else {
		newToken := s.SetSession(c)
		id = s.service.SaveAndGetSessionMap(newToken)
	}
	user = strconv.Itoa(id)
	str := string(body)
	short, err := s.service.GetShort(user, str)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			c.String(http.StatusConflict, short)
			return
		}
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.String(http.StatusCreated, short)
}

func (s *Handler) GetLongURL(c *gin.Context) {
	idOfParam := c.Param("id")
	var user string
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.service.SaveAndGetSessionMap(token)
	} else {
		newToken := s.SetSession(c)
		id = s.service.SaveAndGetSessionMap(newToken)
	}
	user = strconv.Itoa(id)
	long, err := s.service.GetLong(user, idOfParam)
	if err != nil {
		fmt.Println("Error getting")
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
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	var user string
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.service.SaveAndGetSessionMap(token)
	} else {
		newToken := s.SetSession(c)
		id = s.service.SaveAndGetSessionMap(newToken)
	}
	user = strconv.Itoa(id)
	short, err := s.service.GetShort(user, js.URI)
	fmt.Println(short)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			js.Res = short
			js.URI = ""
			bytes, err := json.MarshalIndent(js, "", "    ")
			if err != nil {
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			c.Status(http.StatusConflict)
			c.Header("Content-Type", "application/json")
			c.Writer.Write(bytes)
			return
		}
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	js.Res = short
	js.URI = ""
	bytes, err := json.MarshalIndent(js, "", "    ")
	if err != nil {
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
	var input []jBatch
	var res []jBatch
	err := c.ShouldBindJSON(&input)
	if err != nil {
		_ = fmt.Errorf("error opening file storage %w", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	var user string
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.service.SaveAndGetSessionMap(token)
	} else {
		newToken := s.SetSession(c)
		id = s.service.SaveAndGetSessionMap(newToken)
	}
	user = strconv.Itoa(id)
	for _, i := range input {
		short, err := s.service.SaveWithoutGenerate(user, i.ID, i.Origin)
		if err != nil && errors.Is(err, storage.ErrAlreadyExists) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		i.Short = short
		i.Origin = ""
		res = append(res, i)

	}
	//res = append(res, string(bytes))
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, res)

}

func (s *Handler) GetAll(c *gin.Context) {
	var user string
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.service.SaveAndGetSessionMap(token)
	} else {
		c.Status(http.StatusUnauthorized)
		return
	}
	user = strconv.Itoa(id)
	urlsFrom, err := s.service.GetAllUrls(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error")
		return
	}
	if len(urlsFrom) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	var res []jBatch
	for _, i := range urlsFrom {
		batch := jBatch{}
		batch.Origin = i.Original
		batch.Short = i.Short
		res = append(res, batch)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
}

func (s *Handler) SetSession(c *gin.Context) string {
	uuid := service.GenerateUUID()
	c.SetCookie("token", uuid, 3600, "", "localhost", false, true)
	return uuid
}
