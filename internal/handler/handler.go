package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"url-shortener/config"
	"url-shortener/internal/domains"
	myLog "url-shortener/internal/logger"
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
	token, err := c.Request.Cookie("token")
	if err == nil {
		user = token.Value
	} else {
		user = s.SetTokenAndGetIfExist(c)
	}
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
	id := c.Param("id")
	var user string
	token, err := c.Request.Cookie("token")
	if err == nil {
		user = token.Value
	} else {
		user = s.SetTokenAndGetIfExist(c)
	}
	long, err := s.service.GetLong(user, id)
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
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	var user string
	token, err := c.Request.Cookie("token")
	if err == nil {
		user = token.Value
	} else {
		user = s.SetTokenAndGetIfExist(c)
	}
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
	token, err := c.Request.Cookie("token")
	if err == nil {
		user = token.Value
	} else {
		user = s.SetTokenAndGetIfExist(c)
		c.Status(http.StatusUnauthorized)
	}
	for _, i := range input {
		short, err := s.service.SaveWithoutGenerate(user, i.ID, i.Origin)
		if err != nil && errors.Is(err, storage.ErrAlreadyExists) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
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
	token, err := c.Request.Cookie("token")
	if err == nil {
		user = token.Value
	} else {
		user = s.SetTokenAndGetIfExist(c)
	}
	urlsFrom, err := s.service.GetAllUrls(user)
	if err != nil {
		c.Status(http.StatusNoContent)
		fmt.Errorf("error getting")
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
	c.JSON(http.StatusCreated, res)
}

func (s *Handler) SetTokenAndGetIfExist(c *gin.Context) string {
	newToken := uuid.New().String()
	cookie := &http.Cookie{
		Name:  "token",
		Value: newToken,
		Path:  "/",
	}
	http.SetCookie(c.Writer, cookie)
	return cookie.Value
}
