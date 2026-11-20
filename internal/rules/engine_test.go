package rules

import (
	"testing"
)

func TestMetaphoricalWorks_TheoreticalFails(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}

	// 1. "Theoretical" DOES NOT work!
	theoreticalDisguise := Disguise{
		Adjective: "Theoretical",
		Container: "Envelope",
	}
	verdictTheo := Evaluate(starterLaws, giraffe, theoreticalDisguise, 1)
	if verdictTheo.Passed {
		t.Fatalf("expected Theoretical to FAIL, but it passed!")
	}

	// Verify Theoretical specifically fails NO LIVING THINGS
	foundLivingFail := false
	for _, a := range verdictTheo.Audits {
		if a.Signpost == "NO LIVING THINGS" && !a.Passed {
			foundLivingFail = true
		}
	}
	if !foundLivingFail {
		t.Errorf("expected Theoretical Giraffe to be BUSTED under NO LIVING THINGS")
	}

	// 2. "Metaphorical" WORKS as get-out-of-jail card
	metaDisguise := Disguise{
		Adjective: "Metaphorical",
		Container: "Envelope",
	}
	verdictMeta := Evaluate(starterLaws, giraffe, metaDisguise, 2)
	if !verdictMeta.Passed {
		t.Fatalf("expected Metaphorical in Envelope to CLEAR, but failed: %+v", verdictMeta.Audits)
	}

	// 3. "Imaginary" also works as a synonym on Day 1
	imaginaryDisguise := Disguise{
		Adjective: "Imaginary",
		Container: "Pocket",
	}
	verdictImag := Evaluate(starterLaws, giraffe, imaginaryDisguise, 2)
	if !verdictImag.Passed {
		t.Fatalf("expected Imaginary in Pocket to CLEAR on Day 1, but failed: %+v", verdictImag.Audits)
	}
}

func TestMetaphoricalSynonymBlanketBan(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}

	// Player clears Day 1 using "Metaphorical" in "Envelope"
	winningDisguise := Disguise{
		Adjective: "Metaphorical",
		Container: "Envelope",
	}
	newLaws := CompileLoopholeRules(2, winningDisguise, len(starterLaws))

	if len(newLaws) != 2 {
		t.Fatalf("expected 2 new laws, got %d", len(newLaws))
	}
	// Verify that law 1 is the blanket ban on ALL metaphors
	if newLaws[0].Signpost != "NO METAPHORS OR IMAGINARY CONCEPTS" {
		t.Fatalf("expected 'NO METAPHORS OR IMAGINARY CONCEPTS', got '%s'", newLaws[0].Signpost)
	}
	if newLaws[0].RuleType != "ban_all_metaphors" {
		t.Fatalf("expected rule_type 'ban_all_metaphors', got '%s'", newLaws[0].RuleType)
	}
	if newLaws[1].Signpost != "NOTHING IN AN ENVELOPE" {
		t.Fatalf("expected 'NOTHING IN AN ENVELOPE', got '%s'", newLaws[1].Signpost)
	}

	day2Lawbook := append(starterLaws, newLaws...)

	// On Day 2, trying another synonym like "Imaginary" or "Allegorical" or "Fictional" MUST BE BUSTED!
	synonymsToTest := []string{"Imaginary", "Allegorical", "Symbolic", "Fictional", "Abstract", "Metaphorical"}
	for _, syn := range synonymsToTest {
		v := Evaluate(day2Lawbook, giraffe, Disguise{Adjective: syn, Container: "Pocket"}, 1)
		if v.Passed {
			t.Fatalf("expected synonym '%s' to be BUSTED on Day 2 by blanket ban, but it passed!", syn)
		}
		foundBlanketFail := false
		for _, a := range v.Audits {
			if a.Signpost == "NO METAPHORS OR IMAGINARY CONCEPTS" && !a.Passed {
				foundBlanketFail = true
			}
		}
		if !foundBlanketFail {
			t.Errorf("expected failure on NO METAPHORS OR IMAGINARY CONCEPTS for synonym '%s'", syn)
		}
	}

	// But a concrete disguise like "Inflatable" in "Pocket" STILL PASSES!
	vConcrete := Evaluate(day2Lawbook, giraffe, Disguise{Adjective: "Inflatable", Container: "Pocket"}, 2)
	if !vConcrete.Passed {
		t.Fatalf("expected Inflatable in Pocket to clear Day 2: %+v", vConcrete.Audits)
	}
}

func TestFairCargoEvaluation(t *testing.T) {
	starterLaws := StarterLaws()

	// Contraband with wet trait
	weather := Cargo{
		ID:     "cargo-2",
		Name:   "TOMORROW'S WEATHER",
		Traits: []string{"wet", "oversized"},
	}

	// Solvable fairly with Dehydrated in Pocket
	vDehydrated := Evaluate(starterLaws, weather, Disguise{Adjective: "Dehydrated", Container: "Pocket"}, 1)
	if !vDehydrated.Passed {
		t.Fatalf("expected Dehydrated in Pocket to pass: %+v", vDehydrated.Audits)
	}

	// Solvable fairly with Ziploc Bag (waterproof container)
	vZiploc := Evaluate(starterLaws, weather, Disguise{Adjective: "Unmodified", Container: "Ziploc Bag"}, 1)
	if !vZiploc.Passed {
		t.Fatalf("expected Ziploc Bag to pass for wet cargo: %+v", vZiploc.Audits)
	}

	// Contraband with living trait: solvable with Taxidermied in Matchbox
	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}
	vTaxi := Evaluate(starterLaws, giraffe, Disguise{Adjective: "Taxidermied", Container: "Matchbox"}, 1)
	if !vTaxi.Passed {
		t.Fatalf("expected Taxidermied in Matchbox to pass: %+v", vTaxi.Audits)
	}
}
