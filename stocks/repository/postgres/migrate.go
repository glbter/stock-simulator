package postgres

import (
	"context"
	"embed"
	sqlc "github.com/glbter/currency-ex/sql"
)

//go:embed migrations
var migrations embed.FS

func Migrate(ctx context.Context, db sqlc.DB) (int, error) {
	return db.Migrate(ctx, migrations, "migrations", "migration")
}
