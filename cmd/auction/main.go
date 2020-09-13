package main

import (
	auction "sellerapp-bidding-system/internal/auction"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.auction"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	auction.RegisterAuctionHandler(service.Server(), new(auction.Auction))

	// Register Struct as Subscriber
	// micro.RegisterSubscriber("go.micro.service.auction", service.Server(), new(subscriber.Auction))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
