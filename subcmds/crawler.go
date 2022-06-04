/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package subcmds

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// crawlerCmd represents the crawler command
var (
	// crawlerName is the name of the crawler to act on
	crawlerName string
	crawlerCmd  = &cobra.Command{
		Use:   "crawler",
		Short: "Crawler related commands",
	}
)

func init() {
	logrus.Info("crawler init")
	rootCmd.AddCommand(crawlerCmd)

	crawlerCmd.PersistentFlags().StringVar(&crawlerName, "cname", "", "Name of the crawler")
	//viper.BindPFlag("crawler.name", crawlerCmd.PersistentFlags().Lookup("cname"))
	logrus.WithFields(logrus.Fields{
		"crawlerName": crawlerName,
	}).Info("crawlerName")
}
