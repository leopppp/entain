package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/leopppp/entain/sports/proto/sports"
)

// SportsRepo provides repository access to events.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter, orderBy *sports.ListEventsRequestOrderBy) ([]*sports.Event, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sports repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = s.seed()
	})

	return err
}

func (s *sportsRepo) List(filter *sports.ListEventsRequestFilter, orderBy *sports.ListEventsRequestOrderBy) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = s.applyFilter(query, filter)

	query = s.applyOrderBy(query, orderBy)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	// Only return visible only events
	if filter.GetVisibleOnly() {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (s *sportsRepo) applyOrderBy(query string, orderBy *sports.ListEventsRequestOrderBy) string {
	if orderBy == nil || orderBy.Property == "" {
		return query
	}

	var sortingOrder string
	if orderBy.GetAsc() {
		sortingOrder = "ASC"
	} else {
		sortingOrder = "DESC"
	}
	query = fmt.Sprintf("%s ORDER BY %s %s", query, orderBy.GetProperty(), sortingOrder)

	return query
}

func (s *sportsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &event.Address, &event.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		if err := setAdvertisedStartTime(advertisedStart, &event); err != nil {
			return nil, err
		}

		setStatus(advertisedStart, &event)

		events = append(events, &event)
	}

	return events, nil
}

func setAdvertisedStartTime(advertisedStart time.Time, event *sports.Event) error {
	ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return err
	}
	event.AdvertisedStartTime = ts

	return nil
}

func setStatus(advertisedStart time.Time, event *sports.Event) {
	if time.Now().Before(advertisedStart) {
		// Set status to open if advertised start is in the future
		event.Status = sports.Status_OPEN
	}
}
