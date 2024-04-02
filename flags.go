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
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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

	//Remove http:// from the endpoint if present
	if strings.HasPrefix(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "http://") {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", strings.Replace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "http://", "", 1))
	}

	//Remove https:// from the endpoint if present
	if strings.HasPrefix(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "https://") {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", strings.Replace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "https://", "", 1))
	}

	os.Setenv("OTLP_RESOURCES_ATTRIBUTES", os.ExpandEnv("$OTLP_RESOURCES_ATTRIBUTES"))

	viper.AutomaticEnv() // read in environment variables that match

	//Lndconnect uris
	rootCmd.Flags().String("lndconnecturis", "", "CSV of lndconnect strings to connect to lnd(s)")
	viper.BindPFlag("lndconnecturis", rootCmd.Flags().Lookup("lndconnecturis"))

	//Loopdconnect uris
	rootCmd.Flags().String("loopdconnecturis", "", "CSV of loopdconnect strings to connect to loopd(s)")
	viper.BindPFlag("loopdconnecturis", rootCmd.Flags().Lookup("loopdconnecturis"))

	rootCmd.Flags().String("pollingInterval", "15s", "Interval to poll data")
	viper.BindPFlag("pollingInterval", rootCmd.Flags().Lookup("pollingInterval"))

	rootCmd.Flags().String("logLevel", "info", "Log level from values: {trace, debug, info, warn, error, fatal, panic}")
	viper.BindPFlag("logLevel", rootCmd.Flags().Lookup("logLevel"))

	rootCmd.Flags().String("logFormat", "text", "Log format from: {text, json}")
	viper.BindPFlag("logFormat", rootCmd.Flags().Lookup("logFormat"))

	//Flags for nodeguard grpc endpoint
	rootCmd.Flags().String("nodeguardHost", "", "Hostname:port to connect to nodeguard")
	viper.BindPFlag("nodeguardHost", rootCmd.Flags().Lookup("nodeguardHost"))

	//Swap Publication Offset in minutes
	rootCmd.Flags().String("swapPublicationOffset", "60m", "Swap publication deadline offset (Maximum time for the swap provider to publish the swap)")
	viper.BindPFlag("swapPublicationOffset", rootCmd.Flags().Lookup("swapPublicationOffset"))

	// Retries before applying backoff to the swap
	rootCmd.Flags().Int("retriesBeforeBackoff", 3, "Number of retries before applying backoff to the swap")
	viper.BindPFlag("retriesBeforeBackoff", rootCmd.Flags().Lookup("retriesBeforeBackoff"))

	// Coefficient to apply to the backoff
	rootCmd.Flags().Float64("backoffCoefficient", 0.95, "Coefficient to apply to the backoff")
	viper.BindPFlag("backoffCoefficient", rootCmd.Flags().Lookup("backoffCoefficient"))

	// Limit coefficient of the backoff
	rootCmd.Flags().Float64("backoffLimit", 0.1, "Limit coefficient of the backoff")
	viper.BindPFlag("backoffLimit", rootCmd.Flags().Lookup("backoffLimit"))

	// Limit fees for swaps
	rootCmd.Flags().Float64("limitFees", 0.007, "Limit fees for swaps e.g. 0.01 = 1%")
	viper.BindPFlag("limitFees", rootCmd.Flags().Lookup("limitFees"))

	//Sweep conf
	rootCmd.Flags().String("sweepConfTarget", "400", "Target number of confirmations for swaps, this uses bitcoin core broken estimator, procced with caution")
	viper.BindPFlag("sweepConfTarget", rootCmd.Flags().Lookup("sweepConfTarget"))

	//Sleep between retries
	rootCmd.Flags().Duration("sleepBetweenRetries", 500*time.Millisecond, "Sleep between retries")
	viper.BindPFlag("sleepBetweenRetries", rootCmd.Flags().Lookup("sleepBetweenRetries"))

	//Sleep between retries backoff
	rootCmd.Flags().Float64("sleepBetweenRetriesBackoff", 1.5, "Sleep between retries backoff")
	viper.BindPFlag("sleepBetweenRetriesBackoff", rootCmd.Flags().Lookup("sleepBetweenRetriesBackoff"))

	//Sleep max
	rootCmd.Flags().Duration("sleepMax", 60*time.Second, "Sleep max")
	viper.BindPFlag("sleepMax", rootCmd.Flags().Lookup("sleepMax"))

	//Now we set the global vars

	pollingInterval = viper.GetDuration("pollingInterval")
	nodeguardHost = viper.GetString("nodeguardHost")
	loopdconnectURIs = strings.Split(viper.GetString("loopdconnecturis"), ",")
	lndconnectURIs = strings.Split(viper.GetString("lndconnecturis"), ",")

	retries = viper.GetInt("retriesBeforeBackoff")
	backoffCoefficient = viper.GetFloat64("backoffCoefficient")
	backoffLimit = viper.GetFloat64("backoffLimit")
	limitFees = viper.GetFloat64("limitFees")

	sleepBetweenRetries = viper.GetDuration("sleepBetweenRetries")
	sleepBetweenRetriesBackoff = viper.GetFloat64("sleepBetweenRetriesBackoff")
	sleepMax = viper.GetDuration("sleepMax")

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

	log.Debug("pollingInterval: ", pollingInterval)
	log.Debug("logLevel: ", logLevel)
	log.Debug("logFormat: ", viper.GetString("logFormat"))
	log.Debug("nodeguardHost: ", nodeguardHost)
	log.Debug("retriesBeforeBackoff: ", viper.GetInt("retriesBeforeBackoff"))
	log.Debug("backoffCoefficient: ", viper.GetFloat64("backoffCoefficient"))
	log.Debug("backoffLimit: ", viper.GetFloat64("backoffLimit"))
	log.Debug("limitFees: ", viper.GetFloat64("limitFees"))
	log.Debug("sleepBetweenRetries: ", viper.GetDuration("sleepBetweenRetries"))
	log.Debug("sleepBetweenRetriesBackoff: ", viper.GetFloat64("sleepBetweenRetriesBackoff"))
	log.Debug("sleepMax: ", viper.GetDuration("sleepMax"))
}
