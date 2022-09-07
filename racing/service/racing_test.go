package service

import (
	"context"
	"database/sql"
	"github.com/leopppp/entain/racing/db"
	"testing"
	"time"

	"github.com/leopppp/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
)

var racingDb *sql.DB
var err error

func init() {
	racingDb, err = sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}
}

func TestListRacesWhenFilterEmpty(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 100)
}

func TestListRacesWhenVisibleFalse(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{VisibleOnly: false},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 100)
}

func TestListRacesWhenVisibleTrue(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{VisibleOnly: true},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 54)
}

func TestListRacesWhenMeetingIdsAndVisibleTrue(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: true},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 5)
}

func TestListRacesWhenMeetingIdsAndVisibleFalse(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 15)
}

func TestListRacesWhenOrderByMeetingIdAscWithoutFilter(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "meeting_id", Asc: true},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 100)
	assert.Equal(t, res.Races[0].MeetingId, int64(1))
	assert.Equal(t, res.Races[99].MeetingId, int64(10))
}

func TestListRacesWhenOrderByMeetingIdDescWithoutFilter(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "meeting_id", Asc: false},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 100)
	assert.Equal(t, res.Races[0].MeetingId, int64(10))
	assert.Equal(t, res.Races[99].MeetingId, int64(1))
}

func TestListRacesOrderByAdvertisedStartTimeAscWithFilter(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter:  &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "advertised_start_time", Asc: true},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 15)
	assert.Equal(t, res.Races[0].AdvertisedStartTime.Seconds, int64(1614534476))
	assert.Equal(t, res.Races[14].AdvertisedStartTime.Seconds, int64(1614749300))
}

func TestListRacesWhenOrderByAdvertisedStartTimeDescWithFilter(t *testing.T) {
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter:  &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
		OrderBy: &racing.ListRacesRequestOrderBy{Property: "advertised_start_time", Asc: false},
	}

	res, _ := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 15)
	assert.Equal(t, res.Races[0].AdvertisedStartTime.Seconds, int64(1614749300))
	assert.Equal(t, res.Races[14].AdvertisedStartTime.Seconds, int64(1614534476))
}
