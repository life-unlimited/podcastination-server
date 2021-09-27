package app

import (
	"database/sql"
	nativeerrors "errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/life-unlimited/podcastination-server/embedded"
	"github.com/pkg/errors"
	"log"
)

// defaultMaxDBConnections is the maximum number of database connections that is used when no other one is provided
// in the Config.
const defaultMaxDBConnections = 16

// dbVersion is used for determining the current database version. This is saved in a special table when properly set
// up. If the version does not exist, one can know that the database needs to be initialized. If it is and the latest
// version is greater, migrations can be performed.
type dbVersion string

// dbVersionZero is used when no database version could be found, and therefore we conclude that it has not been
// initialized yet.
const dbVersionZero dbVersion = "0"

// dbMigration is used for performing and checking database migrations. They lie in dbMigrations which is an ordered
// list of versions with their migrations.
type dbMigration struct {
	version dbVersion
	up      string
}

// dbMigrations are the sql migrations in an ordered (!) list. The order is used to determine which migrations need to
// be done when the current database version is not the latest one.
var dbMigrations = []dbMigration{
	{
		version: "1.0",
		up:      embedded.DBMigration1x0,
	},
}

// connectDB connects to the database with the given connection string and returns the connection pool.
func connectDB(connectionStr string, maxDBConnections int) (*sql.DB, error) {
	dbPool, err := sql.Open("pgx", connectionStr)
	dbPool.SetMaxOpenConns(maxDBConnections)
	if err != nil {
		return nil, errors.Wrap(err, "connect to database")
	}
	// Perform test query.

	return dbPool, nil
}

// testDBConnection tests the database connection by simply querying 1.
func testDBConnection(db *sql.DB) error {
	// Build test query.
	q, _, err := goqu.Select(goqu.V(1)).ToSQL()
	if err != nil {
		return errors.Wrap(err, "build test query to sql")
	}
	// Query database.
	result := db.QueryRow(q)
	var got int
	err = result.Scan(&got)
	if err != nil {
		return errors.Wrap(err, "perform test query in database")
	}
	// Assure that we got 1.
	if got != 1 {
		return errors.Wrap(err, "invalid result for test query in database")
	}
	return nil
}

// performDBMigrations performs all needed database migrations according to the (un)set database version.
func performDBMigrations(db *sql.DB) error {
	currentVersion, err := retrieveCurrentDBVersion(db)
	if err != nil {
		return errors.Wrap(err, "retrieve current db version")
	}
	migrationsToDo, err := getDBMigrationsToDo(currentVersion)
	if err != nil {
		return errors.Wrap(err, "get db migrations to do")
	}
	// Check if migrations need to be performed.
	if len(migrationsToDo) == 0 {
		return nil
	}
	// Begin tx for avoiding database destruction if something fails.
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	// Perform migrations.
	var newVersion dbVersion
	for i, migration := range migrationsToDo {
		log.Printf("performing database migration %d/%d...", i+1, len(migrationsToDo))
		// Perform migration according to the version.
		_, err = tx.Exec(migration.up)
		if err != nil {
			rollbackTx(tx, "database migration failed")
			return errors.Wrap(err, "database migration failed")
		}
		newVersion = migration.version
	}
	// Update database version.
	var updateDBVersionQuery string
	if currentVersion == dbVersionZero {
		updateDBVersionQuery, _, err = goqu.Dialect("postgres").Insert(goqu.T("podcastination")).Rows(goqu.Record{
			"key":   "db-version",
			"value": newVersion,
		}).ToSQL()
	} else {
		updateDBVersionQuery, _, err = goqu.Dialect("postgres").Update(goqu.T("podcastination")).
			Set(goqu.Record{"value": newVersion}).
			Where(goqu.C("key").Eq("db-version")).ToSQL()
	}
	if err != nil {
		rollbackTx(tx, "update database version query to sql failed")
		return nil
	}
	_, err = tx.Exec(updateDBVersionQuery)
	if err != nil {
		rollbackTx(tx, "update database version failed")
		return nil
	}
	// Commit tx.
	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit tx")
	}
	// All done.
	return nil
}

// getDBMigrationsToDo retrieves all database migrations that need to be performed. If the version is dbVersionZero, it
// will return all migrations. If the version is unknown, an error will be returned.
func getDBMigrationsToDo(currentVersion dbVersion) ([]dbMigration, error) {
	// Check if empty version.
	if currentVersion == dbVersionZero {
		return dbMigrations, nil
	}
	found := false
	migrationsToDo := make([]dbMigration, 0)
	for _, migration := range dbMigrations {
		if migration.version == currentVersion {
			// Match found.
			if found {
				// This should not happen and is an internal error as the versions are not properly set up. What did you
				// do?
				return nil, errors.New(fmt.Sprintf("duplicate database version %v in available migrations", currentVersion))
			}
			// Set found flag.
			found = true
			// Continue with next one as we already performed everything for this database version.
			continue
		}
		// Append migration to todos.
		migrationsToDo = append(migrationsToDo, migration)
	}
	// Check if found.
	if !found {
		return nil, errors.New(fmt.Sprintf("no database version found matching %v", currentVersion))
	}
	// Done.
	return migrationsToDo, nil
}

// retrieveCurrentDBVersion retrieves the current dbVersion from the given database. If no version could be found,
// dbVersionZero will be returned.
func retrieveCurrentDBVersion(db *sql.DB) (dbVersion, error) {
	versionStr, ok, err := retrieveKeyValFromDB(db, "db-version")
	if err != nil {
		return "", errors.Wrap(err, "retrieve key val from database")
	}
	if !ok {
		return dbVersionZero, nil
	}
	return dbVersion(versionStr), nil
}

// retrieveKeyValFromDB retrieves the value for the given key from the given database. If the table does not exist, we
// do not care and expect the caller to have already checked or expect this.
func retrieveKeyValFromDB(db *sql.DB, key string) (string, bool, error) {
	// Build query.
	q, _, err := goqu.Dialect("postgres").From(goqu.T("podcastination")).
		Select(goqu.C("value")).
		Where(goqu.C("key").Eq(key)).ToSQL()
	if err != nil {
		return "", false, errors.Wrap(err, "query to sql")
	}
	// Exec query and scan value.
	var value string
	err = db.QueryRow(q).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		// Check if error is because relation does not exist as then it's a not-found error.
		var pgErr *pgconn.PgError
		if nativeerrors.As(err, &pgErr) && pgErr.Code == "42P01" {
			return "", false, nil
		} else {
			fmt.Println(nativeerrors.As(err, &pgErr))
		}
		return "", false, errors.Wrap(err, "query and scan row")
	}
	// Done.
	return value, true, nil
}

// rollbackTx rolls back the given sql.Tx. The encapsulation is needed because rolling back might return an error which
// does not need to be returned but definitely logged with the original reason the rollback was performed.
func rollbackTx(tx *sql.Tx, reason string) {
	err := tx.Rollback()
	if err != nil {
		log.Printf("err: rollback tx because of %s: %s", reason, err)
	}
}
