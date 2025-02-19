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

	providerDomain = "my-domain"
	fraExpKey      = "fractional-evaluation"
	colorExpKey    = "color-experiment"
)

type BoolFlag struct {
	Result bool `json:"result"`
}

type StringFlag struct {
	Color string `json:"color"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := openfeature.SetNamedProviderAndWait(providerDomain, flagd.NewProvider())
	if err != nil {
		log.Fatalf("Failed to set the OpenFeature provider: %v", err)
	}
	client := openfeature.NewClient(providerDomain)

	engine := gin.Default()
	registerRoutes(engine, client)

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

func registerRoutes(engine *gin.Engine, client *openfeature.Client) {
	registerFractionEndpoint(engine, client)
	registerColorEndpoint(engine, client)
}

func registerFractionEndpoint(engine *gin.Engine, client *openfeature.Client) {
	engine.GET("/hello", func(c *gin.Context) {
		userId := c.Query("userId")

		evalCtx := openfeature.NewEvaluationContext(
			"",
			map[string]interface{}{
				"userId": userId,
			})

		flagResult, err := client.BooleanValue(
			context.Background(), fraExpKey, false, evalCtx,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if flagResult {
			c.JSON(http.StatusOK, BoolFlag{Result: true})
			return
		}
		c.JSON(http.StatusOK, BoolFlag{Result: false})
	})
}

func registerColorEndpoint(engine *gin.Engine, client *openfeature.Client) {
	engine.GET("/color", func(c *gin.Context) {
		age := c.Query("age")

		evalCtx := openfeature.NewEvaluationContext(
			"",
			map[string]interface{}{
				"age": age,
			})

		flagResult, err := client.StringValue(
			context.Background(), colorExpKey, "", evalCtx,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, StringFlag{Color: flagResult})
	})
}
