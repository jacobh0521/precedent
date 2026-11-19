package rules

import (
	"fmt"
	"strings"
	"time"
)

// Cargo represents the contraband item.
type Cargo struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Traits []string `json:"traits"` // e.g., "living", "wet", "oversized"
}

// HasTrait checks if cargo has a trait (case-insensitive).
func (c Cargo) HasTrait(t string) bool {
	target := strings.ToLower(strings.TrimSpace(t))
	for _, tr := range c.Traits {
		if strings.ToLower(strings.TrimSpace(tr)) == target {
			return true
		}
	}
	return false
}

// Disguise is the 2-click player declaration: 1 Adjective + 1 Container.
type Disguise struct {
	Adjective string `json:"adjective"` // e.g. "Inflatable", "Dehydrated", "Folded", "Metaphorical"
	Container string `json:"container"` // e.g. "Envelope", "Pocket", "Glass Jar", "Lead Crate", "Thermos"
}

// Law represents a 4-word signpost rule.
type Law struct {
	ID          int64     `json:"id"`
	Day         int       `json:"day"`
	Signpost    string    `json:"signpost"` // Blunt 3-4 word signpost (e.g. "NO LIVING THINGS")
	RuleType    string    `json:"rule_type"`
	TargetAdj   string    `json:"target_adj,omitempty"`
	TargetCargo string    `json:"target_cargo,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// RuleAudit represents the green/red light status of a single rule.
type RuleAudit struct {
	LawID    int64  `json:"law_id"`
	Signpost string `json:"signpost"`
	Passed   bool   `json:"passed"`
	Status   string `json:"status"` // "SLIPPED PAST" (green) or "BUSTED" (red)
}

// AuditVerdict is the overall result of a stamp attempt.
type AuditVerdict struct {
	Passed      bool        `json:"passed"`
	Status      string      `json:"status"` // "CLEARED" or "CONFISCATED"
	Audits      []RuleAudit `json:"audits"`
	Attempt     int         `json:"attempt"`
	MaxAttempts int         `json:"max_attempts"`
	Punchline   *Law        `json:"punchline,omitempty"`
}

// StarterLaws returns the 5 blunt 4-word signpost starter rules.
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

// EvaluateRule evaluates a single law against the cargo and player's disguise.
func EvaluateRule(law Law, cargo Cargo, d Disguise) bool {
	adj := strings.ToLower(strings.TrimSpace(d.Adjective))
	cnt := strings.ToLower(strings.TrimSpace(d.Container))

	switch law.RuleType {
	case "no_living":
		// Violated if cargo is living, UNLESS disguised as non-living (Inflatable, Metaphorical, Robotic)
		if cargo.HasTrait("living") {
			if adj == "inflatable" || adj == "metaphorical" {
				return true // Slipped past!
			}
			return false // Busted!
		}
		return true

	case "no_glass":
		// Violated if container is glass
		if strings.Contains(cnt, "glass") {
			return false // Busted!
		}
		return true

	case "pocket_only":
		// Must fit in a pocket.
		// Containers that naturally fit: envelope, pocket, velvet pouch, thermos
		// OR adjectives that make things fit: folded, microscopic, inflatable (deflated!)
		isPocketContainer := cnt == "pocket" || cnt == "envelope" || cnt == "velvet pouch"
		isCompressingAdj := adj == "folded" || adj == "microscopic" || adj == "inflatable"

		if isPocketContainer || isCompressingAdj {
			return true // Slipped past!
		}
		return false // Busted! (e.g. Lead Crate, Cardboard Box without compression)

	case "no_wet":
		// Violated if cargo is wet, UNLESS dehydrated, frozen, metaphorical, or in sealed thermos
		if cargo.HasTrait("wet") {
			if adj == "dehydrated" || adj == "frozen" || adj == "metaphorical" || cnt == "thermos" {
				return true // Slipped past!
			}
			return false // Busted!
		}
		return true

	case "no_metal":
		// Violated if container contains metal
		if cnt == "lead crate" || cnt == "thermos" {
			return false // Busted!
		}
		return true

	case "custom_punchline":
		// Specific ban on winning combo, e.g. "NO INFLATABLE GIRAFFES"
		targetAdj := strings.ToLower(strings.TrimSpace(law.TargetAdj))
		if adj == targetAdj && cargoMatches(cargo.Name, law.TargetCargo) {
			return false // Busted by yesterday's precedent!
		}
		return true

	default:
		return true
	}
}

func cargoMatches(name, target string) bool {
	n := strings.ToLower(strings.TrimSpace(name))
	t := strings.ToLower(strings.TrimSpace(target))
	tSingular := strings.TrimSuffix(t, "s")
	nSingular := strings.TrimSuffix(n, "s")

	return strings.Contains(n, t) || strings.Contains(t, n) ||
		strings.Contains(n, tSingular) || strings.Contains(nSingular, tSingular)
}

// Evaluate evaluates all active laws against cargo and disguise.
func Evaluate(laws []Law, cargo Cargo, d Disguise, attempt int) AuditVerdict {
	allPassed := true
	audits := make([]RuleAudit, 0, len(laws))

	for _, law := range laws {
		passed := EvaluateRule(law, cargo, d)
		status := "SLIPPED PAST"
		if !passed {
			allPassed = false
			status = "BUSTED"
		}
		audits = append(audits, RuleAudit{
			LawID:    law.ID,
			Signpost: law.Signpost,
			Passed:   passed,
			Status:   status,
		})
	}

	verdictStatus := "CONFISCATED"
	if allPassed {
		verdictStatus = "CLEARED"
	}

	return AuditVerdict{
		Passed:      allPassed,
		Status:      verdictStatus,
		Audits:      audits,
		Attempt:     attempt,
		MaxAttempts: 5,
	}
}

// CompilePunchlineLaw creates a 3-4 word punchline ban from the winning loophole.
// e.g. "NO INFLATABLE GIRAFFES", "NO DEHYDRATED WEATHER", "NO METAPHORICAL REGRETS"
func CompilePunchlineLaw(nextDay int, cargo Cargo, d Disguise, lawCount int) Law {
	noun := simplifyCargoName(cargo.Name)
	adj := strings.ToUpper(strings.TrimSpace(d.Adjective))

	signpost := fmt.Sprintf("NO %s %s", adj, noun)

	return Law{
		ID:          int64(lawCount + 1),
		Day:         nextDay,
		Signpost:    signpost,
		RuleType:    "custom_punchline",
		TargetAdj:   d.Adjective,
		TargetCargo: noun,
		CreatedAt:   time.Now(),
	}
}

func simplifyCargoName(name string) string {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "giraffe") {
		return "GIRAFFES"
	}
	if strings.Contains(lower, "weather") {
		return "WEATHER"
	}
	if strings.Contains(lower, "regret") {
		return "REGRETS"
	}
	if strings.Contains(lower, "cat") {
		return "CATS"
	}
	if strings.Contains(lower, "debt") {
		return "DEBT"
	}

	// Fallback: take last word and pluralize/uppercase
	words := strings.Fields(strings.ToUpper(name))
	if len(words) > 0 {
		last := words[len(words)-1]
		if !strings.HasSuffix(last, "S") {
			last += "S"
		}
		return last
	}
	return "ITEMS"
}
