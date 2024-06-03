

go run main.go wire_gen.go

evans --proto internal/infra/grpc/protofiles/order.proto repl




protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles//order.proto




docker container ps

docker exec -it 15d35377dbaf bash


mysql -uroot -p orders
root