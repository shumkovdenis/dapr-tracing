clear:
	rm -rf ./.dapr/logs

run-http: clear
	dapr run -f ./http.yaml

stop-http:
	dapr stop -f ./http.yaml

run-grpc: clear
	dapr run -f ./grpc.yaml

stop-grpc:
	dapr stop -f ./grpc.yaml

echo:
	curl "http://localhost:3501/v1.0/invoke/service-a/method/echo" \
		-H "Content-type: application/json" \
		-d 'hello'
