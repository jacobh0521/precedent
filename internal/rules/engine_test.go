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
	if newLaws[0].Signpost != "NO METAPHORS OR IMAGINARY CONCEPTS" {
		t.Fatalf("expected 'NO METAPHORS OR IMAGINARY CONCEPTS', got '%s'", newLaws[0].Signpost)
	}
	if newLaws[0].RuleType != "ban_all_metaphors" {
		t.Fatalf("expected rule_type 'ban_all_metaphors', got '%s'", newLaws[0].RuleType)
	}
	if newLaws[1].Signpost != "NO PAPER OR ENVELOPE CONTAINERS" {
		t.Fatalf("expected 'NO PAPER OR ENVELOPE CONTAINERS', got '%s'", newLaws[1].Signpost)
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

func TestCreativeMutatedLoopholeRules(t *testing.T) {
	starterLaws := StarterLaws()

	// Player clears with "Microscopic" in "Envelope"
	winDisguise := Disguise{
		Adjective: "Microscopic",
		Container: "Envelope",
	}
	compiled := CompileLoopholeRules(2, winDisguise, len(starterLaws))
	if len(compiled) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(compiled))
	}

	// Law 1: Category ban on shrinking/scaling
	if compiled[0].Signpost != "NO SHRINKING OR SCALING ADJECTIVES" {
		t.Errorf("expected 'NO SHRINKING OR SCALING ADJECTIVES', got '%s'", compiled[0].Signpost)
	}
	// Law 2: Category ban on paper/envelopes
	if compiled[1].Signpost != "NO PAPER OR ENVELOPE CONTAINERS" {
		t.Errorf("expected 'NO PAPER OR ENVELOPE CONTAINERS', got '%s'", compiled[1].Signpost)
	}

	lawbook := append(starterLaws, compiled...)
	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}

	// Trying "Miniature" in "Folder" should be BUSTED by both category bans!
	vAttempt := Evaluate(lawbook, giraffe, Disguise{Adjective: "Miniature", Container: "Folder"}, 1)
	if vAttempt.Passed {
		t.Fatalf("expected Miniature in Folder to be busted by mutated category bans!")
	}

	// But switching strategy to "Inflatable" in "Velvet Pouch" succeeds!
	vCreative := Evaluate(lawbook, giraffe, Disguise{Adjective: "Inflatable", Container: "Velvet Pouch"}, 1)
	if !vCreative.Passed {
		t.Fatalf("expected Inflatable in Velvet Pouch to pass: %+v", vCreative.Audits)
	}
}

func TestCreativeContainersAndAdjectives(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}

	creativeCases := []Disguise{
		{Adjective: "Taxidermied", Container: "Pencil Case"},
		{Adjective: "Origami", Container: "Box of Cards"},
		{Adjective: "Taxidermied", Container: "Journal"},
		{Adjective: "Inflatable", Container: "Lunchbox"},
	}

	for _, d := range creativeCases {
		v := Evaluate(starterLaws, giraffe, d, 1)
		if !v.Passed {
			t.Errorf("expected creative disguise (%s, %s) to pass, but failed: %+v", d.Adjective, d.Container, v.Audits)
		}
	}
}

func TestRegimeShiftDay8(t *testing.T) {
	regime8 := GetRegimeForDay(8)
	if regime8.Number != 2 {
		t.Fatalf("expected Regime 2 on Day 8, got %d", regime8.Number)
	}
	if len(regime8.Codex) != 5 {
		t.Fatalf("expected 5 starter codex laws in Regime 2, got %d", len(regime8.Codex))
	}
	if !IsCycleStartDay(8) {
		t.Fatalf("expected IsCycleStartDay(8) to be true")
	}
	if IsCycleStartDay(7) {
		t.Fatalf("expected IsCycleStartDay(7) to be false")
	}

	// Test Existential Dread under Regime 2
	dread := Cargo{
		ID:     "cargo-day-8",
		Name:   "EXISTENTIAL DREAD",
		Traits: []string{"sentient", "cursed", "oversized"},
	}

	// Fails with soft pocket / sentient form
	vFail := Evaluate(regime8.Codex, dread, Disguise{Adjective: "Unmodified", Container: "Pocket"}, 1)
	if vFail.Passed {
		t.Fatalf("expected Existential Dread with Unmodified in Pocket to fail Regime 2!")
	}

	// Passes with Taxidermied / Pacified in Rigid Lockbox / Tin Can
	vPass := Evaluate(regime8.Codex, dread, Disguise{Adjective: "Taxidermied", Container: "Tin Can"}, 1)
	if !vPass.Passed {
		t.Fatalf("expected Taxidermied in Tin Can to pass Regime 2: %+v", vPass.Audits)
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

func TestFormatDisguiseSentenceGrammar(t *testing.T) {
	cases := []struct {
		adj       string
		cargo     string
		container string
		expected  string
	}{
		{"Metaphorical", "A Live Giraffe", "Envelope", "Metaphorical Live Giraffe in an Envelope"},
		{"Inflatable", "A Live Giraffe", "Plastic Bag", "Inflatable Live Giraffe in a Plastic Bag"},
		{"Dehydrated", "Tomorrow's Weather", "Pocket", "Dehydrated Tomorrow's Weather in a Pocket"},
		{"Taxidermied", "An Angry Goose", "Matchbox", "Taxidermied Angry Goose in a Matchbox"},
		{"Unmodified", "A Live Giraffe", "Bare Hands", "Live Giraffe with bare hands"},
	}

	for _, tc := range cases {
		actual := FormatDisguiseSentence(tc.adj, tc.cargo, tc.container)
		if actual != tc.expected {
			t.Errorf("FormatDisguiseSentence(%q, %q, %q) = %q; expected %q", tc.adj, tc.cargo, tc.container, actual, tc.expected)
		}
	}
}
