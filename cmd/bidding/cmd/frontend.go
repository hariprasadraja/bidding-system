/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"net/http"
	"sellerapp-bidding-system/internal/auction"
	"sellerapp-bidding-system/internal/frontend"
	"sellerapp-bidding-system/internal/user"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"github.com/prometheus/common/log"

	httpServer "github.com/micro/go-plugins/server/http/v2"
	"github.com/spf13/cobra"
)

// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

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

	},
}

func init() {
	serviceCmd.AddCommand(frontendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// frontendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// frontendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
