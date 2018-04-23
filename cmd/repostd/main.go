package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/jahkeup/repost/config"
)

func main() {
	conf, err := ReadConfig("./repost.toml")
	if err != nil {
		logrus.Fatal(err)
	}
	run(conf)
}

func run(config *config.Config) error {
	//	ctx, cancel := Context()
	return nil
}
