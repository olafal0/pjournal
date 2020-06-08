package dbinterface

import (
	"server/config"
)

// A DatabaseModifier is an object that handles database operations.
type DatabaseModifier interface {
	PostsModifier
}

// DatabaseHandler is an object that fulfills the DatabaseModifier interface.
type DatabaseHandler struct {
	Posts
}

// NewModifier creates a new object that satisfies the DatabaseModifier interface.
func NewModifier(cfg *config.PjournalConfig) DatabaseModifier {
	db := DatabaseHandler{
		Posts: Posts{
			Dynamo: DynamoHandler{
				TableName:  cfg.PostsTableName,
				PrimaryKey: "id",
				SortKey:    "userId",
				Indexes: map[string]DynamoIndex{
					"userId-createdAt-index": DynamoIndex{
						PrimaryKey: "userId",
						SortKey:    "createdAt",
					},
				},
			},
			UsercreatedAtIndexName: "userId-createdAt-index",
		},
	}

	return &db
}
