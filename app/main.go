package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// Wrapper from events.APIGatewayProxyRequest to http.Handler.
var httpLambda *httpadapter.HandlerAdapter // nolint:gochecknoglobals // Don't want to remake for each run.

// init Sets up the proxy handler once at the start.
func init() { // nolint: gochecknoinits // Limiting scope for interview only.
	server := getServer()
	httpLambda = httpadapter.New(server)
}

// getServer Create a http.ServerMux that contains our routes.
func getServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/dynamo", dynamoHandler)

	return mux
}

// Handler Take a events.APIGatewayProxyRequest and uses a proxy to access a http.ServerMux
// See https://github.com/awslabs/aws-lambda-go-api-proxy
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req) // nolint:wrapcheck // No reason to wrap this final error.
}

func main() {
	lambda.Start(Handler)
}
