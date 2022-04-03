image:
	docker build -t restaurant-assistant-image:v1 .

container:
	docker run --name restaurant-assistant -p 8080:80 -p 58080:50080/tcp --env-file .env restaurant-assistant-image:v1

proto-ra:
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/order-ra.proto

proto-cs:
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/order-cs.proto

proto-auth:
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/auth.proto

proto-fd:
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/order-fd.proto

proto-manager:
	evans api/proto/manager.proto -p 50080

proto-restaurant:
	protoc -I api/proto --go_out=. --go-grpc_out=. api/proto/restaurant.proto

swag-generate:
	swag init -g cmd/main.go

run:
	go run cmd/main.go
