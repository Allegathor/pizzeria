package api

import (
	"pizzeria/repo"

	"net/http"

	"github.com/gin-gonic/gin"
)

type APIService struct {
	storage repo.Storage
}

func NewAPIService(storage repo.Storage) *APIService {
	return &APIService{
		storage: storage,
	}
}

func (s *APIService) router() *gin.Engine {
	g := gin.New()
	g.GET("/storage", s.Storage)
	g.POST("/storage/add", s.Resupply)

	return g
}

func (s *APIService) Run() {
	g := s.router()

	g.Run(":8000")
}

func (s *APIService) Storage(c *gin.Context) {
	list, err := s.storage.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorJSON{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

type ResupplyRequest struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

func (s *APIService) Resupply(c *gin.Context) {
	var req ResupplyRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorJSON{Error: err.Error()})
		return
	}

	err := s.storage.Add(c, req.Qty, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorJSON{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type ErrorJSON struct {
	Error string `json:"error"`
}
