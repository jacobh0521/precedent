package rules

import "fmt"

// Regime represents a 7-day customs authority cycle with custom starter laws.
type Regime struct {
	Number int
	Name   string
	Codex  []Law
}

// StarterLaws returns the original 5 blunt signpost starter rules for Regime 1.
func StarterLaws() []Law {
	return []Law{
		{
			ID:       1,
			Day:      1,
			Signpost: "NO LIVING THINGS",
			RuleType: "no_living",
		},
		{
			ID:       2,
			Day:      1,
			Signpost: "NO GLASS",
			RuleType: "no_glass",
		},
		{
			ID:       3,
			Day:      1,
			Signpost: "MUST FIT IN A POCKET",
			RuleType: "pocket_only",
		},
		{
			ID:       4,
			Day:      1,
			Signpost: "NOTHING WET",
			RuleType: "no_wet",
		},
		{
			ID:       5,
			Day:      1,
			Signpost: "NO METAL",
			RuleType: "no_metal",
		},
	}
}

// GetRegimeForDay calculates the 7-day regime for any calendar day.
func GetRegimeForDay(day int) Regime {
	if day < 1 {
		day = 1
	}
	cycle := ((day - 1) / 7) + 1

	switch cycle {
	case 1:
		return Regime{
			Number: 1,
			Name:   "SURFACE BORDER PATROL",
			Codex:  StarterLaws(),
		}
	case 2:
		return Regime{
			Number: 2,
			Name:   "BIO & ESOTERIC CUSTOMS REGIME",
			Codex: []Law{
				{ID: 1, Day: day, Signpost: "NO LIQUIDS OR GASES", RuleType: "no_liquid_gas"},
				{ID: 2, Day: day, Signpost: "NO SENTIENT OR CURSED", RuleType: "no_sentient_cursed"},
				{ID: 3, Day: day, Signpost: "CONTAINERS MUST BE RIGID", RuleType: "rigid_container"},
				{ID: 4, Day: day, Signpost: "NO FLAMMABLE MATERIALS", RuleType: "no_flammable"},
				{ID: 5, Day: day, Signpost: "NO GLASS OR CRYSTAL", RuleType: "no_glass"},
			},
		}
	case 3:
		return Regime{
			Number: 3,
			Name:   "SYNTHETIC HIGH-SEC BORDER",
			Codex: []Law{
				{ID: 1, Day: day, Signpost: "NO ORGANIC CONTAINERS", RuleType: "no_organic_container"},
				{ID: 2, Day: day, Signpost: "CONTAINERS MUST BE TRANSPARENT", RuleType: "transparent_container"},
				{ID: 3, Day: day, Signpost: "NO WEAPONS OR PREDATORS", RuleType: "no_predators"},
				{ID: 4, Day: day, Signpost: "NOTHING RADIOACTIVE", RuleType: "no_radiation"},
				{ID: 5, Day: day, Signpost: "MUST DECLARE A REAL CONTAINER", RuleType: "must_use_container"},
			},
		}
	default:
		mod := (cycle - 1) % 3
		if mod == 0 {
			return Regime{
				Number: cycle,
				Name:   fmt.Sprintf("REGIME %d // SURFACE RE-INSPECTION", cycle),
				Codex:  StarterLaws(),
			}
		} else if mod == 1 {
			return Regime{
				Number: cycle,
				Name:   fmt.Sprintf("REGIME %d // BIO & ESOTERIC OVERHAUL", cycle),
				Codex: []Law{
					{ID: 1, Day: day, Signpost: "NO LIQUIDS OR GASES", RuleType: "no_liquid_gas"},
					{ID: 2, Day: day, Signpost: "NO SENTIENT OR CURSED", RuleType: "no_sentient_cursed"},
					{ID: 3, Day: day, Signpost: "CONTAINERS MUST BE RIGID", RuleType: "rigid_container"},
					{ID: 4, Day: day, Signpost: "NO FLAMMABLE MATERIALS", RuleType: "no_flammable"},
					{ID: 5, Day: day, Signpost: "NO GLASS OR CRYSTAL", RuleType: "no_glass"},
				},
			}
		} else {
			return Regime{
				Number: cycle,
				Name:   fmt.Sprintf("REGIME %d // SYNTHETIC HIGH-SEC OVERHAUL", cycle),
				Codex: []Law{
					{ID: 1, Day: day, Signpost: "NO ORGANIC CONTAINERS", RuleType: "no_organic_container"},
					{ID: 2, Day: day, Signpost: "CONTAINERS MUST BE TRANSPARENT", RuleType: "transparent_container"},
					{ID: 3, Day: day, Signpost: "NO WEAPONS OR PREDATORS", RuleType: "no_predators"},
					{ID: 4, Day: day, Signpost: "NOTHING RADIOACTIVE", RuleType: "no_radiation"},
					{ID: 5, Day: day, Signpost: "MUST DECLARE A REAL CONTAINER", RuleType: "must_use_container"},
				},
			}
		}
	}
}

// IsCycleStartDay returns true on Day 8, Day 15, Day 22, etc.
func IsCycleStartDay(day int) bool {
	return day > 1 && (day-1)%7 == 0
}
