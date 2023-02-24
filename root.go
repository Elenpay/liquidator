/*
Copyright © 2023 Clovr Labs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "liquidator",
	Short: "A CLI tool to monitor and automate the liquidity of your LND channels",
	Run: func(cmd *cobra.Command, args []string) {

		//Cobra main

		log.Infoln("Starting liquidator")
		startLiquidator()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//OTEL Expanded vars
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", os.ExpandEnv("$OTEL_EXPORTER_OTLP_ENDPOINT"))
	os.Setenv("OTLP_RESOURCES_ATTRIBUTES", os.ExpandEnv("$OTLP_RESOURCES_ATTRIBUTES"))

	viper.AutomaticEnv() // read in environment variables that match

	rootCmd.Flags().String("nodesHosts", "", "Command separated list of hostname:port to connect to")
	viper.BindPFlag("nodesHosts", rootCmd.Flags().Lookup("nodesHosts"))

	rootCmd.Flags().String("nodesMacaroons", "", "Command separated list of macaroons used in nodesHosts in hex, in the same order of NODESHOSTS")
	viper.BindPFlag("nodesMacaroons", rootCmd.Flags().Lookup("nodesMacaroons"))

	rootCmd.Flags().String("nodesTLSCerts", "", "Command separated list of tls certs from LNDS in base64, in the same order of NODESHOSTS and NODESMACAROONS")
	viper.BindPFlag("nodesTLSCerts", rootCmd.Flags().Lookup("nodesTLSCerts"))

	rootCmd.Flags().String("loopdHosts", "", "Command separated list of hostname:port to connect to loopd, each position corresponds to the managed node in nodesHosts")
	viper.BindPFlag("loopdHosts", rootCmd.Flags().Lookup("loopdHosts"))

	rootCmd.Flags().String("loopdMacaroons", "", "Command separated list of macaroons used in loopdHosts in hex, in the same order of loopdHosts")
	viper.BindPFlag("loopdMacaroons", rootCmd.Flags().Lookup("loopdMacaroons"))

	rootCmd.Flags().String("loopdTLSCerts", "", "Command separated list of tls certs from loopd in base64, in the same order of loopdHosts and loopdMacaroons")
	viper.BindPFlag("loopdTLSCerts", rootCmd.Flags().Lookup("loopdTLSCerts"))

	rootCmd.Flags().String("pollingInterval", "15s", "Interval to poll data")
	viper.BindPFlag("pollingInterval", rootCmd.Flags().Lookup("pollingInterval"))

	rootCmd.Flags().String("logLevel", "info", "Log level from values: {trace, debug, info, warn, error, fatal, panic}")
	viper.BindPFlag("logLevel", rootCmd.Flags().Lookup("logLevel"))

	rootCmd.Flags().String("logFormat", "text", "Log format from: {text, json}")
	viper.BindPFlag("logFormat", rootCmd.Flags().Lookup("logFormat"))

	//Flags for nodeguard grpc endpoint
	rootCmd.Flags().String("nodeguardHost", "", "Hostname:port to connect to nodeguard")
	viper.BindPFlag("nodeguardHost", rootCmd.Flags().Lookup("nodeguardHost"))

	//If nodesHosts length is different than nodesMacaroons, exit
	if len(viper.GetStringSlice("nodesHosts")) != len(viper.GetStringSlice("nodesMacaroons")) {
		log.Fatal("nodesHosts and nodesMacaroons must have the same length")
	}

	//Now we set the global vars
	nodesHosts = strings.Split(viper.GetString("nodesHosts"), ",")
	nodesMacaroons = strings.Split(viper.GetString("nodesMacaroons"), ",")
	nodesTLSCerts = strings.Split(viper.GetString("nodesTLSCerts"), ",")
	pollingInterval = viper.GetDuration("pollingInterval")
	nodeguardHost = viper.GetString("nodeguardHost")
	loopdHosts = strings.Split(viper.GetString("loopdHosts"), ",")
	loopdMacaroons = strings.Split(viper.GetString("loopdMacaroons"), ",")
	loopdTLSCerts = strings.Split(viper.GetString("loopdTLSCerts"), ",")

	// //Check that nodeguardHost is not empty
	// if nodeguardHost == "" {
	// 	log.Fatal("nodeguardHost is empty")
	// }

	//Set log level and format

	logLevel, err := log.ParseLevel(viper.GetString("logLevel"))
	if err != nil {
		log.Fatal("Invalid log level")
	}

	log.SetLevel(logLevel)

	if viper.GetString("logFormat") == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else if viper.GetString("logFormat") == "text" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.Fatal("Invalid log format")
	}

	// Log debug of the config
	log.Debug("nodesHosts: ", nodesHosts)
	log.Debug("nodesMacaroons: ", nodesMacaroons)
	log.Debug("nodesTLSCerts: ", nodesTLSCerts)
	log.Debug("pollingInterval: ", pollingInterval)
	log.Debug("logLevel: ", logLevel)
	log.Debug("logFormat: ", viper.GetString("logFormat"))
	log.Debug("nodeguardHost: ", nodeguardHost)
	log.Debug("loopdHosts: ", loopdHosts)
	log.Debug("loopdMacaroons: ", loopdMacaroons)
	log.Debug("loopdTLSCerts: ", loopdTLSCerts)

}
