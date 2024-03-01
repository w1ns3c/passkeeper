

proto:
	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
            --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      		internal/transport/grpc/protofiles/sessionkey.proto

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
      		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
      		internal/transport/grpc/protofiles/users_service.proto

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
          		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
          		internal/transport/grpc/protofiles/users_service.proto

	/usr/bin/protoc --go_out=. --go_opt=paths=source_relative \
              		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
              		internal/transport/grpc/protofiles/users_change_pwd.proto

