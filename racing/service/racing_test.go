package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/leopppp/entain/racing/db"
	"github.com/leopppp/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
)

var (
	racingDb        *sql.DB
	err             error
	racesRepo       db.RacesRepo
	myRacingService Racing
	ctx             context.Context
	cancel          context.CancelFunc
)

func init() {
	// Ideally, we should mock a database to do the tests.
	racingDb, err = sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}

	racesRepo = db.NewRacesRepo(racingDb)
	myRacingService = NewRacingService(racesRepo)
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
}

func TestListRacesWhenFilterEmpty(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 100)
}

func TestListRacesWhenVisibleOnlyFalse(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{VisibleOnly: false},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 100)
}

func TestListRacesWhenVisibleOnlyTrue(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{VisibleOnly: true},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 54)
}

func TestListRacesWhenMeetingIdsAndVisibleOnlyTrue(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: true},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 5)
}

func TestListRacesWhenMeetingIdsAndVisibleOnyFalse(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 15)
}

func TestListRacesWhenOrderByMeetingIdAscWithoutFilter(t *testing.T) {
	req := &racing.ListRacesRequest{
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "meeting_id", Asc: true},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 100)
	assert.Equal(t, res.Races[0].MeetingId, int64(1))
	assert.Equal(t, res.Races[99].MeetingId, int64(10))
}

func TestListRacesWhenOrderByMeetingIdDescWithoutFilter(t *testing.T) {
	req := &racing.ListRacesRequest{
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "meeting_id", Asc: false},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 100)
	assert.Equal(t, res.Races[0].MeetingId, int64(10))
	assert.Equal(t, res.Races[99].MeetingId, int64(1))
}

func TestListRacesOrderByAdvertisedStartTimeAscWithFilter(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter:  &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "advertised_start_time", Asc: true},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 15)
	assert.Equal(t, res.Races[0].AdvertisedStartTime.Seconds, int64(1614534476))
	assert.Equal(t, res.Races[14].AdvertisedStartTime.Seconds, int64(1614749300))
}

func TestListRacesWhenOrderByAdvertisedStartTimeDescWithFilter(t *testing.T) {
	req := &racing.ListRacesRequest{
		Filter:  &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "advertised_start_time", Asc: false},
	}

	res, err := myRacingService.ListRaces(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, len(res.Races), 15)
	assert.Equal(t, res.Races[0].AdvertisedStartTime.Seconds, int64(1614749300))
	assert.Equal(t, res.Races[14].AdvertisedStartTime.Seconds, int64(1614534476))
}

func TestGetRaceWhenNotExists(t *testing.T) {
	req := &racing.GetRaceRequest{Id: -1}
	res, err := myRacingService.GetRace(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, (*racing.Race)(nil), res.Race)
}

func TestGetRaceWhenExists(t *testing.T) {
	req := &racing.GetRaceRequest{Id: 5}
	res, err := myRacingService.GetRace(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), res.Race.Id)
}
