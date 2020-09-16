# Bidding App

Application to manage Auction Management


## Build

```sh
 go build ./cmd/bidding
```

## CMD

*bidding* command has a sub command called *service* which will start it's child services

so, far we have three services

#### Auction Service
``` sh
protoc --proto_path=.:$GOPATH/src --go_out=./internal --micro_out=./internal internal/auction/auction.proto && go build ./cmd/bidding && ./bidding service auction --config configs/dev.bidding.yaml
```