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
	"tkeelBatchTool/src/conf"
	"tkeelBatchTool/src/create"
	"tkeelBatchTool/src/del"
	"tkeelBatchTool/src/parse"
)

// spaceTreeRegisterCmd represents the spaceTreeRegister command
var spaceTreeCmd = &cobra.Command{
	Use:   "spaceTree",
	Short: "Creat spaceTree from excel",
	Long:  `Creat spaceTree from excel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spaceTree called")

		//config
		if err := conf.InitConfig("./config.json"); err != nil {
			panic(err)
		}

		//parse
		fmt.Println("start parse xlsx\n")
		content, f, err, order := parse.DoParseSpaceTreeExcel(dataFile, startRow, endRow)
		if err != nil {
			fmt.Println("parse xlsx err\n")
			//panic(err)
			return
		}

		if op == "del" {
			// del
			fmt.Println("start del \n")
			if err := del.DelSpaceTree(content,order); err != nil {
				//panic(err)
				return
			}
		} else {
			// create
			fmt.Println("start create \n")
			if err := create.CreateSpaceTree(content, f, order); err != nil {
				//panic(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(spaceTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spaceTreeRegisterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spaceTreeRegisterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
