package wal

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type WalHandler struct {
	walChan  chan string
	filePath string
}

func NewWalHandler(filePath string) *WalHandler {

	return &WalHandler{
		walChan:  make(chan string),
		filePath: filePath,
	}

}

func (walHandler *WalHandler) Run(ctx context.Context) {
	go walHandler.startWalTask(ctx)
}

func (walHandler *WalHandler) startWalTask(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		walHandler.walTask(ctx)
	}()
	wg.Wait()
}

func (walHandler *WalHandler) walTask(ctx context.Context) {
	for {
		select {
		case cmd := <-walHandler.walChan:
			println("Writing WAL")
			walEntry := strconv.FormatInt(time.Now().Unix(), 10) + " " + cmd + "\n"
			walHandler.WriteToFile(walEntry)
		case <-ctx.Done():
			return
		}
	}
}

func (walHandler *WalHandler) WriteToFile(cmd string) {
	f, err := os.OpenFile("wal.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Failed to write file: " + err.Error())
	}
	defer f.Close()
	if _, err = f.WriteString(cmd); err != nil {
		println("Error writing to file")
	}
}

func (walHandler *WalHandler) GetWriteChan() chan<- string {
	return walHandler.walChan
}
