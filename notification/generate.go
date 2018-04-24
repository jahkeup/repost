//go:generate mockgen -package notification -destination sqs_receiver_mock_test.go github.com/jahkeup/repost/notification SQSReceiver
//go:generate mockgen -package notification -destination delivery_handler_mock_test.go github.com/jahkeup/repost/notification DeliveryHandler
package notification
