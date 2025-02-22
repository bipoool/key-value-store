package main

import (
	"context"
	"key-value-store/kvpApi"
	"key-value-store/store"
	"key-value-store/wal"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	numShards := uint8(4)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	walHandler := wal.NewWalHandler("")
	storeManager := store.NewStoreManger(numShards, walHandler.GetWriteChan())
	kvpController := kvpApi.NewKvpController(numShards, storeManager)

	walHandler.Run(ctx)
	storeManager.Run(ctx)

	router := gin.Default()
	router.GET("/get", kvpController.GetController)
	router.POST("/set", kvpController.SetController)
	router.DELETE("/delete", kvpController.DeleteController)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		stop()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
