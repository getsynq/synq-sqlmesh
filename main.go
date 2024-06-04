package main

import (
	"github.com/getsynq/synq-sqlmesh/cmd"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
)

//go:generate bash bin/version.sh

func main() {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "15:04:05",
		LogFormat:       "%time%  %msg%\n",
	})

	if err := cmd.Execute(); err != nil {
		logrus.WithError(err).Error("Error executing command")
		os.Exit(1)
	}
	os.Exit(0)
}
