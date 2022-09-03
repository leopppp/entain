package db

import (
	"testing"

	"github.com/leopppp/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
)

func TestApplyFilterWhenEmpty(t *testing.T) {
	racesRepo := &racesRepo{}
	query, args := racesRepo.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{})

	assert.Equal(t, "SELECT * FROM races", query)
	assert.Nil(t, args)
}

func TestApplyFilterWhenVisibleFalse(t *testing.T) {
	racesRepo := &racesRepo{}
	query, args := racesRepo.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		VisibleOnly: false,
	})

	assert.Equal(t, "SELECT * FROM races", query)
	assert.Nil(t, args)
}

func TestApplyFilterWhenVisibleTrue(t *testing.T) {
	racesRepo := &racesRepo{}
	query, args := racesRepo.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		VisibleOnly: true,
	})

	assert.Equal(t, "SELECT * FROM races WHERE visible = 1", query)
	assert.Nil(t, args)
}

func TestApplyFilterWhenMeetingIds(t *testing.T) {
	racesRepo := &racesRepo{}
	query, args := racesRepo.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		MeetingIds: []int64{5},
	})

	assert.Equal(t, "SELECT * FROM races WHERE meeting_id IN (?)", query)
	assert.ObjectsAreEqualValues([]int64{5}, args)
}

func TestApplyFilterWhenMeetingIdsAndVisibleTrue(t *testing.T) {
	racesRepo := &racesRepo{}
	query, args := racesRepo.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		MeetingIds: []int64{5, 8}, VisibleOnly: true,
	})

	assert.Equal(t, "SELECT * FROM races WHERE meeting_id IN (?,?) AND visible = 1", query)
	assert.ObjectsAreEqualValues([]int64{5, 8}, args)
}

func TestApplyFilterWhenMeetingIdsAndVisibleFalse(t *testing.T) {
	racesRepo := &racesRepo{}
	query, args := racesRepo.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{
		MeetingIds: []int64{5, 8}, VisibleOnly: false,
	})

	assert.Equal(t, "SELECT * FROM races WHERE meeting_id IN (?,?)", query)
	assert.ObjectsAreEqualValues([]int64{5, 8}, args)
}
