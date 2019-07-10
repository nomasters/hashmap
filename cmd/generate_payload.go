/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"io/ioutil"
	"log"
	"time"

	payload "github.com/nomasters/hashmap/pkg/payload"
	sigutil "github.com/nomasters/hashmap/pkg/sig/sigutil"
	"github.com/spf13/cobra"
)

var message string
var ttl string
var timestamp int64
var keysetPath string
var outputPath string

// generatePayloadCmd represents the generatePayload command
var generatePayloadCmd = &cobra.Command{
	Use:   "payload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		signerBytes, err := ioutil.ReadFile(keysetPath)
		if err != nil {
			// TODO write a more meaningful error message
			log.Fatal(err)
		}
		signers, err := sigutil.Decode(signerBytes)
		if err != nil {
			// TODO write a more meaningful error message
			log.Fatal(err)
		}

		t, err := time.ParseDuration(ttl)
		if err != nil {
			// TODO write a more meaningful error message
			log.Fatal(err)
		}

		p, err := payload.Generate(
			[]byte(message),
			signers,
			payload.WithTTL(t),
			payload.WithTimestamp(time.Unix(0, timestamp)),
		)
		if err != nil {
			log.Fatal(err)
		}
		b, err := payload.Marshal(p)
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(outputPath, b, 0600); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(generatePayloadCmd)

	generatePayloadCmd.Flags().StringVarP(&message, "message", "m", "", "The message to be stored in data of payload")
	generatePayloadCmd.Flags().StringVarP(&ttl, "ttl", "t", payload.DefaultTTL.String(), "ttl in XXhXXmXXs string format. Defaults to 24 hours")
	generatePayloadCmd.Flags().Int64VarP(&timestamp, "timestamp", "s", time.Now().UnixNano(), "timestamp for message in unix-nano time. Defaults now")
	generatePayloadCmd.Flags().StringVarP(&keysetPath, "keyset", "k", "hashmap.keyset", "the path for the keyset file. Defaults `./hashamp.keyset`")
	generatePayloadCmd.Flags().StringVarP(&outputPath, "output", "o", "payload.protobuf", "the path for the output payload file. Defaults `./payload.protobuf`")
}
