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
	Adjective string `json:"adjective"` // e.g. "Inflatable", "Dehydrated", "Metaphorical"
	Container string `json:"container"` // e.g. "Envelope", "Pocket", "Matchbox"
}

// Law represents an active checkpoint rule.
type Law struct {
	ID          int64     `json:"id"`
	Day         int       `json:"day"`
	Signpost    string    `json:"signpost"` // Blunt signpost
	RuleType    string    `json:"rule_type"`
	TargetAdj   string    `json:"target_adj,omitempty"`
	TargetCargo string    `json:"target_cargo,omitempty"` // used for target container
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
	Passed       bool        `json:"passed"`
	Status       string      `json:"status"` // "CLEARED" or "CONFISCATED"
	Audits       []RuleAudit `json:"audits"`
	Attempt      int         `json:"attempt"`
	MaxAttempts  int         `json:"max_attempts"` // 3 attempts per day
	EmojiRow     string      `json:"emoji_row"`    // e.g. "🟩🟩🟩🟩🟩"
	CompiledLaws []Law       `json:"compiled_laws,omitempty"` // 2 new laws generated on win
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

// MetaphoricalSynonyms are "get out of jail free" words.
// If any is used, ALL are permanently banned for the rest of the run.
var MetaphoricalSynonyms = []string{
	"metaphorical", "metaphor", "imaginary", "allegorical", "symbolic",
	"fictional", "abstract", "spiritual", "figurative", "illusory",
	"conceptual", "pretend", "hallucinated", "mythical", "phantom",
	"make-believe", "philosophical", "metaphysical",
}

// IsMetaphorical checks if an adjective is metaphorical or one of its synonyms.
// Specifically: "theoretical" explicitly DOES NOT work.
func IsMetaphorical(adj string) bool {
	t := strings.ToLower(strings.TrimSpace(adj))
	if strings.Contains(t, "theoretical") || strings.Contains(t, "theory") {
		return false
	}
	for _, syn := range MetaphoricalSynonyms {
		if strings.Contains(t, syn) || strings.Contains(syn, t) {
			return true
		}
	}
	return false
}

// Dictionaries of semantic keywords for fair evaluation
var (
	NonLivingAdjectives = []string{
		"inflatable", "robotic", "mechanical", "taxidermied",
		"plastic", "wooden", "stuffed", "plush", "fake", "artificial",
		"holographic", "digital", "toy", "ceramic", "replica", "synthetic",
		"balloon", "origami", "painted", "drawn", "ghostly", "spectral",
		"decoy", "paper", "cardboard", "porcelain", "plaster", "fossilized",
		"freeze-dried", "statue", "dummy", "mannequin", "puppet",
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
		"ziploc bag", "sleeve", "sheath", "can", "tin can", "pocketbook",
		"plastic bag", "dry bag", "silicone pouch", "wax paper bag",
	}

	CompressingAdjectives = []string{
		"folded", "microscopic", "inflatable", "miniature", "tiny", "shrunk",
		"compressed", "flattened", "pocket-sized", "nano", "subatomic",
		"compact", "deflated", "origami", "dehydrated", "crushed", "condensed",
		"paper", "freeze-dried",
	}

	DryAdjectives = []string{
		"dehydrated", "dried", "freeze-dried", "frozen", "powdered", "powder",
		"evaporated", "baked", "burnt", "toasted", "cured", "smoked",
		"crystallized", "solidified", "digital", "taxidermied", "wooden",
		"plastic", "paper", "desiccated", "inflatable", "ceramic", "plush",
		"stuffed", "fossilized", "dry",
	}

	WaterproofContainers = []string{
		"ziploc", "ziploc bag", "plastic bag", "dry bag", "silicone pouch",
		"wax paper bag", "sealed pouch", "sealed packet", "capsule", "thermos",
		"vacuum flask", "vacuum bag",
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

	// "theoretical" explicitly does not work as a get-out-of-jail card
	isTheoretical := strings.Contains(adj, "theoretical") || strings.Contains(adj, "theory")
	isMeta := IsMetaphorical(adj) && !isTheoretical

	switch law.RuleType {
	case "no_living":
		if cargo.HasTrait("living") {
			if isTheoretical {
				return false // theoretical doesn't work!
			}
			if isMeta || matchAnyKeyword(adj, NonLivingAdjectives) {
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
		if isTheoretical {
			// theoretical doesn't compress
			return matchAnyKeyword(cnt, PocketContainers)
		}
		if isMeta || matchAnyKeyword(cnt, PocketContainers) || matchAnyKeyword(adj, CompressingAdjectives) {
			return true // Slipped past!
		}
		return false // Busted!

	case "no_wet":
		if cargo.HasTrait("wet") {
			if isTheoretical {
				return false // theoretical doesn't dry!
			}
			if isMeta || matchAnyKeyword(adj, DryAdjectives) || matchAnyKeyword(cnt, WaterproofContainers) {
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

	case "ban_all_metaphors":
		// When player burned their metaphorical card, ALL metaphorical/imaginary synonyms are banned!
		if isMeta || strings.Contains(adj, "metaphor") || strings.Contains(adj, "imaginar") {
			return false // Busted! Metaphorical joker already spent!
		}
		return true

	case "banned_adjective":
		target := strings.ToLower(strings.TrimSpace(law.TargetAdj))
		if target != "" && (strings.Contains(adj, target) || strings.Contains(target, adj)) {
			return false // Busted! Used previous day's loophole adjective
		}
		return true

	case "banned_container":
		target := strings.ToLower(strings.TrimSpace(law.TargetCargo))
		if target != "" && (strings.Contains(cnt, target) || strings.Contains(target, cnt)) {
			return false // Busted! Used previous day's loophole container
		}
		return true

	default:
		return true
	}
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

// CompileLoopholeRules creates TWO brand-new rules from the winning disguise:
// If metaphorical or its synonym was used, bans ALL metaphors for the rest of the run!
func CompileLoopholeRules(nextDay int, d Disguise, existingCount int) []Law {
	adj := strings.TrimSpace(d.Adjective)
	cnt := strings.TrimSpace(d.Container)

	var law1 Law
	if IsMetaphorical(adj) {
		// Used the get-out-of-jail free card -> Ban ALL metaphors & imaginary concepts!
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  "NO METAPHORS OR IMAGINARY CONCEPTS",
			RuleType:  "ban_all_metaphors",
			TargetAdj: "metaphorical",
			CreatedAt: time.Now(),
		}
	} else {
		// Specific adjective ban
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  fmt.Sprintf("NOTHING %s", strings.ToUpper(adj)),
			RuleType:  "banned_adjective",
			TargetAdj: strings.ToLower(adj),
			CreatedAt: time.Now(),
		}
	}

	law2 := Law{
		ID:          int64(existingCount + 2),
		Day:         nextDay,
		Signpost:    formatContainerSignpost(cnt),
		RuleType:    "banned_container",
		TargetCargo: strings.ToLower(cnt),
		CreatedAt:   time.Now(),
	}

	return []Law{law1, law2}
}

// FormatDisguiseSentence builds a grammatically correct sentence describing the smuggling attempt.
// e.g. "Metaphorical Live Giraffe in an Envelope", "Inflatable Live Giraffe in a Plastic Bag"
func FormatDisguiseSentence(adj, cargoName, container string) string {
	cleanCargo := strings.TrimSpace(cargoName)
	upper := strings.ToUpper(cleanCargo)
	if strings.HasPrefix(upper, "A ") {
		cleanCargo = strings.TrimSpace(cleanCargo[2:])
	} else if strings.HasPrefix(upper, "AN ") {
		cleanCargo = strings.TrimSpace(cleanCargo[3:])
	} else if strings.HasPrefix(upper, "THE ") {
		cleanCargo = strings.TrimSpace(cleanCargo[4:])
	}

	cleanAdj := strings.TrimSpace(adj)
	var subject string
	if cleanAdj != "" && !strings.EqualFold(cleanAdj, "Unmodified") {
		subject = fmt.Sprintf("%s %s", cleanAdj, cleanCargo)
	} else {
		subject = cleanCargo
	}

	cleanCnt := strings.TrimSpace(container)
	var containerPhrase string
	if cleanCnt != "" {
		cntLower := strings.ToLower(cleanCnt)
		if strings.HasPrefix(cntLower, "in ") || strings.HasPrefix(cntLower, "with ") {
			containerPhrase = cleanCnt
		} else if strings.EqualFold(cleanCnt, "bare hands") {
			containerPhrase = "with bare hands"
		} else {
			article := "a"
			firstChar := strings.ToLower(string(cleanCnt[0]))
			if strings.Contains("aeiou", firstChar) {
				article = "an"
			}
			containerPhrase = fmt.Sprintf("in %s %s", article, cleanCnt)
		}
	}

	if containerPhrase != "" {
		return fmt.Sprintf("%s %s", subject, containerPhrase)
	}
	return subject
}

func formatContainerSignpost(cnt string) string {
	upper := strings.ToUpper(strings.TrimSpace(cnt))
	if len(upper) == 0 {
		return "NO CONTAINERS"
	}

	firstLetter := string(upper[0])
	article := "A"
	if strings.Contains("AEIOU", firstLetter) {
		article = "AN"
	}
	return fmt.Sprintf("NOTHING IN %s %s", article, upper)
}
