package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler := os.Getenv("_HANDLER")

	if handler == "" {
		// Running locally
		e.Logger.Fatal(e.Start(":8080"))
	} else {
		// Running on lambda
		lambda.Start(echoadapter.New(e).ProxyWithContext)
	}
}
