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
GIT_TAG=`git tag -l | tail -n 1`
GIT_VER=`git log --oneline | head -n1 | cut -f1 -d' '`
BUILD_VER=main.BuildVersion=${GIT_TAG}
BUILD_GIT=main.BuildCommit=${GIT_VER}
BUILD_DATE=main.BuildDate=${DATE}

LD_FLAGS="-X \"${BUILD_DATE}\" -X \"${BUILD_VER}\" -X \"${BUILD_GIT}\""

client:
	mkdir -p builds
	go build -o builds/client.elf -ldflags ${LD_FLAGS} cmd/client/client.go
	GOOS=windows go build -o builds/client.exe cmd/client/client.go

server:
	mkdir -p builds
	go build -o builds/server.elf cmd/server/server.go
	GOOS=windows go build -o builds/server.exe cmd/server/server.go

clean:
	@rm -rf builds

cover:
	#for p in `go list ./... | grep -viE "(cmd|tui|proto)"`; do echo -en "$p/... "; done
	@go test -coverprofile=cover.out passkeeper/internal/entities/... passkeeper/internal/entities/compress/... passkeeper/internal/entities/config/... passkeeper/internal/entities/config/client/... passkeeper/internal/entities/config/server/... passkeeper/internal/entities/hashes/... passkeeper/internal/entities/logger/... passkeeper/internal/entities/myerrors/... passkeeper/internal/entities/structs/... passkeeper/internal/server/... passkeeper/internal/storage/... passkeeper/internal/storage/dbstorage/postgres/... passkeeper/internal/storage/memstorage/... passkeeper/internal/transport/grpc/handlers/... passkeeper/internal/transport/grpc/interceptors/... passkeeper/internal/usecase/cli/... passkeeper/internal/usecase/cli/filesUC/... passkeeper/internal/usecase/srv/blobsUC/... passkeeper/internal/usecase/srv/usersUC/...
	@go tool cover -html=cover.out -o cover.html
	@go tool cover -func=cover.out | grep -i total | tr -s '\t'

MOCKS_DESTINATION="mocks"

mocks:
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	mockgen -source internal/storage/storage.go -destination mocks/storage_mock.go -package=mocks
	#@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done


tools:
	@go install github.com/golang/mock/mockgen@1.6.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

