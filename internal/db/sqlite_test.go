package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreDay8RegimeShift(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "precedent-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Initial day should be 1 with 5 starter laws
	day, err := store.GetCurrentDay()
	if err != nil || day != 1 {
		t.Fatalf("expected day 1, got %d (err: %v)", day, err)
	}
	laws, err := store.GetActiveLaws()
	if err != nil || len(laws) != 5 {
		t.Fatalf("expected 5 starter laws, got %d (err: %v)", len(laws), err)
	}

	// Advance through days 2 to 7, simulating precedent laws adding up
	for d := 1; d <= 6; d++ {
		nextD, err := store.AdvanceToNextDay()
		if err != nil {
			t.Fatalf("failed to advance day: %v", err)
		}
		if nextD != d+1 {
			t.Fatalf("expected day %d, got %d", d+1, nextD)
		}
	}

	// Now on Day 7, let's verify we have advanced
	day7, _ := store.GetCurrentDay()
	if day7 != 7 {
		t.Fatalf("expected Day 7, got %d", day7)
	}

	// Advance from Day 7 to Day 8 -> Trigger Regime Shift!
	day8, err := store.AdvanceToNextDay()
	if err != nil {
		t.Fatalf("failed to advance to day 8: %v", err)
	}
	if day8 != 8 {
		t.Fatalf("expected Day 8, got %d", day8)
	}

	// On Day 8, all past precedent laws should be purged and exactly 5 Regime 2 codex laws active
	activeLaws8, err := store.GetActiveLaws()
	if err != nil {
		t.Fatalf("failed to get active laws on day 8: %v", err)
	}
	if len(activeLaws8) != 5 {
		t.Fatalf("expected exactly 5 active laws on day 8 regime overhaul, got %d", len(activeLaws8))
	}

	// Verify Regime 2 starter signposts
	foundLiquid := false
	for _, l := range activeLaws8 {
		if l.Signpost == "NO LIQUIDS OR GASES" {
			foundLiquid = true
		}
	}
	if !foundLiquid {
		t.Errorf("expected 'NO LIQUIDS OR GASES' in Regime 2 active laws, found: %+v", activeLaws8)
	}
}
