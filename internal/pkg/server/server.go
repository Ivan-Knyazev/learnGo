package server

import (
	"go-storage/internal/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Host    string
	Storage storage.Storage
}

func NewServer(host string, storage storage.Storage) *Server {
	server := Server{
		Host:    host,
		Storage: storage,
	}
	return &server
}

func (r *Server) newAPI() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)  // On release mode
	engine := gin.Default()

	scalarRouter := engine.Group("/scalar")
	scalarRouter.GET("/get/:key", r.scalarGet)
	scalarRouter.PUT("/set/:key", r.scalarSet)

	dictRouter := engine.Group("/dict")
	dictRouter.GET("/get/value/:key/:field", r.dictGetValue)
	dictRouter.GET("/get/:key/", r.dictGet)
	dictRouter.PUT("/set/:key", r.dictSet)

	healthRouter := engine.Group("/health")
	healthRouter.GET("", r.health)

	return engine
}

func (r *Server) StartServer() error {
	return r.newAPI().Run(r.Host)
}

// Controller for health
func (r *Server) health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
