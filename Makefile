all: clean client server done

.PHONY: all client server clean proto mocks cover done

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
BUILDS_FOLDER=builds
LD_FLAGS="-X \"${BUILD_DATE}\" -X \"${BUILD_VER}\" -X \"${BUILD_GIT}\""

MOCKS_DESTINATION="mocks"

client:
	@echo "Building client elf/exe in \"builds\" folder ..."
	@mkdir -p builds
	@go build -o $(BUILDS_FOLDER)/client.elf -ldflags $(LD_FLAGS) cmd/client/client.go
	@GOOS=windows go build -o $(BUILDS_FOLDER)/client.exe -ldflags $(LD_FLAGS) cmd/client/client.go


server:
	@echo "Building server elf/exe in \"builds\" folder ..."
	@mkdir -p builds
	@go build -o $(BUILDS_FOLDER)/server.elf cmd/server/server.go
	@GOOS=windows go build -o $(BUILDS_FOLDER)/server.exe cmd/server/server.go

done:
	@echo Done!

clean:
	@rm -rf builds
	@rm -rf $(MOCKS_DESTINATION)


cover: mocks
	@#for p in `go list ./... | grep -viE "(cmd|tui|proto)"`; do echo -en "$p/... "; done
	@echo "Checking coverage ..."
	@go test -coverprofile=cover.out passkeeper/internal/entities/... passkeeper/internal/entities/compress/... passkeeper/internal/entities/config/... passkeeper/internal/entities/config/client/... passkeeper/internal/entities/config/server/... passkeeper/internal/entities/hashes/... passkeeper/internal/entities/logger/... passkeeper/internal/entities/myerrors/... passkeeper/internal/entities/structs/... passkeeper/internal/server/... passkeeper/internal/storage/... passkeeper/internal/storage/dbstorage/postgres/... passkeeper/internal/storage/memstorage/... passkeeper/internal/transport/grpc/handlers/... passkeeper/internal/transport/grpc/interceptors/... passkeeper/internal/usecase/cli/... passkeeper/internal/usecase/cli/filesUC/... passkeeper/internal/usecase/srv/blobsUC/... passkeeper/internal/usecase/srv/usersUC/...
	@go tool cover -html=cover.out -o cover.html
	@go tool cover -func=cover.out | grep -i total | tr -s '\t'


mocks:
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@mockgen -source internal/storage/storage.go -destination mocks/mockstorage/storage_mock.go -package=mocks
	@mockgen -source internal/transport/grpc/protofiles/proto/users_change_pwd_grpc.pb.go -destination mocks/gservice/user_change_pass.go -package=mocks
	@mockgen -source internal/transport/grpc/protofiles/proto/users_service_grpc.pb.go -destination mocks/gservice/user_auth.go -package=mocks
	@mockgen -source internal/transport/grpc/protofiles/proto/credentials_service_grpc.pb.go -destination=mocks/gservice/blobs.go -package=mocks
	@mockgen -source internal/usecase/srv/blobsUC/blobs_usecase.go -destination mocks/usecase/blobs_usecase/blobs_usecase.go -package mocks
	@mockgen -source internal/usecase/srv/usersUC/users_usecase.go -destination mocks/usecase/users_usecase/users_usecase.go -package mocks
	@#@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done



tools:
	@go install github.com/golang/mock/mockgen@1.6.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

