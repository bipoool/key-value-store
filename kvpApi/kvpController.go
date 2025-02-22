package kvpApi

import (
	"key-value-store/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

type kvpController struct {
	kvpHelper *kvpHelper
}

func NewKvpController(numShards uint8, storeManager *store.StoreManager) *kvpController {
	return &kvpController{
		kvpHelper: NewKvp(numShards, storeManager),
	}
}

func (k *kvpController) GetController(ctx *gin.Context) {
	key := ctx.Query("key")
	value, err := k.kvpHelper.Get(key)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if value == nil {
		ctx.JSON(404, gin.H{"error": "Not Defined"})
		return
	}
	ctx.JSON(200, gin.H{
		"key":   key,
		"value": value.GetValue().(string),
		"ttl":   value.GetTtl(),
	})
}

func (k *kvpController) SetController(ctx *gin.Context) {
	key := ctx.Query("key")
	value := ctx.Query("value")
	ttl, err := strconv.ParseUint(ctx.Query("ttl"), 10, 64)

	if err != nil {
		println("Wrong TTL Value")
		ttl = store.DEFAULT_TTL
	}
	errFromKvp := k.kvpHelper.Set(key, value, ttl)
	if errFromKvp != nil {
		ctx.JSON(500, gin.H{"error": errFromKvp})
		return
	}
	ctx.JSON(200, gin.H{"error": "None"})
}

func (k *kvpController) DeleteController(ctx *gin.Context) {
	key := ctx.Query("key")
	err := k.kvpHelper.Delete(key)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"error": "None"})
}
