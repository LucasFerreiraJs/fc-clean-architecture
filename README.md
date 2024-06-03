# desafio-Clean Architecture


### Rodando o projeto



Execute o comando na raiz do projeto para subir os containers:

```
go mod tidy
docker-compose up -d

```

Será necessário criar a tabela orders:

```
docker exec -it mysql bash
-- ou com id do container listando com docker ps

mysql -uroot -p orders
-- senha root

CREATE DATABASE orders;
CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id));

SELECT * FROM orders;

```




Para rodar a aplicação, acesse a pasta cmd/ordersystem e execute:

```
go run main.go wire_gen.go
```


### gRPC:
Executar evans:
```
evans -r repl
```
chamadas para criar e listar orders:
  call CreateOrder
  call GetOrders


### GraphQL
acesse por`localhost:8080`

mutation para criação de order:
```
mutation createOrder {
  createOrder(input: {id: "", Price: 0.0, Tax: 0.0}) {
    id,
    Price,
    Tax,
    FinalPrice
  }
}
```

Query para listar orders:
```
query orders {
   orders{
    id,
    FinalPrice,
    Tax,
    Price
  }
}

```
### Web

listar orders com `localhost:8000/order`