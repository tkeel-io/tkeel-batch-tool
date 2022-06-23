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
	"io/fs"
	"os"
	path_filepath "path/filepath"
	"tkeelBatchTool/excel_file"
)

var ()

// templateCreateCmd represents the templateCreate command
var downloadExcelCmd = &cobra.Command{
	Use:   "download",
	Short: "Download excel",
	Long:  `Download excel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Download excel called")

		if outputPath == "" {
			outputPath = "."
		}

		if list, err := fs.Glob(excel_file.ExcelFiles, "*"); err == nil {
			os.Mkdir(outputPath, 0755)
			for _, filename := range list {
				if byt, err := fs.ReadFile(excel_file.ExcelFiles, filename); err == nil {
					writeFile(path_filepath.Join(outputPath, filename), byt)
				} else {
					fmt.Printf("Error when read file %v: %v\n", filename, err)
				}
			}
		} else {
			fmt.Printf("Error list file %v\n", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(downloadExcelCmd)
}

func writeFile(filePath string, byt []byte) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	defer file.Close()
	file.Write(byt)
	return nil
}
