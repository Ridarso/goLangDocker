# this will call the protoc library, which is responsibly for compiling your protobuf definitions into code
build:
	# protoc -I. --go_out=plugins=grpc:. proto/consignment/*.proto
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/docker/consignment-service/ \
		proto/consignment/*.proto
	docker build -t consignment-service .
run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns consignment-service