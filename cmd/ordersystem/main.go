package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/fc-clean-architecture/configs"
	"github.com/devfullcycle/fc-clean-architecture/internal/event/handler"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/graph"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/grpc/pb"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/grpc/service"
	"github.com/devfullcycle/fc-clean-architecture/internal/infra/web/webserver"
	"github.com/devfullcycle/fc-clean-architecture/pkg/events"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	webserver.Router.Use(middleware.Logger)
	webserver.Router.Post("/order", webOrderHandler.Create)
	webserver.Router.Get("/order", webOrderHandler.GetOrders)
	// webserver.AddHandler("/order", webOrderHandler.Create)
	// webserver.AddHandler("/order", webOrderHandler.GetOrders)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	createGetOrdersUseCase := NewGetOrdersUseCase(db)

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *createGetOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)

	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase:     *createOrderUseCase,
		CreateGetOrdersUseCase: *createGetOrdersUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(configs *configs.Conf) *amqp.Channel {

	rabbitAccess := fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.RabbitMQUser, configs.RabbitMQPass, configs.RabbitMQHost, configs.RabbitMQPort)
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rabbitAccess)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
