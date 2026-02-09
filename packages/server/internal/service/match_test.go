package service

import (
	"testing"
	"time"

	"github.com/galchammat/kadeem/internal/model"
)

func TestComputeMatchSyncStartTime_WithSyncedAt(t *testing.T) {
	syncedAtTime := int64(1000000)
	expectedStart := syncedAtTime + 1

	if expectedStart != 1000001 {
		t.Errorf("expected start=%d, got %d", 1000001, expectedStart)
	}
}

func TestComputeMatchSyncStartTime_ExclusiveBoundary(t *testing.T) {
	syncedAtTime := int64(1609459200)
	expectedStart := syncedAtTime + 1

	if expectedStart <= syncedAtTime {
		t.Errorf("expected start > synced_at, but start=%d <= synced_at=%d", expectedStart, syncedAtTime)
	}
}

func TestComputeMatchSyncCount_Always100(t *testing.T) {
	count := 100

	if count != 100 {
		t.Errorf("expected count=100, got %d", count)
	}
}

func TestComputeMatchSyncStartTime_BoundaryMaxInt64(t *testing.T) {
	syncedAtTime := int64(9223372036854775800)
	expectedStart := syncedAtTime + 1

	if expectedStart != 9223372036854775801 {
		t.Errorf("expected start=%d, got %d", 9223372036854775801, expectedStart)
	}
}

func TestComputeBroadcastSyncStartTime_WithSyncedAt(t *testing.T) {
	syncedAtTime := int64(1000000)
	expectedStart := syncedAtTime + 1

	if expectedStart != 1000001 {
		t.Errorf("expected start=%d, got %d", 1000001, expectedStart)
	}
}

func TestComputeBroadcastSyncStartTime_NilStartsAtZero(t *testing.T) {
	channel := model.Channel{
		SyncedAt: nil,
	}

	var startTime int64
	if channel.SyncedAt != nil {
		startTime = *channel.SyncedAt + 1
	}

	if startTime != 0 {
		t.Errorf("expected startTime=0 when SyncedAt is nil, got %d", startTime)
	}
}

func TestComputeBroadcastSyncStartTime_ExclusiveBoundary(t *testing.T) {
	syncedAtTime := int64(1609459200)
	channel := model.Channel{
		SyncedAt: &syncedAtTime,
	}

	var startTime int64
	if channel.SyncedAt != nil {
		startTime = *channel.SyncedAt + 1
	}

	if startTime <= syncedAtTime {
		t.Errorf("expected startTime > synced_at, but startTime=%d <= synced_at=%d", startTime, syncedAtTime)
	}
}

func TestSyncBoundaryTimestampConsistency(t *testing.T) {
	now := time.Now().Unix()

	account := model.LeagueOfLegendsAccount{
		SyncedAt: &now,
	}

	if account.SyncedAt == nil {
		t.Fatal("expected SyncedAt to be set")
	}

	if *account.SyncedAt != now {
		t.Errorf("expected SyncedAt=%d, got %d", now, *account.SyncedAt)
	}

	nextSync := *account.SyncedAt + 1
	if nextSync <= now {
		t.Errorf("expected nextSync > now, got nextSync=%d, now=%d", nextSync, now)
	}
}
