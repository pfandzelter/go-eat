package dynamo

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pfandzelter/go-eat/pkg/food"
)

// DB is a DynamoDB service for a particular table.
type DB struct {
	dynamodb *dynamodb.DynamoDB
	table    string
}

// New creates a new DynamoDB session.
func New(region string, table string) (*DB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
	})

	if err != nil {
		return nil, err
	}

	return &DB{
		dynamodb: dynamodb.New(sess),
		table:    table,
	}, nil
}

// PutFood puts one food item into the DynamoDB table.
func (d *DB) PutFood(c string, f []food.Food) error {
	item, err := json.Marshal(struct {
		Canteen string `json:"canteen"`
		Items []food.Food `json:"items"`
	}{
		Canteen: c,
		Items: f,
	})

	if err != nil {
		return err
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &d.table,
	}

	_, err = d.dynamodb.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}
