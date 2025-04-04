//go:build unit
// +build unit

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

type SqlResult struct {
	lastInsertId int64
	rowsAffected int64
}

func DefaultSqlResult() SqlResult {
	return SqlResult{
		lastInsertId: 0,
		rowsAffected: 0,
	}
}

func NewSqlResult(lastInsertId int64, rowsAffected int64) SqlResult {
	return SqlResult{
		lastInsertId: lastInsertId,
		rowsAffected: rowsAffected,
	}
}

// LastInsertId returns the integer generated by the database
// in response to a command. Typically this will be from an
// "auto increment" column when inserting a new row. Not all
// databases support this feature, and the syntax of such
// statements varies.
func (s *SqlResult) LastInsertId() (int64, error) {
	return int64(s.lastInsertId), nil
}

// RowsAffected returns the number of rows affected by an
// update, insert, or delete. Not every database or database
// driver may support this.
func (s *SqlResult) RowsAffected() (int64, error) {
	return int64(s.rowsAffected), nil
}

type StoreMock struct {
	result        SqlResult
	err           error
	getModel      interface{}
	executedCalls int
	executedQuery string
}

func NewStoreMock() StoreMock {
	return StoreMock{
		result:   DefaultSqlResult(),
		err:      nil,
		getModel: nil,
	}
}

func (s StoreMock) WithResult(r SqlResult) StoreMock {
	s.result = r
	return s
}

func (s StoreMock) WithError(e error) StoreMock {
	s.err = e
	return s
}

func (s StoreMock) WithModel(m interface{}) StoreMock {
	s.getModel = m
	return s
}

func (s StoreMock) WithDummyModel() StoreMock {
	s.getModel = 1
	return s
}

func (s StoreMock) WithSuccess(r SqlResult) StoreMock {
	s.result = r
	return s
}

func (s StoreMock) getExecResults() (int, string) {
	return s.executedCalls, s.executedQuery
}

func (s *StoreMock) ResetExecResults() {
	s.executedCalls = 0
	s.executedQuery = ""
}

func (s *StoreMock) MustExec(query string, args ...interface{}) sql.Result {
	s.executedCalls++
	s.executedQuery = query
	return &s.result
}

func (s *StoreMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.executedCalls++
	s.executedQuery = query
	return &s.result, s.err
}

func (s *StoreMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	s.executedCalls++
	s.executedQuery = query
	return &s.result, s.err
}

func (s *StoreMock) Get(modelRef interface{}, query string, args ...interface{}) error {
	s.executedCalls++
	s.executedQuery = query

	if s.getModel == nil {
		return fmt.Errorf("could not find the resource")
	}

	v := reflect.ValueOf(modelRef)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	v.Elem().Set(reflect.ValueOf(s.getModel).Elem())
	return s.err
}

func (s *StoreMock) Select(modelRef interface{}, query string, args ...interface{}) error {
	s.executedCalls++
	s.executedQuery = query

	if s.getModel == nil {
		return fmt.Errorf("could not find the resource")
	}

	v := reflect.ValueOf(modelRef)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to Select destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to Select destination")
	}

	destValue := v.Elem()
	sourceValue := reflect.ValueOf(s.getModel)

	if destValue.Kind() != reflect.Slice {
		return errors.New("destination is not a slice")
	}

	if sourceValue.Kind() != reflect.Slice {
		return errors.New("source is not a slice")
	}

	if destValue.Type().Elem() != sourceValue.Type().Elem() {
		return errors.New("incompatible slice element types")
	}

	for i := 0; i < sourceValue.Len(); i++ {
		destValue.Set(reflect.Append(destValue, sourceValue.Index(i)))
	}

	return s.err
}

func (s *StoreMock) Beginx() (*sqlx.Tx, error) {
	return nil, s.err
}

type TxMock struct {
	result        SqlResult
	err           error
	executedCalls int
	executedQuery string
}

func NewTxMock() *TxMock {
	return &TxMock{
		result: DefaultSqlResult(),
		err:    nil,
	}
}

func (t TxMock) getExecResults() (int, string) {
	return t.executedCalls, t.executedQuery
}

func (t TxMock) ResetExecResults() {
	t.executedCalls = 0
	t.executedQuery = ""
}

func (t TxMock) WithError(e error) TxMock {
	t.err = e
	return t
}

func (t *TxMock) Commit() error {
	return t.err
}

func (t *TxMock) Exec(query string, args ...interface{}) (sql.Result, error) {
	t.executedCalls++
	t.executedQuery = query
	return &t.result, t.err
}

func (t *TxMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	t.executedCalls++
	t.executedQuery = query
	return &t.result, t.err
}

func (t *TxMock) MustExec(query string, args ...interface{}) sql.Result {
	t.executedCalls++
	t.executedQuery = query
	return &t.result
}

func (t *TxMock) Rollback() error {
	t.executedCalls = 0
	return t.err
}

type AuthBackendMock struct {
	tokenStr    string
	isSuperUser bool
	token       jwt.Token
	err         error
}

func NewAuthBackendMock() AuthBackendMock {
	return AuthBackendMock{
		"",
		false,
		*jwt.New(jwt.SigningMethodHS256),
		nil,
	}
}

func (ab AuthBackendMock) NewJWT(userId string, orgId string, isSuperuser bool) (string, error) {
	return ab.tokenStr, ab.err
}

func (ab AuthBackendMock) ParseJWT(r *http.Request) (*jwt.Token, error) {
	return &ab.token, ab.err
}

func (ab AuthBackendMock) WithSuperUser() AuthBackendMock {
	ab.isSuperUser = true
	return ab
}

func (ab AuthBackendMock) WithToken(token jwt.Token) AuthBackendMock {
	ab.token = token
	return ab
}

type queryExecuter interface {
	getExecResults() (int, string)
}

// AssertQueryExecution asserts that an expected query was executed in a given db store
func AssertQueryExecution(t *testing.T, expectedQuery string, exec queryExecuter) {
	numExecs, execQuery := exec.getExecResults()
	if numExecs != 1 {
		t.Errorf("Expected query to be called 1 time, got %d", numExecs)
	} else {
		softCompareQueries(t, expectedQuery, execQuery)
	}
}

// SoftCompareQueries compares two queries, but it takes into consideration
// only the alphanumeric algarisms. This is done to ignore different spaces,
// tabulations and line breaks that would mistakenly qualify two queries as
// different.
func softCompareQueries(t *testing.T, want string, got string) {
	if cleanString(want) != cleanString(got) {
		t.Errorf("Expected query to be:\n%s\nBut got:\n%s", want, got)
	}
}

// cleanString erases all non-alphanumeric chars from a string
func cleanString(str string) string {
	pattern := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	result := pattern.ReplaceAllString(str, "")
	return result
}

// AssertEqualModels compares two structs and resgisters an error if they are not equal.
func AssertEqualModels(t *testing.T, want interface{}, got interface{}) {
	wantVal := reflect.ValueOf(want)
	gotVal := reflect.ValueOf(got)

	if wantVal.Type() != gotVal.Type() {
		t.Errorf("Structs have different types: %s and %s", wantVal.Type(), gotVal.Type())
	}

	for i := 0; i < wantVal.NumField(); i++ {
		wantField := wantVal.Field(i)
		gotField := gotVal.Field(i)
		fieldName := wantVal.Type().Field(i).Name

		if !reflect.DeepEqual(wantField.Interface(), gotField.Interface()) {
			t.Errorf("Field %s does not match. Expected: %v, Got: %v", fieldName, wantField.Interface(), gotField.Interface())
		}
	}
}
