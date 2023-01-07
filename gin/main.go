package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	handler := os.Getenv("_HANDLER")

	if handler == "" {
		// Running locally
		r.Run()
	} else {
		// Running on lambda
		lambda.Start(ginadapter.New(r).ProxyWithContext)
	}
}
