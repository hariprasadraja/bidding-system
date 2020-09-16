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
	"log"
	"sellerapp-bidding-system/internal/auction"

	"github.com/micro/go-micro/v2"
	"github.com/spf13/cobra"
)

// auctionCmd represents the auction command
var auctionCmd = &cobra.Command{
	Use:   "auction",
	Short: "start the auction service server",
	Long:  `./bidding service auction`,
	Run: func(cmd *cobra.Command, args []string) {
		// New Service
		service := micro.NewService(
			micro.Name("go.micro.server.auction"),
			micro.Version("latest"),
		)

		// Initialise service
		service.Init()

		// Register Handler
		auction.RegisterAuctionHandler(service.Server(), new(auction.Handler))

		// Register Struct as Subscriber
		// micro.RegisterSubscriber("go.micro.service.auction", service.Server(), new(subscriber.Auction))

		// Run service
		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serviceCmd.AddCommand(auctionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// auctionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// auctionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
