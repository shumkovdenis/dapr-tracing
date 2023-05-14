package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	client http.Client
)

func runHTTPServer() {
	r := mux.NewRouter()
	r.HandleFunc("/echo", echoHandle).Methods("POST")

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != http.ErrServerClosed {
		log.Fatalln("Error starting HTTP server")
	}
}

func echoHandle(w http.ResponseWriter, r *http.Request) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body:", err.Error())
	}

	log.Println("Received:", string(in))

	traceparent := r.Header.Get("traceparent")

	log.Println("traceparent:", traceparent)

	out := fmt.Sprintf("%s %s", string(in), serviceName)

	if !disableCall {
		if clientMode == "http" {
			out = callHTTPService(traceparent, out)
		} else {
			out = callGRPCService(context.Background(), out)
		}
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		log.Fatal("Error writing the response:", err.Error())
	}
}

func callHTTPService(traceparent, body string) string {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("http://localhost:%s/echo", daprHttpPort),
		strings.NewReader(body),
	)
	if err != nil {
		log.Fatal("Error creating HTTP request:", err.Error())
	}

	req.Header.Add("dapr-app-id", calledService)

	if traceparent != "" {
		req.Header.Add("traceparent", traceparent)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Error calling:", err.Error())
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response:", err.Error())
	}
	response.Body.Close()

	return string(result)
}
