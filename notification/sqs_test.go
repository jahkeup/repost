package notification

import (
	"context"
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

func TestDeliveryHandlerInvoke(t *testing.T) {
	maxTestDuration := 1 * time.Second
	testCtx, cancel := context.WithTimeout(context.Background(), maxTestDuration)
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queueURL := "https://test-queue-url.example.com"

	receiver := NewMockSQSReceiver(ctrl)
	handler := NewMockDeliveryHandler(ctrl)

	receiver.EXPECT().ReceiveMessage(gomock.Any()).MinTimes(1).DoAndReturn(
		// Handle request for messages off of a queue
		func(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
			require.NotNil(t, input)
			assert.Equal(t, queueURL, aws.StringValue(input.QueueUrl))

			response := &sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					&sqs.Message{
						Body: aws.String(""),
					},
				},
			}

			return response, nil
		},
	)

	handler.EXPECT().HandleDelivery(gomock.Any()).MinTimes(1)

	notif := NewSQS(receiver, queueURL)
	notif.SetDeliveryHandler(handler)
	err := notif.Poll(testCtx)
	assert.NoError(t, err)
}
