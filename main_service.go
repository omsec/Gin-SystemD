//go:build service

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func main() {
	fmt.Println("this is the service build.")

	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	registerControllers()

	srv := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		listeners, err := activation.Listeners()
		if err != nil {
			// log.Panicf("cannot retrieve listeners: %s", err)
			log.Fatalf("cannot retrieve listeners: %s", err)
		}
		if len(listeners) != 1 {
			log.Fatalf("unexpected number of socket activation (%d != 1)",
				len(listeners))
		}

		daemon.SdNotify(false, daemon.SdNotifyReady) // readiness

		// liveness (check health)
		go func() {
			interval, err := daemon.SdWatchdogEnabled(false)
			if err != nil || interval == 0 {
				return
			}
			for {
				_, err := http.Get("http://localhost:3000") // ❸
				if err == nil {
					daemon.SdNotify(false, daemon.SdNotifyWatchdog)
				}
				time.Sleep(interval / 3)
			}
		}()

		if err := srv.Serve(listeners[0]); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}

// go build (normal)
// go build -tags service (systemD)
