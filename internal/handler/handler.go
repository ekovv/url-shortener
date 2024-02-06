package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"url-shortener/config"
	"url-shortener/internal/domains"
	myLog "url-shortener/internal/logger"
	"url-shortener/internal/storage"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service        domains.UseCase
	sessionService domains.SessionService
	engine         *gin.Engine
	config         config.Config
}

// NewHandler constructor with initialization gin
func NewHandler(service domains.UseCase, sessionService domains.SessionService, conf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service:        service,
		sessionService: sessionService,
		config:         conf,
		engine:         router,
	}
	router.Use(h.AcceptEncoding())
	router.Use(h.Decompressed())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(myLog.HTTPLogger())
	Route(router, h)
	return h
}

// Start server
func (s *Handler) Start() {
	s.engine.Run(s.config.Host)
}

// UpdateAndGetShort Compress links
func (s *Handler) UpdateAndGetShort(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var id int
	var session string
	token, err := c.Cookie("token")
	if err != nil {
		session, id = s.sessionService.CreateIfNotExists()
		s.SetSession(c, session)
	} else {
		id = s.sessionService.GetID(token)
	}
	str := string(body)
	short, err := s.service.GetShort(id, str)
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

// GetLongURL Get long links
func (s *Handler) GetLongURL(c *gin.Context) {
	idOfParam := c.Param("id")
	var id int
	var session string
	token, err := c.Cookie("token")
	if err == nil {
		id = s.sessionService.GetID(token)
	} else {
		session, id = s.sessionService.CreateIfNotExists()
		s.SetSession(c, session)
	}
	long, err := s.service.GetLong(id, idOfParam)
	if long == "" && err == nil {
		c.Status(http.StatusGone)
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", long)
}

// GetShortByJSON Get short links with json
func (s *Handler) GetShortByJSON(c *gin.Context) {
	var js uriJSON
	err := c.ShouldBindJSON(&js)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var id int
	var session string
	token, err := c.Cookie("token")
	if err == nil {
		id = s.sessionService.GetID(token)
	} else {
		session, id = s.sessionService.CreateIfNotExists()
		s.SetSession(c, session)
	}
	short, err := s.service.GetShort(id, js.URI)
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

// GetConnection Check connection
func (s *Handler) GetConnection(c *gin.Context) {
	err := s.service.CheckConn()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}

// GetBatch Get json with long and short links
func (s *Handler) GetBatch(c *gin.Context) {
	var input []batch
	var res []batch
	err := c.ShouldBindJSON(&input)
	if err != nil {
		_ = fmt.Errorf("error opening file storage %w", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	var id int
	var session string
	token, err := c.Cookie("token")
	if err == nil {
		id = s.sessionService.GetID(token)
	} else {
		session, id = s.sessionService.CreateIfNotExists()
		s.SetSession(c, session)
	}
	for _, i := range input {
		short, err := s.service.SaveWithoutGenerate(id, i.ID, i.Origin)
		if err != nil && errors.Is(err, storage.ErrAlreadyExists) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		i.Short = short
		i.Origin = ""
		res = append(res, i)

	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, res)

}

// GetAll Get all links
func (s *Handler) GetAll(c *gin.Context) {
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.sessionService.GetID(token)
	} else {
		c.Status(http.StatusUnauthorized)
		return
	}
	urlsFrom, err := s.service.GetAllUrls(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error")
		return
	}
	if len(urlsFrom) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	var res []batch
	for _, i := range urlsFrom {
		batch := batch{}
		batch.Origin = i.Original
		batch.Short = s.config.BaseURL + i.Short
		res = append(res, batch)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
}

// SetSession Set session for user
func (s *Handler) SetSession(c *gin.Context, session string) {
	c.SetCookie("token", session, 3600, "", "localhost", false, true)
}

// Del Delete links
func (s *Handler) Del(c *gin.Context) {
	var id int
	token, err := c.Cookie("token")
	if err == nil {
		id = s.sessionService.GetID(token)
	} else {
		c.Status(http.StatusNoContent)
		return
	}
	var inputList []string
	err = c.ShouldBindJSON(&inputList)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	go func() {
		err = s.service.Delete(inputList, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}()
	c.Status(http.StatusAccepted)
}
