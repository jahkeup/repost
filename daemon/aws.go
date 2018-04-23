package daemon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func (d *Daemon) logCallerIdentity(session *session.Session) error {
	c := sts.New(session)
	out, err := c.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	d.log.Infof("aws credentials EUID: %q", aws.StringValue(out.Arn))
	return nil
}
