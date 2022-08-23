# Dealer

![CI](https://github.com/Misphix/dealer/actions/workflows/test.yml/badge.svg)
## Run Service
There are two ways to run the service.
- Run with local machine
- Run with docker-composer

### Run with Local Machine
There are three steps to do if you want to run the service on the local machine. You can use the service at `localhost:8626` after running the following commands.
1. `make env` to run MySQL and RabbitMQ
2. `make` to build the execution file
3. `./dealer` to execute the service

### Run with Docker-Composer
There is only one step to do if you want to run the service in the docker. You can use the service at `localhost:8626` after running the following commands.
1. `make deploy` to run the service and all dependencies


## API

### New an Order
- Method: POST
- Path: `localhost:8626/v1/order`
- Body: json format
    - order_type `int`: order type
        - 1: buy
        - 2: sell
    - quantity `int`: quantity
    - price_type `int`: price type
        - 1: limit price
        - 2: market price
    - price `float64` (optional): price
- Response: json format
    - id `int`: order ID
    - order_type `int`: order type
        - 1: buy
        - 2: sell
    - quantity `int`: quantity
    - remain_quantity `int`: remain quantity
    - price_type `int`: price type
        - 1: limit price
        - 2: market price
    - price `float64`: price
    - is_cancel `bool`: is order cancel

#### Example
```sh
curl --location --request POST 'localhost:8626/v1/order' \
--header 'Content-Type: application/json' \
--data-raw '{
    "order_type": 1,
    "quantity": 1,
    "price_type": 1,
    "price": 5
}'
```

### Cancel an Order
- Method: DELETE
- Path: `localhost:8626/v1/order/:id`

#### Example
```sh
curl --location --request DELETE 'localhost:8626/v1/order/1'
```

## System Design
這個系統分成兩部分，一部分是接收訂單的http server，另一部分是處理訂單的consumer。

接收訂單的http server會把訂單的資訊寫進DB之中，再把訂單資訊publish進RabbitMQ之中。consumer會把訂單的資訊(新增或取消)消費下來，放到系統之中去進行撮合。

在database之中可以看到目前有哪些order和有哪些deal。所有客戶下的單都在order這張table之中查到，包含是否逹成、有沒有被取消。而在deal的table中可以查看有哪些交易。

目前的http server和consumer都寫在同一個main之中，若有需要可以再進行拆分。