package main

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/jahkeup/repost/config"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

func ReadConfig(path string) (*config.Config, error) {
	logrus.Debugf("Loading configuration from %q", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot load config from %q", path)
	}

	var conf config.Config
	logrus.Debugf("Config loaded from %q: %v", conf)
	toml.Unmarshal(data, &conf)
	return &conf, nil
}
