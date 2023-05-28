package postgres

import (
	"context"
	"fmt"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	pgx2 "github.com/glbter/currency-ex/pkg/sql/pgx"
	"github.com/glbter/currency-ex/pkg/stesting"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

const dbSchema = "portfolio_test_schema"

var (
	testDB sqlc.DB
	fx     *stesting.FixtureLoader
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("DSN is required to run tests")
		os.Exit(1)
	}

	ctx := context.Background()

	p, err := pgx2.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db := pgx2.NewDB(p, dbSchema)
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

	fx = stesting.NewFixtureLoader(path.Join(".", "testdata"))
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

func pointerTime(t time.Time) *time.Time {
	return &t
}

func pointerFloat64(f float64) *float64 {
	return &f
}

func pointerString(s string) *string {
	return &s
}
