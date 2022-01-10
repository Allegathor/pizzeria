package api

import (
	"net/http"
	"pizzeria/repo"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type APIService struct {
	logger  *zap.Logger
	order   repo.Order
	storage repo.Storage
}

func NewAPIService(logger *zap.Logger, storage repo.Storage, order repo.Order) *APIService {
	return &APIService{
		logger:  logger,
		order:   order,
		storage: storage,
	}
}

func (s *APIService) router() *gin.Engine {
	g := gin.New()
	gin.SetMode(gin.ReleaseMode)

	g.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
	g.Use(ginzap.RecoveryWithZap(s.logger, true))

	g.GET("/storage", s.Storage)
	g.POST("/storage/add", s.Resupply)

	g.GET("/order", s.Orders)
	g.POST("/order/create", s.CreateOrder)

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
	s.logger.Info("Ingredient ressuplied", zap.Int("quantity", req.Qty), zap.Int("ID", req.ID))

	c.Status(http.StatusOK)
}

func (s *APIService) Orders(c *gin.Context) {
	list, err := s.order.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorJSON{Error: err.Error()})
		return
	}
	s.logger.Info("Orders", zap.Any("list", list))

	c.JSON(http.StatusOK, list)
}

func (s *APIService) CreateOrder(c *gin.Context) {
	var req repo.OrderInfo
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorJSON{Error: err.Error()})
		return
	}

	err := s.order.Create(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorJSON{Error: err.Error()})
		return
	}
	s.logger.Info("Order created", zap.Any("OrderInfo", req))

	c.Status(http.StatusOK)
}

type ErrorJSON struct {
	Error string `json:"error"`
}
