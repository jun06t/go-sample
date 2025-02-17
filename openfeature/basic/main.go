package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

const (
	defaultMessage    = "Hello!"
	newWelcomeMessage = "Hello, welcome to this OpenFeature-enabled website!"
)

func main() {
	// シグナルを受け取るcontextを作成
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := openfeature.SetNamedProviderAndWait("my-domain", flagd.NewProvider())
	if err != nil {
		log.Fatalf("Failed to set the OpenFeature provider: %v", err)
	}
	// Iinitialize the OpenFeature client with domain name
	client := openfeature.NewClient("my-domain")

	engine := gin.Default()

	engine.GET("/hello", func(c *gin.Context) {
		// Evaluate welcome-message feature flag
		welcomeMessage, _ := client.BooleanValue(
			context.Background(), "welcome-message", false, openfeature.EvaluationContext{},
		)

		if welcomeMessage {
			c.JSON(http.StatusOK, newWelcomeMessage)
			return
		} else {
			c.JSON(http.StatusOK, defaultMessage)
			return
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	// graceful shutdown of active providers.
	openfeature.Shutdown()

	log.Println("Server exiting")
}
