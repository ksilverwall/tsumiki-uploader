package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ParameterRepository struct {
	Client *ssm.SSM
}

func (r ParameterRepository) Get(key string, val any) error {
	params, err := r.Client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to load parameter %s: %w", key, err)
	}

	s := params.Parameter.Value

	err = json.Unmarshal([]byte(*s), val)
	if err != nil {
		return fmt.Errorf("failed to parse parameter '%s': %w", *s, err)
	}

	return nil
}
