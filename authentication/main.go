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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := db.Connect(); err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()
	if os.Getenv("ENV_MODE") == "development" {
		config := cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
		router.Use(cors.New(config))
	}

	auth := router.Group("/api/auth")
	auth.POST("/signup", SignUp)
	auth.POST("/signin", SignIn)
	auth.POST("/refresh", Refresh)
	auth.GET("/signout", SignOut)

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
