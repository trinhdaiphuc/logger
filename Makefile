gen-protobuf:
	protoc -Iproto/ proto/*.proto --go_out=plugins=grpc:.

clean:
	rm -rf coverage && mkdir coverage

test: clean
	go test -race -covermode=atomic -coverprofile=coverage/c.out && \
	go tool cover -html=coverage/c.out
