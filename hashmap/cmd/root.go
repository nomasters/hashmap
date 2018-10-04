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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hashmap",
	Short: "a light-weight cryptographically signed key value store inspired by IPNS",
	Long:  `a light-weight cryptographically signed key value store inspired by IPNS`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is in the current directory: ./hashmap.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// get the configuration file from the current directory
		viper.AddConfigPath(".")
		viper.SetConfigName("hashmap")
	}

	// support ENV for hashmap settings
	viper.SetEnvPrefix("hashmap")
	viper.BindEnv("server.host", "HASHMAP_SERVER_HOST")
	viper.BindEnv("server.port", "HASHMAP_SERVER_PORT")
	viper.BindEnv("server.tls", "HASHMAP_SERVER_TLS")
	viper.BindEnv("server.certfile", "HASHMAP_SERVER_CERTFILE")
	viper.BindEnv("server.keyfile", "HASHMAP_SERVER_KEYFILE")
	viper.BindEnv("storage.engine", "HASHMAP_STORAGE_ENGINE")
	viper.BindEnv("storage.endpoint", "HASHMAP_STORAGE_ENDPOINT")
	viper.BindEnv("storage.auth", "HASHMAP_STORAGE_AUTH")
	viper.BindEnv("storage.maxIdle", "HASHMAP_STORAGE_MAXIDLE")
	viper.BindEnv("storage.maxActive", "HASHMAP_STORAGE_MAXACTIVE")
	viper.BindEnv("storage.idleTimeout", "HASHMAP_STORAGE_IDLETIMEOUT")
	viper.BindEnv("storage.wait", "HASHMAP_STORAGE_WAIT")
	viper.BindEnv("storage.maxConnLifetime", "HASHMAP_STORAGE_MAXCONNLIFETIME")
	viper.BindEnv("storage.tls", "HASHMAP_STORAGE_TLS")

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
