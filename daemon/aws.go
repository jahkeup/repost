package daemon

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func (d *Daemon) logCallerIdentity(ctx context.Context, session *session.Session) error {
	d.log.Debug("Pinging STS with provided AWS credentials")
	c := sts.New(session)
	out, err := c.GetCallerIdentityWithContext(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	d.log.Infof("aws credentials EUID: %q", aws.StringValue(out.Arn))
	return nil
}
