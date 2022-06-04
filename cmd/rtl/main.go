package main

import (
	"github.com/rmrfslashbin/aws-cf-rtl/subcmds"
	"github.com/sirupsen/logrus"
)

func main() {
	// Catch errors
	var err error
	defer func() {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("main crashed")
		}
	}()
	if err := subcmds.Execute(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("abend")
	}
}
