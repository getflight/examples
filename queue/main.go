package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getflight/core"
	"github.com/getflight/core/queue"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

// In order to invoke this function, once the deployment is complete, send any message on the trigger queue configured for this function
// For example : to trigger the function in the staging environment, send a message to the queue-trigger-staging SQS queue
// The queue name will vary depending on the name of the project and the environment : {project}-trigger-{environment}
func main() {
	log.SetLevel(log.DebugLevel)

	// Initialize core, this will initialize configuration and allow us to fetch queue messages from AWS
	err := core.Init()

	if err != nil {
		fmt.Printf("error initializing core: %s\n", err)
		os.Exit(1)
	}

	handler := os.Getenv("_HANDLER")

	if handler == "" {
		// Running locally
		// To invoke the function locally, one option is to manually create a specific
		// queue for the local environment in SQS that local environments will listen on.
		queue.Start("queue-trigger-local", localHandler)

	} else {
		// Running on lambda
		// Flight will configure the SQS trigger queue for each environment to allow the function invocation.
		lambda.Start(lambdaHandler)
	}
}

func lambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("no message passed to function")
	}

	for _, msg := range sqsEvent.Records {
		log.Debugf("incomming lambda message %q with body %q", msg.MessageId, msg.Body)
	}

	return nil
}

func localHandler(body string) error {
	log.Debugf("incomming local message with body %q", body)

	return nil
}
