package main

import (
	"flag"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jahkeup/repost/config"
	"github.com/jahkeup/repost/daemon"
)

var (
	flagGenerateTOML = flag.Bool("generate-config", false, "Generate and print example configuration TOML")
)

func init() {
	flag.Parse()
}

func main() {
	if *flagGenerateTOML {
		toml := generateTOML()
		fmt.Printf("%s", toml)
		return
	}

	conf, err := ReadConfig("./repost.toml")
	if err != nil {
		logrus.Fatal(err)
	}
	run(conf)
}

func run(config *config.Config) error {
	mainCtx, cancel := Context()
	daemon, err := daemon.New(mainCtx, config)
	if err != nil {
		return err
	}
	err = daemon.Run(mainCtx)
	cancel()
	if err != nil {
		logrus.Errorf("daemon encountered an error: %s", err)
		return err
	}
	logrus.Info("daemon is exiting")
	return nil
}
