package app

import (
	"database/sql"
	nativeerrors "errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/suite"
	"testing"
)

type dbSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
}

// prepareDB creates a new mock database and sets the according fields in the suite.
func (suite *dbSuite) prepareDB() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	suite.Require().Nil(err, "creating mock database should not fail")
	suite.db = db
	suite.mock = mock
}

type RetrieveKeyValFromDBTestSuite struct {
	dbSuite
	key           string
	retrieveQuery string
}

func (suite *RetrieveKeyValFromDBTestSuite) prepareQuery(key string) {
	q, _, err := goqu.Dialect("postgres").From(goqu.T("podcastination")).
		Select("value").
		Where(goqu.C("key").Eq(key)).ToSQL()
	suite.Require().Nil(err, "query to sql should not fail")
	suite.retrieveQuery = q
	suite.key = key
}

func (suite *RetrieveKeyValFromDBTestSuite) SetupTest() {
	suite.prepareDB()
}

func (suite *RetrieveKeyValFromDBTestSuite) TeardownTest() {
	suite.Assert().Nil(suite.mock.ExpectationsWereMet(), "all mock expectations should be met")
}

func (suite *RetrieveKeyValFromDBTestSuite) TestQueryFail() {
	suite.prepareQuery("hello")
	suite.mock.ExpectQuery(suite.retrieveQuery).WillReturnError(nativeerrors.New("ERROR"))

	_, _, err := retrieveKeyValFromDB(suite.db, suite.key)
	suite.Assert().NotNil(err, "retrieval should fail")
}

func (suite *RetrieveKeyValFromDBTestSuite) TestKeyNotFound() {
	suite.prepareQuery("i-am-unknown")
	suite.mock.ExpectQuery(suite.retrieveQuery).WillReturnRows(sqlmock.NewRows([]string{"value"}))

	_, ok, err := retrieveKeyValFromDB(suite.db, suite.key)
	suite.Assert().Nil(err, "retrieval should fail because of unknown key")
	suite.Assert().Falsef(ok, "should not be ok")
}

func (suite *RetrieveKeyValFromDBTestSuite) TestOK() {
	suite.prepareQuery("i-am-known")
	suite.mock.ExpectQuery(suite.retrieveQuery).WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow("Hello World!"))

	value, ok, err := retrieveKeyValFromDB(suite.db, suite.key)
	suite.Require().Nil(err, "retrieval should not fail")
	suite.Require().Truef(ok, "should be ok")
	suite.Assert().Equal("Hello World!", value, "should be expected value")
}

func Test_retrieveKeyValFromDB(t *testing.T) {
	suite.Run(t, new(RetrieveKeyValFromDBTestSuite))
}

type DBMigrationsTestSuite struct {
	suite.Suite
}

func (suite *DBMigrationsTestSuite) TestCorrectOrder() {
	var prevVersion *version.Version
	for i, migration := range dbMigrations {
		v, err := version.NewVersion(string(migration.version))
		suite.Require().Nil(err, "creating version should not fail")
		if i > 0 {
			suite.Assert().Truef(v.GreaterThan(prevVersion), "version %s must be greater than previous version %s",
				v.String(), prevVersion.String())
		}
		prevVersion = v
	}
}

func (suite *DBMigrationsTestSuite) TestNoVersionZero() {
	for _, migration := range dbMigrations {
		suite.Assert().NotEqual(dbVersionZero, migration.version, "no available version should match version zero")
	}
}

func (suite *DBMigrationsTestSuite) TestNoDuplicates() {
	knownVersions := make(map[dbVersion]struct{})
	for _, migration := range dbMigrations {
		_, found := knownVersions[migration.version]
		suite.Assert().Falsef(found, "version %s should have no duplicate in available ones", migration.version)
		knownVersions[migration.version] = struct{}{}
	}
}

func Test_dbMigrations(t *testing.T) {
	suite.Run(t, new(DBMigrationsTestSuite))
}

type GetDBMigrationsToDoTestSuite struct {
	suite.Suite
}

func (suite *GetDBMigrationsToDoTestSuite) TestLatest() {
	migrations, err := getDBMigrationsToDo(dbMigrations[len(dbMigrations)-1].version)
	suite.Require().Nilf(err, "retrieval should not fail but got %s", err)
	suite.Assert().Len(migrations, 0, "should return no migrations to do because current version is latest")
}

func (suite *GetDBMigrationsToDoTestSuite) TestUnknownVersion() {
	_, err := getDBMigrationsToDo(dbVersion(fmt.Sprintf("%s-unknown-version-lol", dbMigrations[len(dbMigrations)-1])))
	suite.Assert().NotNil(err, "retrieval should fail because of unknown version")
}

func (suite *GetDBMigrationsToDoTestSuite) TestVersionZero() {
	migrations, err := getDBMigrationsToDo(dbVersionZero)
	suite.Require().Nilf(err, "retrieval should not fail but got %s", err)
	suite.Assert().Len(migrations, len(dbMigrations), "should return all migrations")
}

func Test_getDBMigrationsToDo(t *testing.T) {
	suite.Run(t, new(GetDBMigrationsToDoTestSuite))
}

type PerformDBMigrationsTestSuite struct {
	dbSuite
	keyVal               RetrieveKeyValFromDBTestSuite
	updateDBVersionQuery string
}

// prepareMigrationQueries sets the PerformDBMigrationsTestSuite.updateDBVersionQuery with correct zero stuff.
func (suite *PerformDBMigrationsTestSuite) prepareUpdateDBVersionQuery(version dbVersion, isZero bool) {
	var q string
	var err error
	if isZero {
		q, _, err = goqu.Dialect("postgres").Insert(goqu.T("podcastination")).Rows(goqu.Record{
			"key":   "db-version",
			"value": version,
		}).ToSQL()
	} else {
		q, _, err = goqu.Dialect("postgres").Update(goqu.T("podcastination")).Set(goqu.Record{
			"value": version,
		}).Where(goqu.C("key").Eq("db-version")).ToSQL()
	}
	suite.Require().Nilf(err, "update db version query should not fail but got %s", err)
	suite.updateDBVersionQuery = q
}

func (suite *PerformDBMigrationsTestSuite) SetupTest() {
	suite.prepareDB()
	suite.keyVal.dbSuite = suite.dbSuite
}

func (suite *PerformDBMigrationsTestSuite) TeardownTest() {
	suite.Assert().Nil(suite.mock.ExpectationsWereMet(), "all mock expectations should be met")
}

func (suite *PerformDBMigrationsTestSuite) TestCurrentVersionRetrievalFail() {
	suite.keyVal.prepareQuery("db-version")
	suite.mock.ExpectQuery(suite.keyVal.retrieveQuery).WillReturnError(nativeerrors.New("ERROR"))

	err := performDBMigrations(suite.db)
	suite.Assert().NotNil(err, "should fail")
}

func (suite *PerformDBMigrationsTestSuite) TestVersionZero() {
	suite.keyVal.prepareQuery("db-version")
	suite.mock.ExpectQuery(suite.keyVal.retrieveQuery).WillReturnRows(sqlmock.NewRows([]string{"value"}))
	suite.mock.ExpectBegin()
	for _, migration := range dbMigrations {
		suite.mock.ExpectExec(migration.up).WillReturnResult(sqlmock.NewResult(0, 0))
	}
	suite.prepareUpdateDBVersionQuery(dbMigrations[len(dbMigrations)-1].version, true)
	suite.mock.ExpectExec(suite.updateDBVersionQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := performDBMigrations(suite.db)
	suite.Assert().Nilf(err, "should not fail but got %s", err)
}

func (suite *PerformDBMigrationsTestSuite) TestUnknownVersion() {
	suite.keyVal.prepareQuery("db-version")
	suite.mock.ExpectQuery(suite.keyVal.retrieveQuery).WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow("i-am-unknown"))

	err := performDBMigrations(suite.db)
	suite.Assert().NotNil(err, "should fail")
}

func (suite *PerformDBMigrationsTestSuite) TestLatest() {
	suite.keyVal.prepareQuery("db-version")
	suite.mock.ExpectQuery(suite.keyVal.retrieveQuery).WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(dbMigrations[len(dbMigrations)-1].version))

	err := performDBMigrations(suite.db)
	suite.Assert().Nilf(err, "should not fail but got %s", err)
}

func Test_performDBMigrations(t *testing.T) {
	suite.Run(t, new(PerformDBMigrationsTestSuite))
}
