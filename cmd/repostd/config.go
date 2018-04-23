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

func generateTOML() []byte {
	data, err := toml.Marshal(exampleConfig())
	if err != nil {
		panic(err)
	}
	return data
}

func exampleConfig() config.Config {
	conf := config.Config{
		General: config.General{
			LogLevel: "INFO",
			Region:   "us-west-2",
			Profile:  "repost",
		},
		Notification: config.Notification{
			QueueURL: "https://queue-url.com",
		},
		Delivery: config.Delivery{
			Pipe: config.PipeDelivery{
				Command: "tee ./delivery/{{.MessageId}}",
			},
		},
	}
	return conf
}
