package dbinterface

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	uuid "github.com/satori/go.uuid"
)

// DynamoID is a string UUID
type DynamoID string

// HashKey fulfills the dynamo.Keyed interface.
func (d *DynamoID) HashKey() interface{} {
	return string(*d)
}

// RangeKey fulfills the dynamo.Keyed interface.
func (d *DynamoID) RangeKey() interface{} {
	return nil
}

// DynamoKey is a primary key/sort
type DynamoKey struct {
	PrimaryKey string
	SortKey    string
}

// HashKey fulfills the dynamo.Keyed interface.
func (d *DynamoKey) HashKey() interface{} {
	return d.PrimaryKey
}

// RangeKey fulfills the dynamo.Keyed interface.
func (d *DynamoKey) RangeKey() interface{} {
	if d.SortKey == "" {
		return nil
	}
	return d.SortKey
}

// DynamoIndex represents an index on a dynamo table.
type DynamoIndex struct {
	PrimaryKey string
	SortKey    string
}

// DynamoHandler contains dynamo-related configuration for specific tables.
type DynamoHandler struct {
	TableName  string
	PrimaryKey string
	SortKey    string
	Indexes    map[string]DynamoIndex
	client     *dynamo.DB
}

var dynamoDB *dynamo.DB

func (d *DynamoHandler) getClient() (*dynamo.DB, error) {
	if d.client != nil {
		return d.client, nil
	}
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	d.client = dynamo.New(sess)
	return d.client, nil
}

func (d *DynamoHandler) newID() string {
	return uuid.NewV4().String()
}
