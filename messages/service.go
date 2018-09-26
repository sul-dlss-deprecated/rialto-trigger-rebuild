package messages

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
)

// MessageService is an interface for sending messages to the derivative service
type MessageService interface {
	Publish(subjects []string) error
}

// SNSMessageService is a message publishing service for SNS
type SNSMessageService struct {
	conn     *sns.SNS
	topicArn *string
}

// NewSNSMessageService creates a new instance of the message service
func NewSNSMessageService(conn *sns.SNS, topicArn *string) MessageService {
	return &SNSMessageService{
		conn:     conn,
		topicArn: topicArn,
	}
}

// Publish crafts a "touch" SNS message with the given subjects and
// pushes the message to the topic
func (s *SNSMessageService) Publish(subjects []string) error {
	msg := message.Message{
		Action:   "touch",
		Entities: subjects,
	}
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	str := string(json)

	input := &sns.PublishInput{
		Message:  &str,
		TopicArn: s.topicArn,
	}
	_, err = s.conn.Publish(input)
	return err
}
