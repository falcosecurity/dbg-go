package main

import (
	"github.com/fededp/dbg-go/cmd"
	logger "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.WithError(err).Fatal("error executing dbg-go")
	}
}
