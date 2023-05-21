package postgres

import (
	"context"
	"fmt"
	sqlc "github.com/glbter/currency-ex/sql"
	"github.com/glbter/currency-ex/sql/pgx"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const dbSchema = "portfolio_test_schema"

var (
	testDB sqlc.DB
	fx     *FixtureLoader
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("DSN is required to run tests")
		os.Exit(1)
	}

	ctx := context.Background()

	p, err := pgx.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db := pgx.NewDB(p, dbSchema)
	if err := sqlc.PrepareSchema(ctx, db, dbSchema); err != nil {
		log.Fatalln(err.Error())
	}

	if err := db.SetSchema(ctx); err != nil {
		log.Fatalln(err.Error())
	}

	//c, err := Migrate(ctx, db)
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}
	//
	//fmt.Printf("%d migrations applied\n", c)

	fx = NewFixtureLoader(path.Join(".", "testdata"))
	testDB = db

	os.Exit(m.Run())
}

func mustTruncateTables(tb testing.TB, exec sqlc.Executor) {
	tb.Helper()

	tb.Log("clean up")

	_, err := exec.Exec(
		context.Background(),
		`TRUNCATE TABLE 
	portfolio_record,
	split,
	stock_daily, 
	ticker
RESTART IDENTITY CASCADE`)

	if err != nil {
		tb.Fatal(err.Error())
	}
}

func NewFixtureLoader(basePath string) *FixtureLoader {
	return &FixtureLoader{
		basePath: basePath,
	}
}

type FixtureLoader struct {
	basePath string
}

func (fx FixtureLoader) execSQLFixture(db sqlc.DB, fileName string) error {
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

func (fx FixtureLoader) mustExecSQLFixture(t *testing.T, db sqlc.DB, fileName string) {
	t.Helper()
	if err := fx.execSQLFixture(db, fileName); err != nil {
		t.Fatal(err.Error())
	}
}

func pointerFloat64(f float64) *float64 {
	return &f
}

func pointerString(s string) *string {
	return &s
}
