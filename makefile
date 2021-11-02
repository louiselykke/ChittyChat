gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto

cli:
	go run client/client.go

serv:
	go run server/server.go

cli1:
	go run client/client.go -user emre

cli2:
	go run client/client.go -user bjørn

cli3:
	go run client/client.go -user louise