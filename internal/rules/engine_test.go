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

	for _, a := range verdict1.Audits {
		if a.Signpost == "NO LIVING THINGS" && a.Passed {
			t.Errorf("expected NO LIVING THINGS to fail on rental live giraffe")
		}
		if a.Signpost == "NO GLASS" && a.Passed {
			t.Errorf("expected NO GLASS to fail on Glass Jar")
		}
	}

	// Try 2: Winning loophole: "Inflatable" + "Envelope"
	// NO LIVING THINGS: Pass (inflatable!)
	// NO GLASS: Pass (envelope)
	// MUST FIT IN A POCKET: Pass (inflatable & envelope!)
	// NOTHING WET: Pass
	// NO METAL: Pass
	winningDisguise := Disguise{
		Adjective: "Inflatable",
		Container: "Envelope",
	}
	verdict2 := Evaluate(starterLaws, giraffe, winningDisguise, 2)
	if !verdict2.Passed {
		t.Fatalf("expected winning disguise to clear all rules! audits: %+v", verdict2.Audits)
	}

	for _, a := range verdict2.Audits {
		if !a.Passed {
			t.Errorf("rule %s should have passed, but failed", a.Signpost)
		}
		if a.Status != "SLIPPED PAST" {
			t.Errorf("expected SLIPPED PAST status, got %s", a.Status)
		}
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

	// Compile punchline law
	newLaw := CompilePunchlineLaw(2, giraffe, winningDisguise, len(starterLaws))
	if newLaw.Signpost != "NO INFLATABLE GIRAFFES" {
		t.Fatalf("expected signpost 'NO INFLATABLE GIRAFFES', got '%s'", newLaw.Signpost)
	}

	// Add to lawbook
	lawbook := append(starterLaws, newLaw)

	// Re-attempting the exact same disguise on Day 2 must now FAIL Law #6!
	day2Verdict := Evaluate(lawbook, giraffe, winningDisguise, 1)
	if day2Verdict.Passed {
		t.Fatalf("expected Day 2 attempt with Inflatable Giraffe to be BUSTED by new law!")
	}

	foundPunchlineFail := false
	for _, a := range day2Verdict.Audits {
		if a.Signpost == "NO INFLATABLE GIRAFFES" && !a.Passed {
			foundPunchlineFail = true
		}
	}
	if !foundPunchlineFail {
		t.Fatalf("expected audit to report BUSTED for NO INFLATABLE GIRAFFES")
	}

	// But an alternate loophole like "Metaphorical" + "Envelope" should pass!
	altDisguise := Disguise{
		Adjective: "Metaphorical",
		Container: "Envelope",
	}
	altVerdict := Evaluate(lawbook, giraffe, altDisguise, 2)
	if !altVerdict.Passed {
		t.Fatalf("expected Metaphorical Giraffe to slip past day 2 laws: %+v", altVerdict.Audits)
	}
}

func TestWeatherDehydration(t *testing.T) {
	starterLaws := StarterLaws()

	weather := Cargo{
		ID:     "cargo-weather",
		Name:   "Tomorrow's Weather",
		Traits: []string{"wet", "oversized"},
	}

	// Normal weather fails NOTHING WET
	bad := Disguise{
		Adjective: "Ceremonial",
		Container: "Envelope",
	}
	vBad := Evaluate(starterLaws, weather, bad, 1)
	if vBad.Passed {
		t.Fatalf("expected wet weather to fail")
	}

	// Dehydrated weather in Pocket slips past all 5!
	good := Disguise{
		Adjective: "Dehydrated",
		Container: "Pocket",
	}
	vGood := Evaluate(starterLaws, weather, good, 2)
	if !vGood.Passed {
		t.Fatalf("expected Dehydrated weather in pocket to pass: %+v", vGood.Audits)
	}

	punchline := CompilePunchlineLaw(2, weather, good, 5)
	if punchline.Signpost != "NO DEHYDRATED WEATHER" {
		t.Fatalf("expected 'NO DEHYDRATED WEATHER', got '%s'", punchline.Signpost)
	}
}
