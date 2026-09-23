/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http/cookiejar"

	"manabacli/src/config"
	"manabacli/src/logic"
	"manabacli/src/manabaclient"

	"github.com/spf13/cobra"
)

var getreportCmd = &cobra.Command{
	Use:   "getreport <URL>",
	Short: "<URL>にはコースのレポートページのURLを指定してください. レポートの名前と提出先URLの一覧を取得します.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username, password, err := config.Load()
		if err != nil {
			return err
		}

		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}

		client := manabaclient.Client{}
		err = client.Login(jar, username, password)
		if err != nil {
			return err
		}

		url := args[0]
		reportNameToUrl, err := logic.GetReportNameToUrl(jar, url)
		if err != nil {
			return err
		}

		for name, url := range reportNameToUrl {
			fmt.Printf("%s -> %s\n", name, url)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getreportCmd)
}
