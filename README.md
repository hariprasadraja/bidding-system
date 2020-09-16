# Bidding App

A Simple Auction Management System

## Dependencies

Go Version: 1.14
Database: MYSQL 8
Micro Framework - https://github.com/micro/micro (Version 2)

## Build

```sh
 go build ./cmd/bidding
```

## CMD

*bidding* command has a sub command called *service* which will start it's child services

As per requirement, It has three services, you can run the given commands as in order.

#### Micro Server

First Setup the micro server, which provides service discovery and registration

```sh
micro server
```


#### Auction Service

Hope, you have already taken the build, execute the below command and start the auction micro service.

``` sh
protoc --proto_path=.:$GOPATH/src --go_out=./internal --micro_out=./internal internal/auction/auction.proto && go build ./cmd/bidding && ./bidding service auction --config configs/dev.bidding.yaml
```

#### User Service

Run the below command to start the User micro service.

```sh
protoc --proto_path=.:$GOPATH/src --go_out=./internal --micro_out=./internal internal/user/user.proto && go build ./cmd/bidding && ./bidding service user --config configs/dev.bidding.yaml
```

#### API Server

This is the API server which serves apis to the frontend

``` sh
go build ./cmd/bidding && ./bidding service frontend
```


## Docs

Added Insomnia File in the docs folder. you can open it with Insomnia or using Postman.

## Scripts

/scripts directory has sql script to initialize the database