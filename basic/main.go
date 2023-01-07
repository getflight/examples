package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	// Basic root endpoint
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		body := fmt.Sprintf("Hello, World!\n")

		_, err := io.WriteString(writer, body)

		if err != nil {
			fmt.Printf("error writing response: %s\n", err)
		}
	})

	// Variables are converted to snake case in caps.
	// ex : my-key-1 is converted to MY_KEY_1
	http.HandleFunc("/variables", func(writer http.ResponseWriter, request *http.Request) {

		body := fmt.Sprintf("Variables: \n"+
			"variable MY_KEY_1 : %s\n"+
			"variable MY_KEY_2 : %s\n"+
			"variable MY_KEY_2 : %s\n",
			os.Getenv("MY_KEY_1"),
			os.Getenv("MY_KEY_2"),
			os.Getenv("MY_KEY_3"))

		_, err := io.WriteString(writer, body)

		if err != nil {
			fmt.Printf("error writing response: %s\n", err)
		}
	})

	// Files are included in the packaged artifact.
	// With this configuration, all files in the resources directory will be included in the deployment.
	http.HandleFunc("/files", func(writer http.ResponseWriter, request *http.Request) {

		var files []string
		err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

		var body string

		for _, f := range files {
			body = fmt.Sprintf("%s\n%s\n", body, f)
		}

		_, err = io.WriteString(writer, body)

		if err != nil {
			fmt.Printf("error writing response: %s\n", err)
		}
	})

	handler := os.Getenv("_HANDLER")

	if handler == "" {
		// Running locally

		err := http.ListenAndServe(":8080", nil)

		if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}

	} else {
		// Running on lambda
		lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
	}
}
