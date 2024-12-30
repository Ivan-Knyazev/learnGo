package server

import (
	"errors"
	"go-storage/internal/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DictEntry struct {
	Value []string `json:"value" binding:"required"`
}

// Controllers for Dict
func (r *Server) dictGetValue(ctx *gin.Context) {
	key := ctx.Param("key")
	field := ctx.Param("field")

	value, err := r.Storage.GetDictField(key, field)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if value.ScalarValueType == storage.ScalarKindInt {
		ctx.JSON(http.StatusOK, ScalarEntry{Value: strconv.FormatInt(value.ScalarValueInt, 10), ValueType: value.ScalarValueType})
	} else if value.ScalarValueType == storage.ScalarKindString {
		ctx.JSON(http.StatusOK, ScalarEntry{Value: value.ScalarValueString, ValueType: value.ScalarValueType})
	} else {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (r *Server) dictGet(ctx *gin.Context) {
	key := ctx.Param("key")

	data, err := r.Storage.GetDict(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	dict := make(map[string]any)
	for key, value := range data {
		if value.ScalarValueType == storage.ScalarKindInt {
			dict[key] = value.ScalarValueInt
		} else if value.ScalarValueType == storage.ScalarKindString {
			dict[key] = value.ScalarValueString
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("bad value type in scalar"))
		}
	}
	ctx.JSON(http.StatusOK, dict)
}

func (r *Server) dictSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var payload DictEntry

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if count, err := r.Storage.SetDictFields(key, payload.Value...); err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	} else {
		ctx.String(http.StatusOK, "was added %d fields", count)
	}
}
