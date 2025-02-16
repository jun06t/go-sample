package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-feature/go-sdk/openfeature"
)

const defaultMessage = "Hello!"
const newWelcomeMessage = "Hello, welcome to this OpenFeature-enabled website!"

func main() {
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
