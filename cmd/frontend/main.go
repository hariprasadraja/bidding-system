package main

import (
	"net/http"
	"sellerapp-bidding-system/internal/auction"
	"sellerapp-bidding-system/internal/frontend"
	"sellerapp-bidding-system/internal/user"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	httpServer "github.com/micro/go-plugins/server/http/v2"
)

func main() {

	srv := httpServer.NewServer(
		server.Name("go.micro.server.frontend"),
		server.Version("latest"),
		server.Address(":8085"),
	)

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		log.Info("server started")
		_, err := w.Write([]byte(`hello world`))
		if err != nil {
			log.Error("unable to write", err)
		}
	})

	err := srv.Handle(srv.NewHandler(router))
	if err != nil {
		log.Error("unable to register handler", err)
	}

	service := micro.NewService(
		micro.Server(srv),
	)

	service.Init()

	userClient := user.NewUserService("go.micro.server.user", service.Client())
	auctionClient := auction.NewAuctionService("go.micro.server.auction", service.Client())
	frontend.RegisterUserRoutes(router, userClient)
	frontend.RegisterAuctionRoutes(router, userClient, auctionClient)

	err = service.Run()
	if err != nil {
		log.Error("unable to start the service", err)
	}
}
