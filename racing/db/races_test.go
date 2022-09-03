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

func TestApplyOrderByNil(t *testing.T) {
	rr := &racesRepo{}

	query := rr.applyOrderBy("SELECT * FROM races", nil)

	assert.Equal(t, "SELECT * FROM races", query)
}

func TestApplyOrderByPropertyEmpty(t *testing.T) {
	rr := &racesRepo{}

	query := rr.applyOrderBy("SELECT * FROM races", &racing.ListRacesRequestOrderBy{
		Property: "",
		Asc:      true,
	})

	assert.Equal(t, "SELECT * FROM races", query)
}

func TestApplyOrderByAscendingWhenFilterEmpty(t *testing.T) {
	rr := &racesRepo{}

	query := rr.applyOrderBy("SELECT * FROM races", &racing.ListRacesRequestOrderBy{
		Property: "advertised_start_time",
		Asc:      true,
	})

	assert.Equal(t, "SELECT * FROM races ORDER BY advertised_start_time ASC", query)
}

func TestApplyOrderByDescendingWhenFilterEmpty(t *testing.T) {
	rr := &racesRepo{}

	query := rr.applyOrderBy("SELECT * FROM races", &racing.ListRacesRequestOrderBy{
		Property: "meeting_id",
		Asc:      false,
	})

	assert.Equal(t, "SELECT * FROM races ORDER BY meeting_id DESC", query)
}

func TestApplyOrderByAscendingWhenHavingFilter(t *testing.T) {
	rr := &racesRepo{}

	query, _ := rr.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{VisibleOnly: true})

	query = rr.applyOrderBy(query, &racing.ListRacesRequestOrderBy{
		Property: "advertised_start_time",
		Asc:      true,
	})

	assert.Equal(t, "SELECT * FROM races WHERE visible = 1 ORDER BY advertised_start_time ASC", query)
}

func TestApplyOrderByDescendingWhenHavingFilter(t *testing.T) {
	rr := &racesRepo{}

	query, _ := rr.applyFilter("SELECT * FROM races", &racing.ListRacesRequestFilter{VisibleOnly: true})

	query = rr.applyOrderBy(query, &racing.ListRacesRequestOrderBy{
		Property: "visible",
		Asc:      false,
	})

	assert.Equal(t, "SELECT * FROM races WHERE visible = 1 ORDER BY visible DESC", query)
}
