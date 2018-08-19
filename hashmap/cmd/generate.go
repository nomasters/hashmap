// This is free and unencumbered software released into the public domain.

// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.

// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

// For more information, please refer to <http://unlicense.org>

package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nomasters/hashmap"
	"github.com/spf13/cobra"
)

var message string
var ttl int64
var timestamp int64

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates data such as private keys and payloads",
	Long:  `generates data such as private keys and payloads`,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "key":
			fmt.Println(base64.StdEncoding.EncodeToString(hashmap.GenerateKey()))
		case "payload":

			opts := hashmap.Options{
				Message:   message,
				TTL:       ttl,
				Timestamp: timestamp,
			}

			text, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalln(err)
			}
			pk, err := base64.StdEncoding.DecodeString(string(text))
			if err != nil {
				log.Fatal(err)
			}

			payload, err := hashmap.GeneratePayload(opts, pk)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(payload))

		default:
			fmt.Println("invalid input. must use payload or key subcommand")
		}

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&message, "message", "m", "", "The message to be stored in data of payload")
	generateCmd.PersistentFlags().Int64VarP(&ttl, "ttl", "t", hashmap.DataTTLDefault, "ttl in seconds for payload")
	generateCmd.PersistentFlags().Int64VarP(&timestamp, "timestamp", "s", time.Now().Unix(), "timestamp for message in unix-time")
}
