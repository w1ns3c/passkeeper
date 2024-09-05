all: client server

proto:
	rm -f internal/transport/grpc/protofiles/*.go

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
            --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      		internal/transport/grpc/protofiles/sessionkey.proto

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
      		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
      		internal/transport/grpc/protofiles/credentials_service.proto

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
          		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
          		internal/transport/grpc/protofiles/users_service.proto

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
              		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
              		internal/transport/grpc/protofiles/users_change_pwd.proto

	mv internal/transport/grpc/protofiles/*.go internal/transport/grpc/protofiles/proto/

client:
	mkdir -p builds
	# client
	go build -o builds/client.elf cmd/client/client.go
	GOOS=windows go build -o builds/client.exe cmd/client/client.go
server:
	mkdir -p builds
	go build -o builds/server.elf cmd/server/server.go
	GOOS=windows go build -o builds/server.exe cmd/server/server.go