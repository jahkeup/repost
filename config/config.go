package config

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
)

// Config options
type Config struct {
	General      General
	Notification Notification
	Delivery     Delivery
}

func (c *Config) Apply() error {
	c.General.apply()
	logrus.Debugf("Loaded configuration:\n%s", c.String())
	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}
