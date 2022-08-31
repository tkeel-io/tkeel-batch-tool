/*
Copyright 2021 The tKeel Authors.

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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/dapr/cli/pkg/print"
	"github.com/spf13/cobra"
	"os"
	"tkeelBatchTool/src/conf"
	"tkeelBatchTool/src/http"
	mock "tkeelBatchTool/src/mock_device"
)

var (
	deviceID     string
	deviceToken  string
	template     string
	templateMode string
	interval     int64
	host         string
	broker       string
)

// templateCreateCmd represents the templateCreate command
// go run main.go  mock --broker preview.tkeel.io:31883 --device iotd-d33e5f4e-2d99-4aa3-9a1b-45af10bb9c0d --token Y2M3YTk2NzMtNmZkYS0zMDU0LWEwMDEtNzE3ZjBiNzNhOWQ5 --template ./telemetry.mock
var mockSendDataCmd = &cobra.Command{
	Use:   "mock",
	Short: "Mock device",
	Long:  `Mock device.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mock device SendData")

		//config
		if err := conf.InitConfig("./config.json"); err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}
		if broker == "" {
			prompt := &survey.Input{Message: "Please enter broker: ",Default: conf.DefaultConfig.Broker}
			if err := survey.AskOne(prompt, &broker); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read broker from stdin")
				return
			}
		}
		if deviceID == "" {
			prompt := &survey.Input{Message: "Please enter deviceID: ", Default: conf.DefaultConfig.DeviceID}
			if err := survey.AskOne(prompt, &deviceID); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read deviceID from stdin")
				return
			}
		}
		if deviceToken == "" {
			prompt := &survey.Input{Message: "Please enter device Token: ",Default: conf.DefaultConfig.DeviceToken}
			if err := survey.AskOne(prompt, &deviceToken); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read DeviceToken from stdin")
				return
			}
		}

		if templateMode == "" {
			prompt := &survey.Select{Message: "Please enter Mode: ", Options: []string{
				"telemetry", "attributes",
			}}
			if err := survey.AskOne(prompt, &templateMode); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read deviceID from stdin")
				return
			}
		}
		if template == "" {
			prompt := &survey.Input{Message: "Please enter template file: ",Default: fmt.Sprintf("%s.mock",templateMode)}
			if err := survey.AskOne(prompt, &template); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read template file from stdin")
				return
			}
		}

		if checkHostAndDeviceID() {
			cmd.UsageString()
			return
		}


		tpl, err := http.DeviceDataTemplate(template)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Read Device Data Template(%s) fail:%s", err.Error())
			return
		}


		conf.DefaultConfig.Template = template
		conf.DefaultConfig.TemplateMode = templateMode
		conf.DefaultConfig.DeviceToken = deviceToken
		conf.DefaultConfig.Broker = broker
		conf.DefaultConfig.DeviceID = deviceID


		//config
		if err := conf.SaveConfig("./config.json"); err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}

		err = mock.RunDataSender(mock.ClientOptions{
			Interval:    interval,
			Template:    tpl,
			Host:        broker,
			Mode:        templateMode,
			DeviceID:    deviceID,
			DeviceToken: deviceToken,
			//Host:,
		})
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Run Data Sender fail:%s", err.Error())
			return
		}

	},
}

// templateCreateCmd represents the templateCreate command
var mockTemplateCmd = &cobra.Command{
	Use:   "mock-sample",
	Short: "Creat mock sample template",
	Long:  `Creat mock sample template`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mock template creat")

		//config
		if err := conf.InitConfig("./config.json"); err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}


		if host == "" {
			prompt := &survey.Input{Message: "Please enter host: ", Default: conf.DefaultConfig.Host}
			if err := survey.AskOne(prompt, &host); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read broker from stdin")
				return
			}
		}
		if deviceID == "" {
			prompt := &survey.Input{Message: "Please enter deviceID: ", Default: conf.DefaultConfig.DeviceID}
			if err := survey.AskOne(prompt, &deviceID); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read deviceID from stdin")
				return
			}
		}

		if templateMode == "" {
			prompt := &survey.Select{Message: "Please enter Mode: ", Options: []string{
				"telemetry", "attributes",
			}}
			if err := survey.AskOne(prompt, &templateMode); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read deviceID from stdin")
				return
			}
		}

		if checkHostAndDeviceID() {
			cmd.UsageString()
			return
		}

		conf.DefaultConfig.Host = host
		conf.DefaultConfig.DeviceID = deviceID
		conf.DefaultConfig.DeviceToken = deviceToken
		conf.DefaultConfig.TemplateMode = templateMode
		//config
		if err := conf.SaveConfig("./config.json"); err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}

		ret, err := http.GenDeviceTemplate(host, deviceID, templateMode)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Gen Device mock sample file fail: %s", err)
			return
		}

		if template == "" {
			template = fmt.Sprintf("%s.mock", templateMode)
		}

		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, ret, "", "\t")
		if error != nil {
			print.FailureStatusEvent(os.Stdout, "JSON parse error: ", error)
			return
		}

		if file, err := os.Create(template); err != nil {
			print.FailureStatusEvent(os.Stdout, "error open file(%s)", template)
			return
		} else {
			defer file.Close()
			if count, err := file.Write(prettyJSON.Bytes()); err != nil {
				print.FailureStatusEvent(os.Stdout, "error close file(%s)", template)
				return
			} else {
				print.SuccessStatusEvent(os.Stdout, "Write mock sample file(%s), size %d byte", template, count)
			}
		}
	},
}

func checkHostAndDeviceID() bool {
	if deviceID == "" {
		print.FailureStatusEvent(os.Stdout, "deviceID unset, --device iotd-43b3acd1-de84-4a5b-afb1-b9b88f87708a")
		return true
	}

	if templateMode != "telemetry" && templateMode != "attributes" {
		print.FailureStatusEvent(os.Stdout, "template mode unset, --mode attributes or --mode telemetry")
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(mockSendDataCmd)
	mockSendDataCmd.PersistentFlags().StringVar(&templateMode, "mode", "", "The template config mode(attributes,telemetry)")
	mockSendDataCmd.PersistentFlags().Int64Var(&interval, "interval", 1000, "The interval of message, default is 1000(1 second)")
	mockSendDataCmd.PersistentFlags().StringVar(&broker, "broker", "", "The broker of tkeel server(like preview.tkeel.io:31883)")
	mockSendDataCmd.PersistentFlags().StringVar(&template, "template", "", "The template config filename")
	mockSendDataCmd.PersistentFlags().StringVar(&deviceID, "device", "", "The device ID")
	mockSendDataCmd.PersistentFlags().StringVar(&deviceToken, "token", "", "The device ID")
	rootCmd.AddCommand(mockTemplateCmd)
	mockTemplateCmd.PersistentFlags().StringVar(&templateMode, "mode", "", "The template config mode(attributes,telemetry)")
	mockTemplateCmd.PersistentFlags().StringVar(&host, "host", "", "The host of tkeel server(like http://preview.tkeel.io:30080/)")
	mockTemplateCmd.PersistentFlags().StringVar(&template, "template", "", "The template config filename")
	mockTemplateCmd.PersistentFlags().StringVar(&deviceID, "device", "", "The device ID")
}
