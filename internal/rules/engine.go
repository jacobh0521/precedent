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

// Disguise is the player's 2-part disguise (Adjective + Container).
type Disguise struct {
	Adjective string `json:"adjective"` // e.g. "Inflatable", "Dehydrated", "Taxidermied"
	Container string `json:"container"` // e.g. "Envelope", "Pocket", "Matchbox"
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
	MaxAttempts int         `json:"max_attempts"` // 3 attempts per day
	EmojiRow    string      `json:"emoji_row"`    // e.g. "🟩🟩🟩🟩🟩" or "🟥🟩🟩🟩🟩"
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

// Broad dictionaries of semantic keywords for open text user entry
var (
	NonLivingAdjectives = []string{
		"inflatable", "metaphorical", "robotic", "mechanical", "taxidermied",
		"plastic", "wooden", "stuffed", "plush", "fake", "artificial",
		"holographic", "digital", "toy", "ceramic", "replica", "synthetic",
		"imaginary", "fictional", "balloon", "origami", "painted", "drawn",
		"ghostly", "spectral", "abstract", "symbolic", "spiritual", "allegorical",
		"virtual", "decoy", "paper", "cardboard", "porcelain", "plaster",
	}

	GlassKeywords = []string{
		"glass", "jar", "bottle", "vial", "flagon", "test tube", "flask",
		"beaker", "bulb", "aquarium", "terrarium", "decanter", "goblet",
		"chalice", "mason jar", "petri dish", "ampoule",
	}

	PocketContainers = []string{
		"pocket", "envelope", "velvet pouch", "pouch", "wallet", "sock",
		"thimble", "matchbox", "glove", "baggie", "capsule", "coin purse",
		"locket", "pillbox", "packet", "folder", "thermos", "ziploc",
		"sleeve", "sheath", "can", "tin can", "pocketbook",
	}

	CompressingAdjectives = []string{
		"folded", "microscopic", "inflatable", "miniature", "tiny", "shrunk",
		"compressed", "flattened", "pocket-sized", "nano", "subatomic",
		"compact", "deflated", "metaphorical", "abstract", "digital", "virtual",
		"origami", "dehydrated", "crushed", "condensed",
	}

	DryAdjectives = []string{
		"dehydrated", "dried", "freeze-dried", "frozen", "powdered", "powder",
		"metaphorical", "evaporated", "baked", "burnt", "toasted", "cured",
		"smoked", "crystallized", "solidified", "digital", "abstract",
		"taxidermied", "wooden", "plastic", "paper", "desiccated",
	}

	MetalKeywords = []string{
		"lead crate", "thermos", "metal", "lead", "steel", "iron", "tin",
		"aluminum", "safe", "vault", "can", "tin can", "cage", "toolbox",
		"lockbox", "ammo can", "brass box", "bronze urn", "copper pot",
		"foil", "titanium", "silver", "gold", "cast iron",
	}
)

func matchAnyKeyword(text string, keywords []string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	for _, kw := range keywords {
		if strings.Contains(t, kw) || strings.Contains(kw, t) {
			return true
		}
	}
	return false
}

// EvaluateRule evaluates a single law against the cargo and player's disguise.
func EvaluateRule(law Law, cargo Cargo, d Disguise) bool {
	adj := strings.ToLower(strings.TrimSpace(d.Adjective))
	cnt := strings.ToLower(strings.TrimSpace(d.Container))

	switch law.RuleType {
	case "no_living":
		if cargo.HasTrait("living") {
			// Slipped past if declared as non-living (e.g. Inflatable, Taxidermied, Robotic, Metaphorical)
			if matchAnyKeyword(adj, NonLivingAdjectives) {
				return true
			}
			return false // Busted!
		}
		return true

	case "no_glass":
		if matchAnyKeyword(cnt, GlassKeywords) {
			return false // Busted!
		}
		return true

	case "pocket_only":
		// Must fit in a pocket: either container fits in pocket OR adjective compresses it
		if matchAnyKeyword(cnt, PocketContainers) || matchAnyKeyword(adj, CompressingAdjectives) {
			return true // Slipped past!
		}
		return false // Busted!

	case "no_wet":
		if cargo.HasTrait("wet") {
			// Slipped past if dehydrated, frozen, dry adjective, or sealed container (like thermos)
			if matchAnyKeyword(adj, DryAdjectives) || strings.Contains(cnt, "thermos") || strings.Contains(cnt, "vacuum") {
				return true
			}
			return false // Busted!
		}
		return true

	case "no_metal":
		if matchAnyKeyword(cnt, MetalKeywords) {
			return false // Busted!
		}
		return true

	case "custom_punchline":
		targetAdj := strings.ToLower(strings.TrimSpace(law.TargetAdj))
		if strings.Contains(adj, targetAdj) && cargoMatches(cargo.Name, law.TargetCargo) {
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

// Evaluate evaluates all active laws against cargo and disguise (Max 3 attempts).
func Evaluate(laws []Law, cargo Cargo, d Disguise, attempt int) AuditVerdict {
	allPassed := true
	audits := make([]RuleAudit, 0, len(laws))
	var emojiBuilder strings.Builder

	for _, law := range laws {
		passed := EvaluateRule(law, cargo, d)
		status := "SLIPPED PAST"
		if !passed {
			allPassed = false
			status = "BUSTED"
			emojiBuilder.WriteString("🟥")
		} else {
			emojiBuilder.WriteString("🟩")
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
		MaxAttempts: 3, // 3 attempts per user per day
		EmojiRow:    emojiBuilder.String(),
	}
}

// CompilePunchlineLaw creates a 3-4 word punchline ban from the winning loophole.
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
