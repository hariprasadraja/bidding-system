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
	"fmt"
	"log"
	"sellerapp-bidding-system/internal/app/sellerapp/auction"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	httpServer "github.com/micro/go-plugins/server/http/v2"
	"github.com/spf13/cobra"
)

// auctionCmd represents the auction command
var auctionCmd = &cobra.Command{
	Use:   "auction",
	Short: "start a new micro service to provide auction apis",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auction called")
		srv := httpServer.NewServer(
			server.Name("micro.sellerapp.auction"),
			server.Version("1.0"),
		)

		router := httprouter.New()
		auction.AddRoutes(router)

		// mux := http.NewServeMux()
		// mux.
		// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	w.Write([]byte(`hello world`))
		// })

		hd := srv.NewHandler(router)

		err := srv.Handle(hd)
		if err != nil {
			log.Println(err)
			return
		}

		service := micro.NewService(
			micro.Name("micro.sellerapp"),
			micro.Server(srv),
		)

		service.Init()
		err = service.Run()
		if err != nil {
			log.Println(err)
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
