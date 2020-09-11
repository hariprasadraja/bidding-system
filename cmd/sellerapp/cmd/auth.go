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
	"sellerapp-bidding-system/internal/app/sellerapp/auth"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	httpServer "github.com/micro/go-plugins/server/http/v2"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth starts the authentication http server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := httpServer.NewServer(
			server.Name("auth"),
			server.Version("1.0"),
		)

		router := httprouter.New()
		auth.AddRoutes(router)

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
			micro.Name("go.micro"),
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
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
