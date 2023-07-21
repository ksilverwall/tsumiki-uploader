package repositories

import (
	"catch-all/models"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type internalItem struct {
	ID    string `json:"id"`
	Value string `json:"name"`
}

func (i *internalItem) toModel() (models.Transaction, error) {
	var m models.Transaction
	if err := json.Unmarshal([]byte(i.Value), &m); err != nil {
		return models.Transaction{}, fmt.Errorf("transaction cannnot decode from to json: %w", err)
	}

	return m, nil
}

func newInternalItem(m models.Transaction) (internalItem, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return internalItem{}, fmt.Errorf("transaction cannnot encode to json: %w", err)
	}

	return internalItem{
		ID:    m.ID,
		Value: string(b),
	}, nil
}

type Transaction struct {
	Dynamodb  *dynamodb.DynamoDB
	TableName string
}

func (c Transaction) Get(key string) (models.Transaction, error) {
	db := c.Dynamodb

	params := &dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	}

	result, err := db.GetItem(params)
	if err != nil {
		return models.Transaction{}, err
	}

	item := internalItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return models.Transaction{}, err
	}

	return item.toModel()
}

func (c Transaction) Put(key string, value models.Transaction) error {
	item, err := newInternalItem(value)
	if err != nil {
		return fmt.Errorf("failed to encode item: %w", err)
	}

	av, err := dynamodbattribute.MarshalMap(item)
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
