        export PATH=$PATH:~/go/bin

proto:

        protoc -I . \
        --go_out ./pb/. --go_opt paths=source_relative \
        --go-grpc_out ./pb/ --go-grpc_opt paths=source_relative \
        ./proto/api.proto

server:

        go run server/gServer.go

client:

        go run client/gClient.go -sender Dipesh -channel default
