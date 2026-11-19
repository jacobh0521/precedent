package rules

import (
	"testing"
)

func TestSignpostRules_GiraffeSmuggling(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
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

	// Try 2: User enters "Metaphorical" in "Envelope" -> Should CLEAR all 5!
	winningDisguise := Disguise{
		Adjective: "Metaphorical",
		Container: "Envelope",
	}
	verdict2 := Evaluate(starterLaws, giraffe, winningDisguise, 2)
	if !verdict2.Passed {
		t.Fatalf("expected Metaphorical in Envelope to clear all rules! audits: %+v", verdict2.Audits)
	}
	if verdict2.EmojiRow != "🟩🟩🟩🟩🟩" {
		t.Fatalf("expected emoji row 🟩🟩🟩🟩🟩, got %s", verdict2.EmojiRow)
	}
}

func TestTwoRuleLoopholeCompilation(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}

	winningDisguise := Disguise{
		Adjective: "Metaphorical",
		Container: "Envelope",
	}

	// Day 1 clearance generates TWO rules for Day 2:
	// Rule 1: NOTHING METAPHORICAL
	// Rule 2: NOTHING IN AN ENVELOPE
	newLaws := CompileLoopholeRules(2, winningDisguise, len(starterLaws))
	if len(newLaws) != 2 {
		t.Fatalf("expected 2 new laws, got %d", len(newLaws))
	}
	if newLaws[0].Signpost != "NOTHING METAPHORICAL" {
		t.Fatalf("expected 'NOTHING METAPHORICAL', got '%s'", newLaws[0].Signpost)
	}
	if newLaws[1].Signpost != "NOTHING IN AN ENVELOPE" {
		t.Fatalf("expected 'NOTHING IN AN ENVELOPE', got '%s'", newLaws[1].Signpost)
	}

	lawbook := append(starterLaws, newLaws...)
	if len(lawbook) != 7 {
		t.Fatalf("expected 7 laws in Day 2 lawbook, got %d", len(lawbook))
	}

	// Test Day 2: Trying Metaphorical in Pocket must fail Law #6 (NOTHING METAPHORICAL)
	day2FailAdj := Evaluate(lawbook, giraffe, Disguise{Adjective: "Metaphorical", Container: "Pocket"}, 1)
	if day2FailAdj.Passed {
		t.Fatalf("expected Metaphorical to be BUSTED by Law #6")
	}

	// Test Day 2: Trying Inflatable in Envelope must fail Law #7 (NOTHING IN AN ENVELOPE)
	day2FailCnt := Evaluate(lawbook, giraffe, Disguise{Adjective: "Inflatable", Container: "Envelope"}, 2)
	if day2FailCnt.Passed {
		t.Fatalf("expected Envelope to be BUSTED by Law #7")
	}

	// Test Day 2: Trying Inflatable in Pocket clears all 7 laws!
	day2Success := Evaluate(lawbook, giraffe, Disguise{Adjective: "Inflatable", Container: "Pocket"}, 3)
	if !day2Success.Passed {
		t.Fatalf("expected Inflatable in Pocket to clear all 7 laws! audits: %+v", day2Success.Audits)
	}
	if day2Success.EmojiRow != "🟩🟩🟩🟩🟩🟩🟩" {
		t.Fatalf("expected 7 green squares, got %s", day2Success.EmojiRow)
	}
}
