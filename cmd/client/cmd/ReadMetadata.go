/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"challenge/pkg/api"
	"challenge/util"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

// ReadMetadataCmd represents the ReadMetadata command
var ReadMetadataCmd = &cobra.Command{
	Use:   "ReadMetadata",
	Short: "Method for reading metadata from the request context.",
	Long: `Send any content as an argument, here it serves as a placeholder, 
and get the metadata corresponding to the "i-am-random-key" key.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ReadMetadata called")
		if len(args) != 1 {
			log.Fatal("incorrect number of arguments")
		}

		data := args[0]

		config, err := util.LoadConfig(".")
		if err != nil {
			log.Fatal("cannot load config:", err)
		}

		conn, err := grpc.Dial(config.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}

		c := api.NewChallengeServiceClient(conn)

		ctx := metadata.NewIncomingContext(
			context.Background(),
			metadata.Pairs("i-am-random-key", "i-am-random-key-value"),
		)

		resp, err := c.ReadMetadata(ctx, &api.Placeholder{Data: data})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Metadata: %s", resp.GetData())
	},
}

func init() {
	rootCmd.AddCommand(ReadMetadataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ReadMetadataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ReadMetadataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
