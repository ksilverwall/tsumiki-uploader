package repositories

import (
	"catch-all/models"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Queue struct {
	SQS      *sqs.SQS
	QueueUrl string
}

func (q Queue) Push(req models.ThumbnailRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	sendMsgInput := sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonData)),
		QueueUrl:    &q.QueueUrl,
	}

	_, err = q.SQS.SendMessage(&sendMsgInput)
	if err != nil {
		return err
	}

	return nil
}
