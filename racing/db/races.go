package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/leopppp/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, orderBy *racing.ListRacesRequestOrderBy) ([]*racing.Race, error)

	// Get will return a single race by its ID.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, orderBy *racing.ListRacesRequestOrderBy) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.applyOrderBy(query, orderBy)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var args []interface{}

	query := getRaceQueries()[racesList] + " WHERE id = ?"
	args = append(args, id)

	row := r.db.QueryRow(query, args...)

	return r.scanRace(row)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// Only return visible only races
	if filter.GetVisibleOnly() {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applyOrderBy(query string, orderBy *racing.ListRacesRequestOrderBy) string {
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

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		if err := setAdvertisedStartTime(advertisedStart, &race); err != nil {
			return nil, err
		}

		setStatus(advertisedStart, &race)

		races = append(races, &race)
	}

	return races, nil
}

func (m *racesRepo) scanRace(row *sql.Row) (*racing.Race, error) {
	var race racing.Race
	var advertisedStart time.Time

	if err := row.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if err := setAdvertisedStartTime(advertisedStart, &race); err != nil {
		return nil, err
	}

	setStatus(advertisedStart, &race)

	return &race, nil
}

func setAdvertisedStartTime(advertisedStart time.Time, race *racing.Race) error {
	ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return err
	}
	race.AdvertisedStartTime = ts

	return nil
}

func setStatus(advertisedStart time.Time, race *racing.Race) {
	if time.Now().Before(advertisedStart) {
		// Set status to open if advertised start is in the future
		race.Status = racing.Status_OPEN
	}
}
