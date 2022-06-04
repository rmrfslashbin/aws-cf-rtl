/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package subcmds

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCrawlerStatus represents the status command
var getCrawlerStatus = &cobra.Command{
	Use:   "status",
	Short: "Status the crawler",
	RunE: func(cmd *cobra.Command, args []string) error {
		if crawlerName == "" {
			fmt.Println("crawler name is required")
			return fmt.Errorf("crawler name not specified")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running crawler %s\n", crawlerName)
	},
}

func init() {
	logrus.Info("status init")
	crawlerCmd.AddCommand(getCrawlerStatus)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
