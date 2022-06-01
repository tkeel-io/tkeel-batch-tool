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

var ()

// templateCreateCmd represents the templateCreate command
var templateCreateCmd = &cobra.Command{
	Use:   "template",
	Short: "Creat template from excel",
	Long:  `Creat template from excel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("templateCreate called")

		//config
		if err := conf.InitConfig("./config.json"); err != nil {
			panic(err)
		}

		if dataFile == "" || dataFile == " " {
			fmt.Println("input data file is empty\n")
			return
		}

		//parse
		fmt.Println("start parse xlsx\n")
		content, _, err := parse.DoParseTemplateExcel(dataFile)
		if err != nil {
			fmt.Println("parse xlsx err\n")
			//panic(err)
			return
		}

		// create
		if op == "del" {
			fmt.Println("start del \n")
			if err := del.DelTemplate(content); err != nil {
				//panic(err)
				return
			}

		} else {
			fmt.Println("start create \n")
			if err := create.CreateTemplate(content); err != nil {
				//panic(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCreateCmd)
}
