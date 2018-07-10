package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Handler is the Lambda function handler
func Handler(ctx context.Context) error {
	message := fmt.Sprintf("{\"action\": %s, \"message\": %s}", os.Getenv("REBUILD_ACTION"), os.Getenv("REBUILD_MESSAGE"))
	topicArn := os.Getenv("RIALTO_TOPIC_ARN")
	endpoint := os.Getenv("RIALTO_SNS_ENDPOINT")
	snsConn := sns.New(session.New(), aws.NewConfig().
		WithDisableSSL(false).
		WithEndpoint(endpoint))
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: &topicArn,
	}
	_, err := snsConn.Publish(input)
	if err != nil {
		log.Printf("Error publishing rebuild message to %v: %v", topicArn, err)
	}
	return err
}

func main() {
	lambda.Start(Handler)
}
