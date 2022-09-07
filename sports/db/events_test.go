package db

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/leopppp/entain/sports/proto/sports"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplyFilterWhenEmpty(t *testing.T) {
	sportsRepo := &sportsRepo{}
	query, args := sportsRepo.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{})

	assert.Equal(t, "SELECT * FROM events", query)
	assert.Nil(t, args)
}

func TestApplyFilterWhenVisibleOnlyFalse(t *testing.T) {
	sportsRepo := &sportsRepo{}
	query, args := sportsRepo.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{
		VisibleOnly: false,
	})

	assert.Equal(t, "SELECT * FROM events", query)
	assert.Nil(t, args)
}

func TestApplyFilterWhenVisibleOnlyTrue(t *testing.T) {
	sportsRepo := &sportsRepo{}
	query, args := sportsRepo.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{
		VisibleOnly: true,
	})

	assert.Equal(t, "SELECT * FROM events WHERE visible = 1", query)
	assert.Nil(t, args)
}

func TestApplyOrderByNil(t *testing.T) {
	sportsRepo := &sportsRepo{}

	query := sportsRepo.applyOrderBy("SELECT * FROM events", nil)

	assert.Equal(t, "SELECT * FROM events", query)
}

func TestApplyOrderByPropertyEmpty(t *testing.T) {
	sportsRepo := &sportsRepo{}

	query := sportsRepo.applyOrderBy("SELECT * FROM events", &sports.ListEventsRequestOrderBy{
		Property: "",
		Asc:      true,
	})

	assert.Equal(t, "SELECT * FROM events", query)
}

func TestApplyOrderByAscendingWhenFilterEmpty(t *testing.T) {
	sportsRepo := &sportsRepo{}

	query := sportsRepo.applyOrderBy("SELECT * FROM events", &sports.ListEventsRequestOrderBy{
		Property: "advertised_start_time",
		Asc:      true,
	})

	assert.Equal(t, "SELECT * FROM events ORDER BY advertised_start_time ASC", query)
}

func TestApplyOrderByDescendingWhenFilterEmpty(t *testing.T) {
	sportsRepo := &sportsRepo{}

	query := sportsRepo.applyOrderBy("SELECT * FROM events", &sports.ListEventsRequestOrderBy{
		Property: "name",
		Asc:      false,
	})

	assert.Equal(t, "SELECT * FROM events ORDER BY name DESC", query)
}

func TestApplyOrderByAscendingWhenFilteringVisibleOnly(t *testing.T) {
	sportsRepo := &sportsRepo{}

	query, _ := sportsRepo.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{VisibleOnly: true})

	query = sportsRepo.applyOrderBy(query, &sports.ListEventsRequestOrderBy{
		Property: "advertised_start_time",
		Asc:      true,
	})

	assert.Equal(t, "SELECT * FROM events WHERE visible = 1 ORDER BY advertised_start_time ASC", query)
}

func TestApplyOrderByDescendingWhenFilteringVisibleOnly(t *testing.T) {
	sportsRepo := &sportsRepo{}

	query, _ := sportsRepo.applyFilter("SELECT * FROM events", &sports.ListEventsRequestFilter{VisibleOnly: true})

	query = sportsRepo.applyOrderBy(query, &sports.ListEventsRequestOrderBy{
		Property: "advertised_start_time",
		Asc:      false,
	})

	assert.Equal(t, "SELECT * FROM events WHERE visible = 1 ORDER BY advertised_start_time DESC", query)
}

func TestScanEventsWhenStatusClosed(t *testing.T) {
	sportsRepo := &sportsRepo{}
	mockRows := sqlmock.NewRows([]string{"id", "name", "address", "visible", "advertised_start_time"}).AddRow(123, "Fake Name", "MCG", true, timestamppb.New(time.Now().AddDate(0, 0, -10)).AsTime())
	sqlRows := mockSqlRows(mockRows)
	events, err := sportsRepo.scanEvents(sqlRows)

	assert.NoError(t, err)
	assert.Equal(t, sports.Status_CLOSED, events[0].Status)
}

func TestScanEventsWhenStatusOpen(t *testing.T) {
	sportsRepo := &sportsRepo{}
	mockRows := sqlmock.NewRows([]string{"id", "name", "address", "visible", "advertised_start_time"}).AddRow(123, "Fake Name", "GESAC", true, timestamppb.New(time.Now().AddDate(0, 0, 10)).AsTime())
	sqlRows := mockSqlRows(mockRows)
	events, err := sportsRepo.scanEvents(sqlRows)

	assert.NoError(t, err)
	assert.Equal(t, sports.Status_OPEN, events[0].Status)
}

func mockSqlRows(mockRows *sqlmock.Rows) *sql.Rows {
	db, mock, _ := sqlmock.New()
	mock.ExpectQuery("select").WillReturnRows(mockRows)
	rows, _ := db.Query("select")
	return rows
}
