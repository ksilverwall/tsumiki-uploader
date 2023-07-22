package repositories

import (
	"catch-all/models"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type StateMachine struct {
	Client          *sfn.SFN
	StateMachineArn string
}

func (s StateMachine) Execute(input models.ThumbnailRequest) error {
	b, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed marchial input: %w", err)
	}

	_, err = s.Client.StartExecution(&sfn.StartExecutionInput{
		StateMachineArn: aws.String(s.StateMachineArn),
		Input:           aws.String(string(b)),
	})
	if err != nil {
		return fmt.Errorf("failed to execute state machine: %w", err)
	}

	return nil
}
