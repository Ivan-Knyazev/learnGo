package server

import (
	"go-storage/internal/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value string `json:"value" binding:"required"`
}

func NewServer(host string, storage *storage.Storage) *Server {
	server := Server{
		host:    host,
		storage: storage,
	}
	return &server
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.Default()

	scalarRouter := engine.Group("/scalar")
	scalarRouter.GET("/get/:key", r.scalarGet)
	scalarRouter.PUT("/set/:key", r.scalarSet)

	healthRouter := engine.Group("/health")
	healthRouter.GET("", r.health)

	return engine
}

func (r *Server) StartServer() error {
	return r.newAPI().Run(r.host)
}

func (r *Server) scalarGet(ctx *gin.Context) {
	key := ctx.Param("key")

	value, ok := r.storage.Get(key)
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: value})
}

func (r *Server) scalarSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var payload Entry

	// err := json.NewDecoder(ctx.Request.Body).Decode(&payload)
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Set(key, payload.Value)

	ctx.Status(http.StatusOK)
}

func (r *Server) health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
