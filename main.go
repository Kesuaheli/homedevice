package main

import (
	"context"
	"fmt"
	"homedevice/webserver"
	logger "log"
	"os/signal"
	"syscall"
)

var log logger.Logger = *logger.New(logger.Writer(), "[MAIN] ", logger.LstdFlags|logger.Lmsgprefix)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	webserver.Run()

	<-ctx.Done()
	fmt.Println()
	fmt.Println()
	log.Print("Shutting down...")
}
