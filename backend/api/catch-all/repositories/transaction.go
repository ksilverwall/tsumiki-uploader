package repositories

import (
	"catch-all/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type internalValue string

func (v internalValue) toModel() (models.Transaction, error) {
	var m models.Transaction
	if err := json.Unmarshal([]byte(v), &m); err != nil {
		return models.Transaction{}, fmt.Errorf("transaction cannnot decode from json value '%s': %w", v, err)
	}

	return m, nil
}

func newInternalValue(m models.Transaction) (internalValue, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return internalValue(""), fmt.Errorf("transaction cannnot encode to json: %w", err)
	}

	return internalValue(string(b)), nil
}

type DynamoDBItem struct {
	ID             models.TransactionID `json:"id"`
	Value          string               `json:"name"`
	ExpirationTime int64                `json:"expiration_time"`
}

type Transaction struct {
	Dynamodb  *dynamodb.DynamoDB
	TableName string
}

func (c Transaction) Get(id models.TransactionID) (models.Transaction, error) {
	db := c.Dynamodb

	params := &dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(string(id)),
			},
		},
	}

	result, err := db.GetItem(params)
	if err != nil {
		return models.Transaction{}, err
	}

	item := DynamoDBItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.Transaction{}, err
	}

	m, err := internalValue(item.Value).toModel()
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed to decode dynamodb value: %w", err)
	}

	return m, nil
}

func (c Transaction) Put(id models.TransactionID, value models.Transaction, currentTimeSec time.Time) error {
	expiredAt := currentTimeSec.Add(10 * time.Minute)

	v, err := newInternalValue(value)
	if err != nil {
		return fmt.Errorf("transaction cannnot encode to json: %w", err)
	}

	av, err := dynamodbattribute.MarshalMap(DynamoDBItem{
		ID:             id,
		Value:          string(v),
		ExpirationTime: expiredAt.Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed to map to dynamodb object: %w", err)
	}

	p := dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(c.TableName),
	}

	_, err = c.Dynamodb.PutItem(&p)
	if err != nil {
		return fmt.Errorf("failed to map to put object: %w", err)
	}

	return nil
}
