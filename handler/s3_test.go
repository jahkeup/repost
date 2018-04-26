package handler

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/jahkeup/repost/delivery"
	noti "github.com/jahkeup/repost/notification"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testNotification = `{
  "notificationType": "Received",
  "mail": {
    "timestamp": "2018-04-20T02:38:36.148Z",
    "source": "you@sender-domain.com",
    "messageId": "st81p82fvjpm324a7enhdvnhfufc1bsedpfspsg1",
    "destination": [
      "test@subdomain.example.net"
    ],
    "headersTruncated": false,
    "headers": [
      {
        "name": "Return-Path",
        "value": "<you@sender-domain.com>"
      },
      {
        "name": "Received",
        "value": "from mail-oi0-f53.google.com (mail-oi0-f53.google.com [209.85.218.53]) by inbound-smtp.us-west-2.amazonaws.com with SMTP id st81p82fvjpm324a7enhdvnhfufc1bsedpfspsg1 for test@subdomain.example.net; Fri, 20 Apr 2018 02:38:36 +0000 (UTC)"
      },
      {
        "name": "X-SES-Spam-Verdict",
        "value": "PASS"
      },
      {
        "name": "X-SES-Virus-Verdict",
        "value": "PASS"
      },
      {
        "name": "Received-SPF",
        "value": "pass (spfCheck: domain of sender-domain.com designates 209.85.218.53 as permitted sender) client-ip=209.85.218.53; envelope-from=you@sender-domain.com; helo=mail-oi0-f53.google.com;"
      },
      {
        "name": "Authentication-Results",
        "value": "amazonses.com; spf=pass (spfCheck: domain of sender-domain.com designates 209.85.218.53 as permitted sender) client-ip=209.85.218.53; envelope-from=you@sender-domain.com; helo=mail-oi0-f53.google.com; dkim=pass header.i=@sender-domain.com;"
      },
      {
        "name": "X-SES-RECEIPT",
        "value": "AEFBQUFBQUFBQUFHdXZITnptZ3JtYTNYc0RuSzFEOHBqMEJmOS8xU0RqQWE2bjFrQ0cxUTZpMHByY0Vrb1IxQjd0UjZud0Z4dEdnekdRNk1LbEJqQ1I1UnEwSmdGK2ZyelAzYXB3K2hCcERwZ2lib1NxdTlYTUNqMy9HKzZMZm5lZVRjL2t0MGJEN0lEQXVjc1pBampXVnBYeWpMYk9sNjRzMFQvV3ZxWlRvdFY2TVBmTG9XanlJRXdtTitIZUVxTVMzV3BkUzRnVDhjUElNaVNsYVBEQjdHaVlrMWZ4V2hNanB4NndoQjlhdEVsubdomaini84SCtGY3U2V0taN1FSYkhQc2lJN2o4bVFmNHZnaEM1SHkwUE5nRVc5YU5tSlZocC9mSXFMUmEzbXhGZXdQTTlvVmFlcDRBb2c9PQ=="
      },
      {
        "name": "X-SES-DKIM-SIGNATURE",
        "value": "a=rsa-sha256; q=dns/txt; b=MPv5JBZs1ntdvr2xFa6MqqYrl6Xp+oYunzEAu+i1Swu60p/J3k7JDMTpLr+Md3TBqW2Vgw1ukR5zyGHQRsfSubdomainBtkDHdL3mH43yuxg/FJg+m3LgshsmQyeyeiPZkTHwunZBG4cLuKcHZsyemBscLv1MyMg1T+HpTt0GwCeVnevVU=; c=relaxed/simple; s=7v7vs6w47njt4pimodk5mmttbegzsi6n; d=amazonses.com; t=1524191916; v=1; bh=3pdvB7ME7i0I+vMJ8epoY1Ag89sdn1KnNxvpA/J4LAU=; h=From:To:Cc:Bcc:Subject:Date:Message-ID:MIME-Version:Content-Type:X-SES-RECEIPT;"
      },
      {
        "name": "Received",
        "value": "by mail-oi0-f53.google.com with SMTP id p62-v6so6730848oie.10 for <test@subdomain.example.net>; Thu, 19 Apr 2018 19:38:35 -0700 (PDT)"
      },
      {
        "name": "DKIM-Signature",
        "value": "v=1; a=rsa-sha256; c=relaxed/relaxed; d=sender-domain.com; s=google; h=mime-version:from:date:message-id:subject:to; bh=Jl+R73h+5DXursl3A0XEOzalYkAMCHCLWEcyITHqOBE=; b=F+l0uAD0n69yF4f+7RgR5VsOaJwdiVpg7E2TnXZenwxqKQjLBp09VVOt/Rw9ZeiZ5IVH7haLnrqN2t/N/Q4rcG225TgZRn4yMr9wJx8A70zKx3NQYtV3Xf7NtD+Q9IL6cXrXp3I0EHuf3UL8f6LsZW5CJDkHB1dUnvoUgIeZ/RI="
      },
      {
        "name": "X-Google-DKIM-Signature",
        "value": "v=1; a=rsa-sha256; c=relaxed/relaxed; d=1e100.net; s=20161025; h=x-gm-message-state:mime-version:from:date:message-id:subject:to; bh=Jl+R73h+5DXursl3A0XEOzalYkAMCHCLWEcyITHqOBE=; b=ClSLw6kvhh2/Bmya/iRueggXCRAqKQL/TGld/NmqXsXH3SfQ6t3wW/h2xepRIadGdq f/OHNH+1TKZe/ZQhySw5DmwsPJRHSbJOPMwd5RRPJHANHqAWIk7hJkJA3GmRHKS5alYv mj7yD4T7ebjXtGqlXRM/Lh/8ThexAcJEtFhYxzWvnezwykeWlhPw2GRl8okI0Xty9neJ QctAWo//S8DTjZ5YYP5vpJUfo+WVIh1YfIBnbCQztoT5rAQzyHmrr4vokiL9cNoy2kXd 66CY6/j0guR0bZup6AnTtd7trwUEF1v04/DIgqbKW88YwWXVOmFqrBbqChT1hfAefn+j 6LHA=="
      },
      {
        "name": "X-Gm-Message-State",
        "value": "ALQs6tBuZsJZftzu7sZLjBVETjGJTf7qReG15EFzIJfPYkvFrbD1aVdm frFJfFm+GL+e04lq3TSWEADVglE+2FN5pGFc8mfb/VN+1gQ="
      },
      {
        "name": "X-Google-Smtp-Source",
        "value": "AIpwx4/SlXVEDzlBICwbVrwAt16sPrHUwmaUuX2QajmsfuGsPj5UteOFjLPLNFSBQCJpPMWY1TdofW7TcHaqnON5VV4="
      },
      {
        "name": "X-Received",
        "value": "by 2002:aca:4e15:: with SMTP id c21-v6mr4155961oib.254.1524191914929; Thu, 19 Apr 2018 19:38:34 -0700 (PDT)"
      },
      {
        "name": "MIME-Version",
        "value": "1.0"
      },
      {
        "name": "Received",
        "value": "by 10.201.41.44 with HTTP; Thu, 19 Apr 2018 19:38:34 -0700 (PDT)"
      },
      {
        "name": "X-Originating-IP",
        "value": "[73.140.244.167]"
      },
      {
        "name": "From",
        "value": "Jacob Vallejo <you@sender-domain.com>"
      },
      {
        "name": "Date",
        "value": "Thu, 19 Apr 2018 19:38:34 -0700"
      },
      {
        "name": "Message-ID",
        "value": "<CAP1tTkNSfK6zSTanzGBRBQUdcDpCiLuQsUAwhQiMC2Bjn5mSqg@mail.gmail.com>"
      },
      {
        "name": "Subject",
        "value": "check check"
        },
      {
        "name": "To",
        "value": "test@subdomain.example.net"
      },
      {
        "name": "Content-Type",
        "value": "multipasubdomain/alternative; boundary=\"000000000000a8c537056a3e9555\""
      },
      {
        "name": "X-Example-Received",
        "value": "SES"
      }
    ],
    "commonHeaders": {
      "returnPath": "you@sender-domain.com",
      "from": [
        "Jacob Vallejo <you@sender-domain.com>"
      ],
      "date": "Thu, 19 Apr 2018 19:38:34 -0700",
      "to": [
        "test@subdomain.example.net"
      ],
      "messageId": "<CAP1tTkNSfK6zSTanzGBRBQUdcDpCiLuQsUAwhQiMC2Bjn5mSqg@mail.gmail.com>",
      "subject": "check check"
    }
  },
  "receipt": {
    "timestamp": "2018-04-20T02:38:36.148Z",
    "processingTimeMillis": 511,
    "recipients": [
      "test@subdomain.example.net"
    ],
    "spamVerdict": {
      "status": "PASS"
    },
    "virusVerdict": {
      "status": "PASS"
    },
    "spfVerdict": {
      "status": "PASS"
    },
    "dkimVerdict": {
      "status": "PASS"
    },
    "dmarcVerdict": {
      "status": "GRAY"
    },
    "action": {
      "type": "S3",
      "topicArn": "arn:aws:sns:us-west-2:11111111111:subdomain-example-net-inbound",
      "bucketName": "subdomain-example-net-ses-inbound",
      "objectKeyPrefix": "mailbox/",
      "objectKey": "mailbox/st81p82fvjpm324a7enhdvnhfufc1bsedpfspsg1"
    }
  }
}`

func TestDeliveryBucketObject(t *testing.T) {
	var deliveryNotif noti.DeliveryNotification
	err := json.Unmarshal([]byte(testNotification), &deliveryNotif)
	require.NoError(t, err)

	bucket, objectKey, err := deliveryBucketObject(deliveryNotif)
	assert.NoError(t, err, "should be able to determine s3 object from notification")
	assert.NotEmpty(t, bucket)
	assert.NotEmpty(t, objectKey)
}

type closableBufferReader struct {
	buffer *bytes.Buffer
	closed bool
}

func newClosableBufferReader(buf []byte) *closableBufferReader {
	return &closableBufferReader{
		buffer: bytes.NewBuffer(buf),
		closed: false,
	}
}

func (cb *closableBufferReader) Read(p []byte) (int, error) {
	if cb.closed {
		return 0, errors.New("buffer is closed")
	}
	return cb.buffer.Read(p)
}

func (cb *closableBufferReader) Close() error {
	if cb.closed {
		return errors.New("already closed")
	}
	cb.closed = true
	return nil
}

func TestS3Handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3client := NewMockS3Client(ctrl)
	in := []byte("this is some data")

	// Record and assert requested object
	var (
		reqObject string
		reqBucket string
	)
	gomock.InOrder(
		// Get should return some object data
		s3client.EXPECT().GetObject(gomock.Any()).Times(1).DoAndReturn(
			func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				reqObject = aws.StringValue(input.Key)
				reqBucket = aws.StringValue(input.Bucket)
				require.NotEmpty(t, reqObject)
				require.NotEmpty(t, reqBucket)

				return &s3.GetObjectOutput{
					Body: newClosableBufferReader(in),
				}, nil
			},
		),

		// Delete should follow after being handled.
		s3client.EXPECT().DeleteObject(gomock.Any()).Times(1).DoAndReturn(
			func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
				assert.Equal(t, reqObject, aws.StringValue(input.Key))
				assert.Equal(t, reqBucket, aws.StringValue(input.Bucket))
				return &s3.DeleteObjectOutput{}, nil
			},
		),
	)

	// Run Handler
	var deliveryNotif noti.DeliveryNotification
	err := json.Unmarshal([]byte(testNotification), &deliveryNotif)
	require.NoError(t, err)

	capture := delivery.NewCapture()

	// vender will only be used for one vend.
	vender := NewFuncVender(func() delivery.Deliverer {
		return capture
	})

	s3handler := NewS3(s3client, vender)
	s3handler.log = s3handler.log.WithField("test", t.Name())

	err = s3handler.HandleDelivery(deliveryNotif)
	assert.NoError(t, err)
	assert.Equal(t, in, capture.Data())
}

func TestS3HandlerKeepMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s3client := NewMockS3Client(ctrl)
	in := []byte("this is some data")

	// Get should return some object data, but should NOT call delete.
	s3client.EXPECT().GetObject(gomock.Any()).Times(1).DoAndReturn(
		func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: newClosableBufferReader(in),
			}, nil
		},
	)

	// Run Handler
	var deliveryNotif noti.DeliveryNotification
	err := json.Unmarshal([]byte(testNotification), &deliveryNotif)
	require.NoError(t, err)

	capture := delivery.NewCapture()

	// vender will only be used for one vend.
	vender := NewFuncVender(func() delivery.Deliverer {
		return capture
	})

	s3handler := NewS3(s3client, vender).KeepMessages(true)
	s3handler.log = s3handler.log.WithField("test", t.Name())

	err = s3handler.HandleDelivery(deliveryNotif)
	assert.NoError(t, err)
	assert.Equal(t, in, capture.Data())
}
