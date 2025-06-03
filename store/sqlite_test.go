package store

import (
	"context"
	"dig-inv/ent"
	"testing"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func TestSQLite(t *testing.T) {
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}

	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection to sqlite: %v", err)
		}
	}(client)

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}
}
