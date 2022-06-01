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
	"github.com/spf13/cobra"
	"os"
)

var (
	confFile string
	dataFile string
	op       string
	endRow   int
	startRow int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tkeelBatchTool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(onInitialize)
	rootCmd.PersistentFlags().StringVarP(&confFile, "conf", "c", "", "The iot api config")
	rootCmd.PersistentFlags().StringVarP(&dataFile, "file", "f", "", "The data excel")
	rootCmd.PersistentFlags().StringVarP(&op, "op", "o", "", "add or del")
    //rootCmd.PersistentFlags().IntVarP(&endRow, "end_row", "e", 0, "op end row")
	//rootCmd.PersistentFlags().IntVarP(&startRow, "start_row", "s", 0, "op start row")
}

// initConfig reads in config file and ENV variables if set.
func onInitialize() {
}
