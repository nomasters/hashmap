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
	"fmt"
	"os"
	"time"

	"github.com/nomasters/hashmap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the hashmap server",
	Long:  `Runs the hashmap server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("starting hashmap server version: %v running on port: %v\n", hashmap.Version, viper.GetInt("server.port"))

		seCode, err := hashmap.GetStorageEngineCode(viper.GetString("storage.engine"))
		if err != nil {
			fmt.Println("invalide storage engine type, exiting")
			os.Exit(1)
		}

		opts := hashmap.ServerOptions{
			Host:     viper.GetString("server.host"),
			Port:     viper.GetInt("server.port"),
			TLS:      viper.GetBool("server.TLS"),
			CertFile: viper.GetString("server.certfile"),
			KeyFile:  viper.GetString("server.keyfile"),
			Storage: hashmap.StorageOptions{
				Engine:          seCode,
				Endpoint:        viper.GetString("storage.endpoint"),
				Auth:            viper.GetString("storage.auth"),
				MaxIdle:         viper.GetInt("storage.maxidle"),
				MaxActive:       viper.GetInt("storage.maxactive"),
				IdleTimeout:     time.Duration(viper.GetInt("storage.idletimeout")) * time.Second,
				Wait:            viper.GetBool("storage.wait"),
				MaxConnLifetime: time.Duration(viper.GetInt("storage.maxconnlifetime")) * time.Second,
				TLS:             viper.GetBool("storage.tls"),
			},
		}
		hashmap.Run(opts)
	},
}

// server struct is used for holding flag values
var server struct {
	host     string
	port     int
	tls      bool
	certFile string
	keyFile  string
}

// stroage struct is used for holding flag values
var storage struct {
	engine          string
	endpoint        string
	auth            string
	maxIdle         int
	maxActive       int
	idleTimeout     int
	wait            bool
	maxConnLifetime int
	tls             bool
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().IntVarP(&server.port, "port", "", hashmap.DefaultPort, "The server port.")
	viper.BindPFlag("server.port", runCmd.PersistentFlags().Lookup("port"))

	runCmd.PersistentFlags().StringVarP(&server.host, "host", "", "", "The server hostname or ip address.")
	viper.BindPFlag("server.host", runCmd.PersistentFlags().Lookup("host"))

	runCmd.PersistentFlags().BoolVarP(&server.tls, "tls", "", true, "Boolean switch for running server in TLS mode.")
	viper.BindPFlag("server.tls", runCmd.PersistentFlags().Lookup("tls"))

	runCmd.PersistentFlags().StringVarP(&server.certFile, "certfile", "", "", "Path the TLS cert file.")
	viper.BindPFlag("server.certfile", runCmd.PersistentFlags().Lookup("certfile"))

	runCmd.PersistentFlags().StringVarP(&server.keyFile, "keyfile", "", "", "Path the TLS key file.")
	viper.BindPFlag("server.keyfile", runCmd.PersistentFlags().Lookup("keyfile"))

	runCmd.PersistentFlags().StringVarP(&storage.engine, "engine", "", "memory", "Storage Engine Mode.")
	viper.BindPFlag("storage.engine", runCmd.PersistentFlags().Lookup("engine"))

	runCmd.PersistentFlags().StringVarP(&storage.endpoint, "endpoint", "", "", "Storage endpoint string.")
	viper.BindPFlag("storage.endpoint", runCmd.PersistentFlags().Lookup("endpoint"))

	runCmd.PersistentFlags().StringVarP(&storage.auth, "auth", "", "", "Storage auth string.")
	viper.BindPFlag("storage.auth", runCmd.PersistentFlags().Lookup("auth"))

	runCmd.PersistentFlags().IntVarP(&storage.maxIdle, "max-idle", "", 0, "Storage max idle in seconds.")
	viper.BindPFlag("storage.maxidle", runCmd.PersistentFlags().Lookup("max-idle"))

	runCmd.PersistentFlags().IntVarP(&storage.maxActive, "max-active", "", 0, "Storage max active in seconds.")
	viper.BindPFlag("storage.maxactive", runCmd.PersistentFlags().Lookup("max-active"))

	runCmd.PersistentFlags().IntVarP(&storage.idleTimeout, "idle-timeout", "", 0, "Storage session idle timeout in seconds.")
	viper.BindPFlag("storage.idletimeout", runCmd.PersistentFlags().Lookup("idle-timeout"))

	runCmd.PersistentFlags().BoolVarP(&storage.wait, "wait", "", true, "Storage session wait boolean.")
	viper.BindPFlag("server.wait", runCmd.PersistentFlags().Lookup("wait"))

	runCmd.PersistentFlags().IntVarP(&storage.maxConnLifetime, "max-conn-lifetime", "", 0, "Storage max connection lifetime in seconds.")
	viper.BindPFlag("storage.maxconnlifetime", runCmd.PersistentFlags().Lookup("max-conn-lifetime"))

	runCmd.PersistentFlags().BoolVarP(&storage.tls, "storage-tls", "", false, "Boolean switch for running storage in TLS mode.")
	viper.BindPFlag("storage.tls", runCmd.PersistentFlags().Lookup("storage-tls"))
}
