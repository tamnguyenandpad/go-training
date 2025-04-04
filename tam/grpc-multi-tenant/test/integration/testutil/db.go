package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
)

type env struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

func InitDB(t *testing.T) (DB *sql.DB, dbName string) {
	env := &env{
		DBUser: "root",
		DBPass: "password",
		DBHost: "localhost",
		DBPort: "3307",
		DBName: "test_db",
	}

	db, dbCloseFunc, err := createCleanDB(env)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := dbCloseFunc(); err != nil {
			t.Fatal(err)
		}

		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	return db, dbName
}

func createCleanDB(e *env) (db *sql.DB, closeFunc func() error, err error) {
	dbName := fmt.Sprintf("%s_%s", e.DBName, uuid.NewString())
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", e.DBUser, e.DBPass, e.DBHost, e.DBPort)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	fmt.Println("create db", dbName)

	createDB(sqlDB, dbName)

	createTables(sqlDB)

	loadFixtures(sqlDB)

	closeFunc = func() error {
		err = dropDatabase(sqlDB, dbName)
		if err != nil {
			return fmt.Errorf("drop error %w", err)
		}

		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("sqlDB.Close() error %w", err)
		}

		return nil
	}
	testDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", e.DBUser, e.DBPass, e.DBHost, e.DBPort, dbName))
	if err != nil {
		panic(err)
	}
	return testDB, closeFunc, nil
}

func createDB(sqlDB *sql.DB, dbName string) {
	createFmt := fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbName)
	_, err := sqlDB.Exec(createFmt)
	if err != nil {
		panic(err)
	}

	if _, err := sqlDB.Exec(fmt.Sprintf("USE `%s`;", dbName)); err != nil {
		panic(err)
	}
}

func createTables(sqlDB *sql.DB) {
	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller error")
	}

	dumpPath := filepath.Join(filepath.Dir(thisFilePath), "..", "..", "..", "..")
	dumpPath = filepath.Join(dumpPath, "grpc-multi-tenant", "database", "import", "create_tables.sql")
	fmt.Println(dumpPath)

	dumpBytes, err := os.ReadFile(dumpPath)
	if err != nil {
		panic(err)
	}

	regexpNewline := regexp.MustCompile(`\r\n|\r|\n`)
	dumpStr := regexpNewline.ReplaceAllString(string(dumpBytes), "")

	for _, stmt := range strings.Split(dumpStr, ";") {
		if stmt == "" {
			continue
		}
		if _, err := sqlDB.Exec(stmt); err != nil {
			panic(err)
		}
	}
}

func dropDatabase(sqlDB *sql.DB, dbName string) error {
	_, err := sqlDB.Exec(fmt.Sprintf("DROP DATABASE `%s`;", dbName))
	if err != nil {
		return fmt.Errorf("drop error %w", err)
	}
	return nil
}

func loadFixtures(sqlDB *sql.DB) {
	// create data
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	fixturesPath := fmt.Sprintf("%s/testdata", dir)
	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(fixturesPath),
	)
	if err != nil {
		panic(err)
	}

	err = fixtures.Load()
	if err != nil {
		panic(err)
	}
}
