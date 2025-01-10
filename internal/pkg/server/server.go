package server

import (
	"fmt"
	"go-storage/internal/pkg/storage"
	"log"
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
	engine := gin.Default()

	scalarRouter := engine.Group("/scalar")
	scalarRouter.GET("/get/:key", r.scalarGet)
	scalarRouter.PUT("/set/:key", r.scalarSet)

	dictRouter := engine.Group("/dict")
	dictRouter.GET("/get/value/:key/:field", r.dictGetValue)
	dictRouter.GET("/get/:key/", r.dictGet)
	dictRouter.PUT("/set/:key", r.dictSet)

	sliceRouter := engine.Group("/slice")
	sliceRouter.GET("/get/:key/", r.sliceGet)
	sliceRouter.GET("/get/value/:key/:index", r.sliceGetValue)
	sliceRouter.PUT("/set/value/:key/:index", r.sliceSetValue)
	sliceRouter.POST("/push/left/:key", r.sliceLeftPush)
	sliceRouter.POST("/push/right/:key", r.sliceRightPush)
	sliceRouter.POST("/push/right/unique/:key", r.sliceRightUniquePush)
	sliceRouter.DELETE("/pop/left/:key", r.sliceLeftPop)
	sliceRouter.DELETE("/pop/right/:key", r.sliceRightPop)

	healthRouter := engine.Group("/health")
	healthRouter.GET("", r.health)

	return engine
}

func (r *Server) StartServer() *http.Server {
	server := &http.Server{Addr: r.Host, Handler: r.newAPI()}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	r.Storage.WriteLog(fmt.Sprintf("server was started on %s", r.Host))

	return server
	// return r.newAPI().Run(r.Host)
}

// Controller for health
func (r *Server) health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
