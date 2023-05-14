package main

import (
	"os"
)

var (
	daprHttpPort  string
	daprGrpcPort  string
	port          string
	serviceName   string
	calledService string
	disableCall   bool
	severMode     string
	clientMode    string
)

func init() {
	daprHttpPort = os.Getenv("DAPR_HTTP_PORT")
	if daprHttpPort == "" {
		daprHttpPort = "3500"
	}
	daprGrpcPort = os.Getenv("DAPR_GRPC_PORT")
	if daprGrpcPort == "" {
		daprGrpcPort = "50001"
	}
	port = os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	serviceName = os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "service"
	}
	calledService = os.Getenv("CALLED_SERVICE")
	if calledService == "" {
		calledService = "service"
	}
	disableCall = os.Getenv("DISABLE_CALL") == "true"
	severMode = os.Getenv("SERVER_MODE")
	if severMode == "" {
		severMode = "http"
	}
	clientMode = os.Getenv("CLIENT_MODE")
	if clientMode == "" {
		clientMode = "http"
	}
}

func main() {
	if severMode == "http" {
		runHTTPServer()
	} else {
		runGRPCServer()
	}
}
