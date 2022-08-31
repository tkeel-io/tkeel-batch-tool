package cmd

import (
	"fmt"
	"os"
	"tkeelBatchTool/src/conf"
	"tkeelBatchTool/src/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dapr/cli/pkg/print"
	"github.com/spf13/cobra"
)

var (
	tenant    string
	username  string
	password  string
	printable bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login with password",
	Example: `
 tkeelBatchTool login http://tkeel.io:30080/ --tenant <your tenant name> --username <your username> --password <your password>
 - http://tkeel.io:30080/ is tKeel platform endpoint.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		//config
		if err := conf.InitConfig("./config.json"); err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}

		host := args[0]
		if host == "" {
			prompt := &survey.Input{Message: "Please enter host: ", Default: conf.DefaultConfig.Host}
			if err := survey.AskOne(prompt, &host); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read broker from stdin")
				return
			}
		}
		if tenant == "" {
			prompt := &survey.Input{Message: "Please enter your tenant: ", Default: conf.DefaultConfig.Tenant}
			if err := survey.AskOne(prompt, &tenant); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read tenant from stdin")
				return
			}
		}
		tenantID, err := http.GetTenantID(host, tenant)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "failed to get tenant id")
			return
		}

		if username == "" {
			prompt := &survey.Input{Message: "Please enter your username: ", Default: conf.DefaultConfig.Username}
			if err := survey.AskOne(prompt, &username); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read username from stdin")
				return
			}
		}
		if password == "" {
			prompt := &survey.Input{Message: "Please enter your password: ", Default: conf.DefaultConfig.Password}
			if err := survey.AskOne(prompt, &password); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read password from stdin")
				return
			}
		}
		accessToken, refreshToken, err := http.GetTenantLoginToken(host, tenantID, username, password)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}

		conf.DefaultConfig.IotUrl = fmt.Sprintf("%s/apis/tkeel-device", host)
		conf.DefaultConfig.Token = accessToken
		conf.DefaultConfig.RefreshToken = refreshToken
		conf.DefaultConfig.Host = host
		conf.DefaultConfig.Tenant = tenant
		conf.DefaultConfig.TenantID = tenantID
		conf.DefaultConfig.Username = username
		conf.DefaultConfig.Password = password
		//config
		if err := conf.SaveConfig("./config.json"); err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
		}

		print.SuccessStatusEvent(os.Stdout, "You are Login as %s in tenant %s!", username, tenant)
		print.SuccessStatusEvent(os.Stdout, "AccessToken is [%s]", accessToken)
		print.SuccessStatusEvent(os.Stdout, "RefreshToken is [%s]", refreshToken)
		print.SuccessStatusEvent(os.Stdout, "Login Token save in ./config.json!")
		print.SuccessStatusEvent(os.Stdout, "You are Login Success!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVar(&tenant, "tenant", "", "input your tenant")
	loginCmd.Flags().StringVar(&username, "username", "", "input your username")
	loginCmd.Flags().StringVar(&password, "password", "", "input your password")
}
