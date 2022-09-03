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

func TestListRacesWhenFilterEmpty(t *testing.T) {
	racingDb, err := sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{},
	}

	res, err := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 100)
}

func TestListRacesWhenVisibleFalse(t *testing.T) {
	racingDb, err := sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{VisibleOnly: false},
	}

	res, err := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 100)
}

func TestListRacesWhenVisibleTrue(t *testing.T) {
	racingDb, err := sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{VisibleOnly: true},
	}

	res, err := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 54)
}

func TestListRacesWhenMeetingIdsAndVisibleTrue(t *testing.T) {
	racingDb, err := sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: true},
	}

	res, err := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 5)
}

func TestListRacesWhenMeetingIdsAndVisibleFalse(t *testing.T) {
	racingDb, err := sql.Open("sqlite3", "../db/racing.db")
	if err != nil {
		panic(err)
	}
	racesRepo := db.NewRacesRepo(racingDb)
	racingService := NewRacingService(racesRepo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{MeetingIds: []int64{5}, VisibleOnly: false},
	}

	res, err := racingService.ListRaces(ctx, req)
	assert.Equal(t, len(res.Races), 15)
}
