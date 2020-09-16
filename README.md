# Bidding App

A Simple Auction Management System

## Dependencies

Go Version: 1.14
Database: MYSQL 8

## Build

```sh
 go build ./cmd/bidding
```

## CMD

*bidding* command has a sub command called *service* which will start it's child services

As per requirement, It has three services

#### Auction Service
``` sh
protoc --proto_path=.:$GOPATH/src --go_out=./internal --micro_out=./internal internal/auction/auction.proto && go build ./cmd/bidding && ./bidding service auction --config configs/dev.bidding.yaml
```

#### User Service

```sh
protoc --proto_path=.:$GOPATH/src --go_out=./internal --micro_out=./internal internal/user/user.proto && go build ./cmd/bidding && ./bidding service user --config configs/dev.bidding.yaml
```

#### API Server
``` sh
go build ./cmd/bidding && ./bidding service frontend
```


## Docs

Add Insomnia File in the docs folder. you can open it with Insomnia or using Postman.