all: clean client server

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


DATE=`date -u '+%Y-%m-%d %H:%M:%S'`
GIT_TAG="${git tag -l | tail -n 1}"
VER=`if [ "${GIT_TAG}" = "" ]; then echo "N/A"; else echo ${GIT_TAG}; fi`
BUILD_VER=main.BuildVersion=${VER}
BUILD_DATE=main.BuildDate=${DATE}

client:
	mkdir -p builds
	go build -o builds/client.elf -ldflags "-X \"${BUILD_DATE}\" -X \"${BUILD_VER}\"" cmd/client/client.go
	GOOS=windows go build -o builds/client.exe cmd/client/client.go

server:
	mkdir -p builds
	go build -o builds/server.elf cmd/server/server.go
	GOOS=windows go build -o builds/server.exe cmd/server/server.go

clean:
	rm -rf builds