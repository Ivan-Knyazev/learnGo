package server

import (
	"go-storage/internal/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScalarEntry struct {
	Value     string             `json:"value" binding:"required"`
	ValueType storage.ScalarKind `json:"valueType"`
	TTL       int64              `json:"ttl" binding:"required"`
}

type ScalarResponse struct {
	Value     string             `json:"value" binding:"required"`
	ValueType storage.ScalarKind `json:"valueType" binding:"required"`
	ExpiresAt int64              `json:"expiresAt" binding:"required"`
}

// Controllers for Scalar
func (r *Server) scalarGet(ctx *gin.Context) {
	key := ctx.Param("key")

	value, ok, expiresAt := r.Storage.GetScalar(key)
	kind := r.Storage.GetScalarKind(key)
	if !ok || kind == storage.ScalarKindUndefined {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, ScalarResponse{Value: value, ValueType: kind, ExpiresAt: expiresAt})
}

func (r *Server) scalarSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var payload ScalarEntry

	// err := json.NewDecoder(ctx.Request.Body).Decode(&payload)
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.Storage.SetScalar(payload.TTL, key, payload.Value)

	ctx.Status(http.StatusOK)
}
