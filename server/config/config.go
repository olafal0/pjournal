package config

import "os"

const userPoolEnvKey = "COGNITO_USER_POOL_ID"
const postsTableNameEnvKey = "POSTS_TABLE_NAME"

// PjournalConfig is an object containing global configuration for the app.
type PjournalConfig struct {
	UserPoolID     string
	PostsTableName string
}

// NewFromEnv returns a PjournalConfig object initialized from environment variables.
func NewFromEnv() *PjournalConfig {
	config := &PjournalConfig{}
	config.UserPoolID = os.Getenv(userPoolEnvKey)
	config.PostsTableName = os.Getenv(postsTableNameEnvKey)

	return config
}
