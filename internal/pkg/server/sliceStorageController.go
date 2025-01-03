package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SliceEntry struct {
	Value []int `json:"value" binding:"required"`
}

type PopSettings struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type SlicePopEntry struct {
	Result int   `json:"result" binding:"required"`
	Value  []int `json:"value" binding:"required"`
}

type IntEntry struct {
	Value int `json:"value" binding:"required"`
}

// func NewPopSettings() *PopSettings {
// 	return &PopSettings{Start: 0, End: 0}
// }

// Controllers for Slice
func (r *Server) sliceGet(ctx *gin.Context) {
	key := ctx.Param("key")

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SliceEntry{Value: value})
}

func (r *Server) sliceGetValue(ctx *gin.Context) {
	key := ctx.Param("key")
	if index, err := strconv.Atoi(ctx.Param("index")); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	} else {
		value, err := r.Storage.GetSliceValue(key, index)
		if err != nil {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusOK, IntEntry{Value: value})
	}
}

func (r *Server) sliceSetValue(ctx *gin.Context) {
	key := ctx.Param("key")
	index, err := strconv.Atoi(ctx.Param("index"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var newElement IntEntry
	err = ctx.ShouldBindJSON(&newElement)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.Storage.SetSliceValue(key, index, newElement.Value); err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SliceEntry{Value: value})
}

func (r *Server) sliceLeftPush(ctx *gin.Context) {
	key := ctx.Param("key")

	var newElements SliceEntry
	err := ctx.ShouldBindJSON(&newElements)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.Storage.LeftPushIntoSlice(key, newElements.Value...); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SliceEntry{Value: value})
}

func (r *Server) sliceRightPush(ctx *gin.Context) {
	key := ctx.Param("key")

	var newElements SliceEntry
	err := ctx.ShouldBindJSON(&newElements)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.Storage.RightPushIntoSlice(key, newElements.Value...); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SliceEntry{Value: value})
}

func (r *Server) sliceRightUniquePush(ctx *gin.Context) {
	key := ctx.Param("key")

	var newElements SliceEntry
	err := ctx.ShouldBindJSON(&newElements)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.Storage.RightUniquePushIntoSlice(key, newElements.Value...); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SliceEntry{Value: value})
}

func (r *Server) sliceLeftPop(ctx *gin.Context) {
	key := ctx.Param("key")

	var inputPopSettings PopSettings
	err := ctx.ShouldBindJSON(&inputPopSettings)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.Storage.LeftPopFromSlice(key, []int{inputPopSettings.Start, inputPopSettings.End}...)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SlicePopEntry{Result: result, Value: value})
}

func (r *Server) sliceRightPop(ctx *gin.Context) {
	key := ctx.Param("key")

	var inputPopSettings PopSettings
	err := ctx.ShouldBindJSON(&inputPopSettings)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.Storage.RightPopFromSlice(key, []int{inputPopSettings.Start, inputPopSettings.End}...)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	value, err := r.Storage.GetSlice(key)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, SlicePopEntry{Result: result, Value: value})
}
