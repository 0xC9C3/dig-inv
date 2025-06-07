package store

import (
	"context"
	"dig-inv/ent"
	"dig-inv/log"
	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

var Client *ent.Client

// GetClient @todo check docs if this does connection pooling etc?
func GetClient() (*ent.Client, error) {
	if Client != nil {
		return Client, nil
	}

	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}

	Client = client
	return Client, nil
}

func InitializeSchema(ctx context.Context) error {
	client, err := GetClient()
	if err != nil {
		log.S.Errorw("Failed to get store client", "error", err)
		return err
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.S.Errorw("Failed to create schema resources", "error", err)
		return err
	}

	log.S.Infow("Schema initialized successfully")

	return nil
}
