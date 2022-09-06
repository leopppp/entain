package service

import (
	"context"
	"database/sql"
	"github.com/leopppp/entain/sports/proto/sports"
	"testing"
	"time"

	"github.com/leopppp/entain/sports/db"
	"github.com/stretchr/testify/assert"
)

var (
	sportsDb        *sql.DB
	err             error
	sportsRepo      db.SportsRepo
	mySportsService Sports
	ctx             context.Context
	cancel          context.CancelFunc
)

func init() {
	// Ideally, we should mock a database to do the tests.
	sportsDb, err = sql.Open("sqlite3", "../db/sports.db")
	if err != nil {
		panic(err)
	}

	sportsRepo = db.NewSportsRepo(sportsDb)
	mySportsService = NewSportsService(sportsRepo)
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
}

func TestListEventsWhenFilterEmpty(t *testing.T) {
	req := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{},
	}

	res, err := mySportsService.ListEvents(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Events), 100)
}

func TestListEventsWhenVisibleOnlyFalse(t *testing.T) {
	req := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{VisibleOnly: false},
	}

	res, err := mySportsService.ListEvents(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Events), 100)
}

func TestListEventsWhenVisibleOnlyTrue(t *testing.T) {
	req := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{VisibleOnly: true},
	}

	res, err := mySportsService.ListEvents(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Events), 49)
}
