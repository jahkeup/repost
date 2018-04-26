package notification

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func testEnveloped(message string) string {
	d, err := json.Marshal(sqsBodyEnvelope{
		Type:      "notification",
		MessageId: "test-message-id",
		TopicArn:  "arn:aws:::",
		Subject:   "Test Message Received",
		Message:   message,
	})
	if err != nil {
		panic("envelopes stink")
	}
	return string(d)
}

func TestDeliveryHandlerInvoke(t *testing.T) {
	cases := []struct {
		name         string
		keepMessages bool
	}{
		{"KeepMessages", true},
		{"NoKeepMessages", false},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			maxTestDuration := 1 * time.Second
			testCtx, cancel := context.WithTimeout(context.Background(), maxTestDuration)
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			queueURL := "https://sqs.us-west-2.amazonaws.com/111111111111/test-queue"

			receiver := NewMockSQSReceiver(ctrl)
			handler := NewMockDeliveryHandler(ctrl)

			// Data to submit to notification poller
			sentDN := DeliveryNotification{}
			sentJSON, err := json.Marshal(sentDN)
			require.NoError(t, err)

			gomock.InOrder(
				receiver.EXPECT().ReceiveMessage(gomock.Any()).MinTimes(1).DoAndReturn(
					// Handle request for messages off of a queue
					func(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
						require.NotNil(t, input)
						assert.Equal(t, queueURL, aws.StringValue(input.QueueUrl))

						response := &sqs.ReceiveMessageOutput{
							Messages: []*sqs.Message{
								&sqs.Message{
									MessageId:     aws.String("test-message-id"),
									ReceiptHandle: aws.String("test-message-receipt-handle"),
									Body:          aws.String(testEnveloped(string(sentJSON))),
								},
							},
						}

						return response, nil
					},
				),
				handler.EXPECT().HandleDelivery(gomock.Any()).MinTimes(1),
			)

			if !tc.keepMessages {
				receiver.EXPECT().DeleteMessage(gomock.Any()).Times(1)
			}

			notif := NewSQS(receiver, queueURL).KeepMessages(tc.keepMessages)
			notif.SetDeliveryHandler(handler)
			err = notif.Poll(testCtx)
			require.NoError(t, err)
		})
	}
}
