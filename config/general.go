package config

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// General configuration options for program run.
type General struct {
	apply sync.Once

	LogLevel string

	// AWS Region
	Region string
	// AWS Credential Profile
	Profile string
}

func (g *General) Apply() {
	g.apply.Do(func() {
		lvl, err := logrus.ParseLevel(g.LogLevel)
		if err != nil {
			logrus.SetLevel(logrus.InfoLevel)
			logrus.Warnf("Provided loglevel %q was invalid, falling back to INFO", g.LogLevel)
		} else {
			logrus.SetLevel(lvl)
		}
	})
}

func (c *Config) Session() (*session.Session, error) {
	return c.General.session()
}

func (g *General) session() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(g.Region),
		},
		Profile: g.Profile,
	})
	return sess, err
}
