package golexa

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

// Start turns on the appropriate listener to handle incoming requests for the given skill. It will
// automatically detect if you're just running the process on your local machine or some testing cluster
// or if you are actually running it as an AWS Lambda function. When not running in Lambda it will simply
// fire up an HTTP server that listens on port 20123 for equivalent incoming blocks of JSON as you'd
// receive from the Alexa API to your Lambda function. You can use a different port by setting
// the GOLEXA_HTTP_PORT environment variable.
//
// You should only call this once per process!
func Start(skill Skill) {
	wrappedHandler := lambda.NewHandler(skill.Handle)

	switch {
	case runningInLambda():
		startLambda(wrappedHandler)
	default:
		startHttp(wrappedHandler)
	}
}

// runningInLambda determines if we're running in a live AWS Lambda environment or if we appear
// to be running locally (thus should fall back to an HTTP listener).
func runningInLambda() bool {
	// https://docs.aws.amazon.com/lambda/latest/dg/lambda-environment-variables.html
	executionEnv := os.Getenv("AWS_EXECUTION_ENV")
	return executionEnv != "" && strings.HasPrefix(executionEnv, "AWS_Lambda_")
}

// startLambda begins listening for requests to your handler as though you're running on AWS Lambda.
func startLambda(handlerFunc lambda.Handler) {
	lambda.StartHandler(handlerFunc)
}

// startHttp fires up an HTTP server that routes all requests to your handler function. It assumes
// that all of your inputs/outputs are JSON.
func startHttp(handlerFunc lambda.Handler) {
	listener := httpListener{handlerFunc: handlerFunc}
	port := listener.port()

	fmt.Println(fmt.Sprintf(`{"logger": "golexa", "msg": "Starting golexa on HTTP port %v"}`, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), listener))
}

// httpListener is only used for local development to simulate the transport work done by AWS Lambda.
type httpListener struct {
	handlerFunc lambda.Handler
}

func (listener httpListener) port() uint16 {
	portVar := os.Getenv("GOLEXA_HTTP_PORT")
	if portVar == "" {
		return 20123
	}
	port, err := strconv.ParseUint(portVar, 10, 16)
	if err != nil {
		return 20123
	}
	return uint16(port)
}

func (listener httpListener) ServeHTTP(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf(`{"logger": "golexa", "msg": "panic: %v"}`, err))
			http.Error(httpWriter, "unexpected panic in skill", http.StatusInternalServerError)
		}
	}()

	// Obtain the raw bytes from the request body - equivalent to the raw bytes lambda would receive.
	jsonInput, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		msg := fmt.Sprintf("Unable to read input [%v]", err)
		http.Error(httpWriter, msg, http.StatusBadRequest)
		return
	}

	// The Lambda library's wrapped handlers already take the raw Go struct returned by the
	// original h and marshal it into raw JSON bytes. We just need to send those bytes
	// back to the caller.
	jsonOutput, err := listener.handlerFunc.Invoke(context.Background(), jsonInput)
	if err != nil {
		msg := fmt.Sprintf("Unable handle request [%v]", err)
		http.Error(httpWriter, msg, http.StatusInternalServerError)
		return
	}

	httpWriter.Header().Set("Content-Type", "application/json")
	httpWriter.WriteHeader(200)
	_, _ = httpWriter.Write(jsonOutput)
}
