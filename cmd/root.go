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
package cmd

import (
	"log"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "liquidator",
	Short: "A CLI tool to monitor and automate the liquidity of your LND channels",
	Run: func(cmd *cobra.Command, args []string) { 
		ErrorLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
		InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
		DebugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile)
		WarnLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)

		InfoLog.Println("Starting liquidator")
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.


	//OTEL Expanded vars
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", os.ExpandEnv("$OTEL_EXPORTER_OTLP_ENDPOINT"))
	os.Setenv("OTLP_RESOURCES_ATTRIBUTES", os.ExpandEnv("$OTLP_RESOURCES_ATTRIBUTES"))
	
	viper.AutomaticEnv() // read in environment variables that match

	rootCmd.Flags().StringSlice("nodesHosts", []string{}, "Command separated list of hostname:port to connect to")
	viper.BindPFlag("nodesHosts", rootCmd.Flags().Lookup("nodesHosts"))

	rootCmd.Flags().StringSlice("nodesMacaroons", []string{}, "Command separated list of macaroons used in nodesHosts, in the same order")
	viper.BindPFlag("nodesMacaroons", rootCmd.Flags().Lookup("nodesMacaroons"))

	rootCmd.Flags().Int64("pollingInterval", 10, "Interval to poll data in seconds")
	viper.BindPFlag("pollingInterval", rootCmd.Flags().Lookup("pollingInterval"))

	//If nodesHosts length is different than nodesMacaroons, exit
	if len(viper.GetStringSlice("nodesHosts")) != len(viper.GetStringSlice("nodesMacaroons")) {
		ErrorLog.Fatal("nodesHosts and nodesMacaroons must have the same length")
	}
	//Now we set the global vars
	nodesHosts = viper.GetStringSlice("nodesHosts")
	nodesMacaroons = viper.GetStringSlice("nodesMacaroons")
	pollingInterval = viper.GetInt64("pollingInterval")

}


// initConfig reads in config file and ENV variables if set.
func initConfig() {


}
