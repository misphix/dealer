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
