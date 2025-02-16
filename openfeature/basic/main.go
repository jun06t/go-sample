package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

const defaultMessage = "Hello!"
const newWelcomeMessage = "Hello, welcome to this OpenFeature-enabled website!"

func main() {
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

	engine.Run()
}
