package pgx

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/tern/migrate"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (db *DB) fullMigrationTableName(table string) string {
	if strings.Contains(table, ".") {
		return table
	}

	return db.schema + "." + table
}

func (db *DB) Migrate(ctx context.Context, fs embed.FS, dir, table string) (int, error) {
	table = db.fullMigrationTableName(table)

	c, err := db.conpool.pool.Acquire(ctx)
	if err != nil {
		return -1, fmt.Errorf("acquire database connection %w", err)
	}
	defer c.Release()

	opts := &migrate.MigratorOptions{MigratorFS: fsAdapter{FS: fs, root: dir}}

	migrator, err := migrate.NewMigratorEx(ctx, c.Conn(), table, opts)
	if err != nil {
		return -1, fmt.Errorf("create a migrator: %w", err)
	}

	migrator.OnStart = func(_ int32, _, _, _ string) {
		if err := db.setSchema(ctx, c); err != nil {
			panic(fmt.Errorf("change schema inside migrator.OnStart: %w", err))
		}
	}

	if err := migrator.LoadMigrations(dir); err != nil {
		return -1, fmt.Errorf("load migrations: %w", err)
	}

	if err := migrator.Migrate(ctx); err != nil {
		return -1, fmt.Errorf("migrate: %w", err)
	}

	v, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return -1, fmt.Errorf("get current schema version: %w", err)
	}

	return int(v), nil
}

type fsAdapter struct {
	embed.FS
	root string
}

func (a fsAdapter) ReadDir(dirname string) ([]os.FileInfo, error) {
	entries, err := a.FS.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	infos := make([]fs.FileInfo, 0, len(entries))
	for _, e := range entries {
		fi, err := e.Info()
		if err != nil {
			continue // no file
		}

		infos = append(infos, fi)
	}

	return infos, nil
}

func (a fsAdapter) Glob(pattern string) ([]string, error) {
	des, err := a.FS.ReadDir(a.root)
	if err != nil {
		return nil, fmt.Errorf("read from pattern as path: %w", err)
	}

	files := make([]string, 0, len(des))
	for _, e := range des {
		matches, err := filepath.Match(pattern, e.Name())
		if err != nil {
			return nil, err
		}

		if !matches {
			continue
		}

		files = append(files, filepath.Join(a.root, e.Name()))
	}

	return files, nil

}
