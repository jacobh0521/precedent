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
	newLaws := CompileLoopholeRules(2, winningDisguise, starterLaws)

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

// Issue 1: Adjectives explicitly remove conflicting base traits
func TestAdjectivesExplicitlyRemoveConflictingBaseTraits(t *testing.T) {
	giraffe := Cargo{
		ID:     "cargo-giraffe",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized", "flammable"},
	}

	// 1. Cremated: removes living, wet, oversized, flammable
	traitsCremated := ResolveCargoTraits(giraffe, "Cremated")
	if traitsCremated["living"] {
		t.Errorf("Cremated must remove living trait")
	}
	if traitsCremated["oversized"] {
		t.Errorf("Cremated must remove oversized trait")
	}
	if traitsCremated["flammable"] {
		t.Errorf("Cremated must remove flammable trait")
	}
	if !traitsCremated["dry"] || !traitsCremated["ash"] {
		t.Errorf("Cremated must add dry and ash traits")
	}

	// 2. Dehydrated: removes wet, liquid
	weather := Cargo{
		ID:     "cargo-weather",
		Name:   "TOMORROW'S WEATHER",
		Traits: []string{"wet", "liquid", "oversized"},
	}
	traitsDehydrated := ResolveCargoTraits(weather, "Dehydrated")
	if traitsDehydrated["wet"] || traitsDehydrated["liquid"] {
		t.Errorf("Dehydrated must remove wet and liquid traits")
	}
	if !traitsDehydrated["dry"] {
		t.Errorf("Dehydrated must add dry trait")
	}

	// 3. Inflatable: removes living, sentient, predator, oversized
	traitsInflatable := ResolveCargoTraits(giraffe, "Inflatable")
	if traitsInflatable["living"] {
		t.Errorf("Inflatable must remove living trait")
	}
	if traitsInflatable["oversized"] {
		t.Errorf("Inflatable must remove oversized trait (deflated/pocketable)")
	}
	if !traitsInflatable["toy"] {
		t.Errorf("Inflatable must add toy trait")
	}

	// 4. Metaphorical: removes all physical traits
	traitsMeta := ResolveCargoTraits(giraffe, "Metaphorical")
	if traitsMeta["living"] || traitsMeta["oversized"] || traitsMeta["flammable"] {
		t.Errorf("Metaphorical must remove all physical constraints")
	}

	// 5. Frozen: removes liquid, wet, gas; adds solid, dry
	traitsFrozen := ResolveCargoTraits(weather, "Frozen")
	if traitsFrozen["liquid"] || traitsFrozen["wet"] {
		t.Errorf("Frozen must remove liquid and wet traits")
	}
	if !traitsFrozen["solid"] || !traitsFrozen["dry"] {
		t.Errorf("Frozen must add solid and dry traits")
	}
}

// Issue 2: MUST FIT IN A POCKET fails if cargo still has large/tall/oversized trait even in pocket container
func TestMustFitInPocketFailsIfCargoStillOversized(t *testing.T) {
	starterLaws := StarterLaws()

	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}

	// 1. Unmodified inside Pocket -> MUST FAIL (a 5.5m live giraffe cannot fit in a pocket!)
	vUnmodified := Evaluate(starterLaws, giraffe, Disguise{Adjective: "Unmodified", Container: "Pocket"}, 1)
	if vUnmodified.Passed {
		t.Fatalf("expected Unmodified Giraffe in Pocket to FAIL, but passed!")
	}
	foundPocketFail := false
	for _, a := range vUnmodified.Audits {
		if a.Signpost == "MUST FIT IN A POCKET" && !a.Passed {
			foundPocketFail = true
		}
	}
	if !foundPocketFail {
		t.Errorf("expected MUST FIT IN A POCKET to fail for Unmodified Giraffe in Pocket")
	}

	// 2. Taxidermied inside Pocket -> Living is removed, BUT giraffe is still oversized -> MUST FAIL MUST FIT IN A POCKET!
	vTaxidermied := Evaluate(starterLaws, giraffe, Disguise{Adjective: "Taxidermied", Container: "Pocket"}, 1)
	if vTaxidermied.Passed {
		t.Fatalf("expected full-sized Taxidermied Giraffe in Pocket to FAIL MUST FIT IN A POCKET!")
	}
	foundPocketFailTaxi := false
	for _, a := range vTaxidermied.Audits {
		if a.Signpost == "MUST FIT IN A POCKET" && !a.Passed {
			foundPocketFailTaxi = true
		}
	}
	if !foundPocketFailTaxi {
		t.Errorf("expected MUST FIT IN A POCKET to fail for Taxidermied Giraffe in Pocket")
	}

	// 3. Microscopic inside Pocket -> Oversized trait is removed by shrinking -> PASSES MUST FIT IN A POCKET!
	lawPocket := Law{ID: 3, Day: 1, Signpost: "MUST FIT IN A POCKET", RuleType: "pocket_only"}
	if !EvaluateRule(lawPocket, giraffe, Disguise{Adjective: "Microscopic", Container: "Pocket"}) {
		t.Errorf("expected Microscopic Giraffe in Pocket to pass pocket_only")
	}

	// 4. Cremated inside Pocket -> Ash fits in pocket -> PASSES MUST FIT IN A POCKET!
	if !EvaluateRule(lawPocket, giraffe, Disguise{Adjective: "Cremated", Container: "Pocket"}) {
		t.Errorf("expected Cremated Giraffe in Pocket to pass pocket_only")
	}

	// 5. Inflatable inside Pocket -> Deflated toy fits in pocket -> PASSES MUST FIT IN A POCKET!
	if !EvaluateRule(lawPocket, giraffe, Disguise{Adjective: "Inflatable", Container: "Pocket"}) {
		t.Errorf("expected Inflatable Giraffe in Pocket to pass pocket_only")
	}
}

// Issue 3: Clean up manifest display names so adjectives replace "LIVE"
func TestFormatDisguiseSentenceReplacesLive(t *testing.T) {
	cases := []struct {
		adj       string
		cargo     string
		container string
		expected  string
	}{
		// Adjective replaces "LIVE"
		{"Cremated", "A LIVE GIRAFFE", "Envelope", "Cremated Giraffe in an Envelope"},
		{"Inflatable", "A LIVE GIRAFFE", "Pocket", "Inflatable Giraffe in a Pocket"},
		{"Metaphorical", "A LIVE GIRAFFE", "Envelope", "Metaphorical Giraffe in an Envelope"},
		{"Taxidermied", "A LIVE GIRAFFE", "Matchbox", "Taxidermied Giraffe in a Matchbox"},
		{"Frozen", "A LIVE GIRAFFE", "Cardboard Box", "Frozen Giraffe in a Cardboard Box"},

		// Unmodified retains "Live Giraffe"
		{"Unmodified", "A LIVE GIRAFFE", "Bare Hands", "Live Giraffe with bare hands"},
		{"", "A LIVE GIRAFFE", "Bare Hands", "Live Giraffe with bare hands"},

		// Items without "LIVE" format normally
		{"Dehydrated", "TOMORROW'S WEATHER", "Pocket", "Dehydrated Tomorrow's Weather in a Pocket"},
		{"Taxidermied", "AN ANGRY GOOSE", "Matchbox", "Taxidermied Angry Goose in a Matchbox"},
	}

	for _, tc := range cases {
		actual := FormatDisguiseSentence(tc.adj, tc.cargo, tc.container)
		if actual != tc.expected {
			t.Errorf("FormatDisguiseSentence(%q, %q, %q) = %q; expected %q", tc.adj, tc.cargo, tc.container, actual, tc.expected)
		}
	}
}

// Issue 4: Max attempts constant synchronization
func TestMaxDailyAttemptsConstant(t *testing.T) {
	if MaxDailyAttempts != 3 {
		t.Fatalf("expected MaxDailyAttempts constant to equal 3, got %d", MaxDailyAttempts)
	}

	starterLaws := StarterLaws()
	giraffe := Cargo{ID: "c1", Name: "Giraffe", Traits: []string{"living"}}
	verdict := Evaluate(starterLaws, giraffe, Disguise{Adjective: "Unmodified", Container: "Pocket"}, 1)
	if verdict.MaxAttempts != MaxDailyAttempts {
		t.Errorf("expected verdict.MaxAttempts (%d) to equal MaxDailyAttempts (%d)", verdict.MaxAttempts, MaxDailyAttempts)
	}
}

// Issue 5: Solvability check guarantees at least 1 winning combo remains for tomorrow's cargo
func TestSolvabilityCheckGuaranteesWinningCombo(t *testing.T) {
	starterLaws := StarterLaws()

	// Today: Giraffe cleared with Folded in Envelope
	winDisguise := Disguise{
		Adjective: "Folded",
		Container: "Envelope",
	}

	// Tomorrow: Tomorrow's Weather (wet, liquid, oversized)
	tomorrowCargo := Cargo{
		ID:     "cargo-2",
		Name:   "TOMORROW'S WEATHER",
		Traits: []string{"wet", "liquid", "oversized"},
	}

	compiled := CompileLoopholeRules(2, winDisguise, starterLaws, tomorrowCargo)
	if len(compiled) == 0 {
		t.Fatalf("expected compiled laws, got 0")
	}

	// Verify that tomorrow's cargo is SOLVABLE under the combined lawbook!
	combinedLaws := append(starterLaws, compiled...)
	solvable := CheckSolvability(combinedLaws, tomorrowCargo)
	if !solvable {
		t.Fatalf("expected tomorrow's cargo to remain solvable after CompileLoopholeRules, but got 0 winning combos!")
	}
}

func TestCreativeMutatedLoopholeRules(t *testing.T) {
	starterLaws := StarterLaws()

	// Player clears with "Microscopic" in "Envelope"
	winDisguise := Disguise{
		Adjective: "Microscopic",
		Container: "Envelope",
	}
	compiled := CompileLoopholeRules(2, winDisguise, starterLaws)
	if len(compiled) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(compiled))
	}

	if compiled[0].Signpost != "NO SHRINKING OR SCALING ADJECTIVES" {
		t.Errorf("expected 'NO SHRINKING OR SCALING ADJECTIVES', got '%s'", compiled[0].Signpost)
	}
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
		{Adjective: "Origami", Container: "Pencil Case"},
		{Adjective: "Origami", Container: "Box of Cards"},
		{Adjective: "Cremated", Container: "Journal"},
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

	dread := Cargo{
		ID:     "cargo-day-8",
		Name:   "EXISTENTIAL DREAD",
		Traits: []string{"sentient", "cursed", "oversized"},
	}

	vFail := Evaluate(regime8.Codex, dread, Disguise{Adjective: "Unmodified", Container: "Pocket"}, 1)
	if vFail.Passed {
		t.Fatalf("expected Existential Dread with Unmodified in Pocket to fail Regime 2!")
	}

	vPass := Evaluate(regime8.Codex, dread, Disguise{Adjective: "Taxidermied", Container: "Tin Can"}, 1)
	if !vPass.Passed {
		t.Fatalf("expected Taxidermied in Tin Can to pass Regime 2: %+v", vPass.Audits)
	}
}

func TestFairCargoEvaluation(t *testing.T) {
	starterLaws := StarterLaws()

	weather := Cargo{
		ID:     "cargo-2",
		Name:   "TOMORROW'S WEATHER",
		Traits: []string{"wet", "oversized"},
	}

	// Solvable with Dehydrated in Pocket (dehydrated removes wet; if small/packed)
	// Note: Weather is wet and oversized. Dehydrated removes wet. To fit in pocket, compressed/dehydrated powder fits.
	vDehydrated := Evaluate(starterLaws, weather, Disguise{Adjective: "Freeze-Dried", Container: "Pocket"}, 1)
	if !vDehydrated.Passed {
		t.Fatalf("expected Freeze-Dried in Pocket to pass: %+v", vDehydrated.Audits)
	}

	// Solvable fairly with Inflatable in Matchbox for Giraffe
	giraffe := Cargo{
		ID:     "cargo-1",
		Name:   "A LIVE GIRAFFE",
		Traits: []string{"living", "oversized"},
	}
	vInflatable := Evaluate(starterLaws, giraffe, Disguise{Adjective: "Inflatable", Container: "Matchbox"}, 1)
	if !vInflatable.Passed {
		t.Fatalf("expected Inflatable in Matchbox to pass: %+v", vInflatable.Audits)
	}

	// Solvable with Cremated in Pocket
	vCremated := Evaluate(starterLaws, giraffe, Disguise{Adjective: "Cremated", Container: "Pocket"}, 1)
	if !vCremated.Passed {
		t.Fatalf("expected Cremated in Pocket to pass: %+v", vCremated.Audits)
	}
}
