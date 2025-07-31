package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimfu/finch/authentication/db"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := db.Connect(); err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()

	// main routes
	auth := router.Group("auth")
	auth.POST("/signup", Signup)
	auth.POST("/signin", Signin)
	auth.POST("/refresh", Refresh)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGABRT)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen :%s\n", err)
		}
	}()

	// waits for termination
	<-sigchan
	log.Println("Received termination signal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("Gracefully shutdown the server")
}
