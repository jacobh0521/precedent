package rules

import (
	"testing"
)

func TestSignpostRules_GiraffeSmuggling(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-giraffe",
		Name:   "A Live Giraffe",
		Traits: []string{"living", "oversized"},
	}

	// Try 1: Unmodified in a Glass Jar -> Should fail NO LIVING THINGS and NO GLASS
	badDisguise1 := Disguise{
		Adjective: "Rental",
		Container: "Glass Jar",
	}
	verdict1 := Evaluate(starterLaws, giraffe, badDisguise1, 1)
	if verdict1.Passed {
		t.Fatalf("expected bad disguise 1 to be confiscated")
	}
	if verdict1.MaxAttempts != 3 {
		t.Fatalf("expected MaxAttempts to be 3, got %d", verdict1.MaxAttempts)
	}

	// Try 2: User enters "Taxidermied" in "Matchbox" -> Should CLEAR all 5!
	// NO LIVING THINGS: Pass (taxidermied)
	// NO GLASS: Pass (matchbox)
	// MUST FIT IN A POCKET: Pass (matchbox fits in pocket)
	// NOTHING WET: Pass
	// NO METAL: Pass (cardboard matchbox)
	winningDisguise := Disguise{
		Adjective: "Taxidermied",
		Container: "Matchbox",
	}
	verdict2 := Evaluate(starterLaws, giraffe, winningDisguise, 2)
	if !verdict2.Passed {
		t.Fatalf("expected Taxidermied in Matchbox to clear all rules! audits: %+v", verdict2.Audits)
	}
	if verdict2.EmojiRow != "🟩🟩🟩🟩🟩" {
		t.Fatalf("expected emoji row 🟩🟩🟩🟩🟩, got %s", verdict2.EmojiRow)
	}
}

func TestPunchlineCompilation(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-giraffe",
		Name:   "A Live Giraffe",
		Traits: []string{"living", "oversized"},
	}

	winningDisguise := Disguise{
		Adjective: "Inflatable",
		Container: "Envelope",
	}

	newLaw := CompilePunchlineLaw(2, giraffe, winningDisguise, len(starterLaws))
	if newLaw.Signpost != "NO INFLATABLE GIRAFFES" {
		t.Fatalf("expected signpost 'NO INFLATABLE GIRAFFES', got '%s'", newLaw.Signpost)
	}

	lawbook := append(starterLaws, newLaw)

	// Re-attempting the exact same disguise on Day 2 must now FAIL Law #6!
	day2Verdict := Evaluate(lawbook, giraffe, winningDisguise, 1)
	if day2Verdict.Passed {
		t.Fatalf("expected Day 2 attempt with Inflatable Giraffe to be BUSTED by new law!")
	}

	// Alternate loophole "Robotic" + "Pocket" should pass!
	altDisguise := Disguise{
		Adjective: "Robotic",
		Container: "Pocket",
	}
	altVerdict := Evaluate(lawbook, giraffe, altDisguise, 2)
	if !altVerdict.Passed {
		t.Fatalf("expected Robotic Giraffe in Pocket to slip past day 2 laws: %+v", altVerdict.Audits)
	}
}

func TestWeatherDehydration(t *testing.T) {
	starterLaws := StarterLaws()

	weather := Cargo{
		ID:     "cargo-weather",
		Name:   "Tomorrow's Weather",
		Traits: []string{"wet", "oversized"},
	}

	good := Disguise{
		Adjective: "Dehydrated",
		Container: "Pocket",
	}
	vGood := Evaluate(starterLaws, weather, good, 2)
	if !vGood.Passed {
		t.Fatalf("expected Dehydrated weather in pocket to pass: %+v", vGood.Audits)
	}
	if vGood.MaxAttempts != 3 {
		t.Fatalf("expected max attempts 3, got %d", vGood.MaxAttempts)
	}
}
