package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimfu/finch/authentication/controllers"
	"github.com/dimfu/finch/authentication/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBMiddleware(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {
	db, err := db.New()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// middlewares
	router.Use(DBMiddleware(db))

	// main routes
	auth := router.Group("auth")
	auth.POST("/signup", controllers.Signup)
	auth.POST("/signin", controllers.Signin)

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
