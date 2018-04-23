package config

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

// General configuration options for program run.
type General struct {
	apply sync.Once

	LogLevel string

	// AWS Region
	Region string
	// AWS Credential Profile
	Profile string
	// AWS Assume Role ARN
	RoleArn string
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

func (g *General) awsConfig() aws.Config {
	return aws.Config{
		Region: aws.String(g.Region),
	}
}

func (g *General) session() (*session.Session, error) {
	conf := g.awsConfig()
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  conf,
		Profile: g.Profile,
	})

	if g.RoleArn != "" {
		return g.assumeSession(sess)
	}

	return sess, err
}

func (g *General) assumeSession(initSess *session.Session) (*session.Session, error) {

	arn, err := arn.Parse(g.RoleArn)
	if err != nil {
		return nil, errors.Wrap(err, "RoleArn could not be parsed")
	}

	baseConf := g.awsConfig()
	stscreds := stscreds.NewCredentials(initSess, arn.String())
	sess, err := session.NewSession(&baseConf, &aws.Config{
		Credentials: stscreds,
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
