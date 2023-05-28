package stesting

import (
	"context"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"os"
	"path/filepath"
	"testing"
)

func NewFixtureLoader(basePath string) *FixtureLoader {
	return &FixtureLoader{
		basePath: basePath,
	}
}

type FixtureLoader struct {
	basePath string
}

func (fx FixtureLoader) ExecSQLFixture(db sqlc.DB, fileName string) error {
	res, err := os.ReadFile(filepath.Join(fx.basePath, filepath.Clean(fileName)))
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), string(res))
	if err != nil {
		return err
	}

	return nil
}

func (fx FixtureLoader) MustLoadStringFixture(t *testing.T, fileName string) string {
	t.Helper()
	res, err := os.ReadFile(filepath.Join(fx.basePath, "fixture", filepath.Clean(fileName)))
	if err != nil {
		t.Fatal(err.Error())
	}

	return string(res)
}

func (fx FixtureLoader) MustExecSQLFixture(t *testing.T, db sqlc.DB, fileName string) {
	t.Helper()
	if err := fx.ExecSQLFixture(db, fileName); err != nil {
		t.Fatal(err.Error())
	}
}
