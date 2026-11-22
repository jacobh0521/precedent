package rules

import (
	"fmt"
	"strings"
	"time"
)

// MaxDailyAttempts is the single authoritative constant for user attempts per day.
const MaxDailyAttempts = 3

// Cargo represents the contraband item.
type Cargo struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Traits []string `json:"traits"` // e.g., "living", "wet", "oversized", "liquid", "gas", "sentient", "cursed", "flammable", "predator", "radiation"
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

// Dictionaries of rich semantic keywords for fair, flexible evaluation
var (
	NonLivingAdjectives = []string{
		"inflatable", "robotic", "mechanical", "taxidermied", "taxidermy",
		"plastic", "wooden", "stuffed", "plush", "fake", "artificial",
		"holographic", "digital", "toy", "ceramic", "replica", "synthetic",
		"balloon", "origami", "painted", "drawn", "ghostly", "spectral",
		"decoy", "paper", "cardboard", "porcelain", "plaster", "fossilized",
		"freeze-dried", "statue", "dummy", "mannequin", "puppet", "carved",
		"sculpted", "wax", "clay", "inert", "dormant", "petrified", "mummified",
		"cybernetic", "automated", "electronic", "deactivated", "sedated",
		"cremated", "ash", "ashes",
	}

	GlassKeywords = []string{
		"glass", "jar", "bottle", "vial", "flagon", "test tube", "flask",
		"beaker", "bulb", "aquarium", "terrarium", "decanter", "goblet",
		"chalice", "mason jar", "petri dish", "ampoule", "crystal", "crystal vase",
	}

	PocketContainers = []string{
		"pocket", "envelope", "velvet pouch", "pouch", "wallet", "sock",
		"thimble", "matchbox", "glove", "baggie", "capsule", "coin purse",
		"locket", "pillbox", "packet", "folder", "thermos", "ziploc",
		"ziploc bag", "sleeve", "sheath", "can", "tin can", "pocketbook",
		"plastic bag", "dry bag", "silicone pouch", "wax paper bag",
		"pencil case", "box of cards", "card case", "journal", "book",
		"hollow book", "spectacle case", "glasses case", "tin", "mint tin",
		"altoids tin", "film canister", "lunchbox", "bento box", "cigar box",
	}

	CompressingAdjectives = []string{
		"folded", "microscopic", "miniature", "tiny", "shrunk",
		"shrunken", "compressed", "flattened", "pocket-sized", "nano", "subatomic",
		"compact", "deflated", "origami", "dehydrated", "crushed", "condensed",
		"paper", "freeze-dried", "mini", "pocket", "micro", "flat", "concentrated",
		"minified", "reduced", "cremated", "ash", "powdered",
	}

	DryAdjectives = []string{
		"dehydrated", "dried", "freeze-dried", "frozen", "powdered", "powder",
		"evaporated", "baked", "burnt", "toasted", "cured", "smoked",
		"crystallized", "solidified", "digital", "taxidermied", "wooden",
		"plastic", "paper", "desiccated", "inflatable", "ceramic", "plush",
		"stuffed", "fossilized", "dry", "solid", "petrified", "congealed",
		"jellied", "calcified", "cremated", "ash", "ashes",
	}

	WaterproofContainers = []string{
		"ziploc", "ziploc bag", "plastic bag", "dry bag", "silicone pouch",
		"wax paper bag", "sealed pouch", "sealed packet", "capsule", "thermos",
		"vacuum flask", "vacuum bag", "tin can", "lockbox", "cooler",
	}

	MetalKeywords = []string{
		"lead crate", "thermos", "metal", "lead", "steel", "iron", "tin",
		"aluminum", "safe", "vault", "can", "tin can", "cage", "toolbox",
		"lockbox", "ammo can", "brass box", "bronze urn", "copper pot",
		"foil", "titanium", "silver", "gold", "cast iron", "metallic",
	}

	RigidContainers = []string{
		"box", "cardboard box", "lead crate", "crate", "tin", "metal box",
		"toolbox", "lockbox", "safe", "trunk", "chest", "vault", "cage",
		"barrel", "canister", "briefcase", "suitcase", "cooler", "thermos",
		"urn", "casket", "coffin", "mason jar", "jar", "tank", "drum",
		"pencil case", "box of cards", "card case", "hollow book", "lunchbox",
		"bento box", "shoebox", "cigar box", "mint tin", "altoids tin",
	}

	SoftContainers = []string{
		"pocket", "envelope", "velvet pouch", "pouch", "bag", "plastic bag",
		"ziploc", "ziploc bag", "backpack", "sack", "duffel", "tote bag", "tote",
		"sock", "glove", "sleeve", "baggie", "silicone pouch", "coin purse",
		"paper bag", "burlap sack", "mesh bag", "trash bag", "grocery bag",
	}

	TransparentContainers = []string{
		"glass", "jar", "clear bag", "ziploc", "ziploc bag", "plastic bag",
		"mason jar", "aquarium", "terrarium", "vial", "test tube",
		"acrylic box", "cellophane", "bottle", "petri dish", "beaker",
		"clear container", "transparent box", "display case",
	}

	FlammableMaterials = []string{
		"paper", "cardboard", "envelope", "folder", "wood", "wooden",
		"matchbox", "cloth", "cotton", "velvet", "burlap", "leaf",
	}

	OrganicContainers = []string{
		"paper", "cardboard", "envelope", "folder", "wood", "wooden",
		"leather", "cotton", "velvet", "burlap", "silk", "leaf", "coconut",
		"bamboo", "sock", "glove",
	}

	PacifyingAdjectives = []string{
		"plush", "toy", "taxidermied", "stuffed", "tame", "domesticated",
		"sleeping", "sedated", "hibernating", "docile", "inflatable", "wooden",
		"plastic", "origami", "replica", "inert", "pacified", "gentle",
		"trained", "fake", "ceramic", "statue", "calm", "dead", "cremated",
	}

	ShieldingAdjectives = []string{
		"lead", "shielded", "inert", "heavy", "dense", "insulated", "grounded",
	}

	// Solvability test pools
	SolvabilityAdjectivePool = []string{
		"Metaphorical", "Inflatable", "Folded", "Origami", "Dehydrated",
		"Frozen", "Microscopic", "Miniature", "Taxidermied", "Robotic",
		"Digital", "Cremated", "Plush", "Ceramic", "Lead", "Sedated",
		"Unmodified",
	}

	SolvabilityContainerPool = []string{
		"Pocket", "Envelope", "Matchbox", "Velvet Pouch", "Cardboard Box",
		"Lead Crate", "Glass Jar", "Ziploc Bag", "Plastic Bag", "Thermos",
		"Tin Can", "Briefcase", "Pencil Case", "Lunchbox", "Lockbox",
		"Wooden Chest", "Bare Hands",
	}
)

func matchAnyKeyword(text string, keywords []string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	words := strings.Fields(t)
	for _, kw := range keywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if t == kw {
			return true
		}
		if strings.Contains(kw, " ") || strings.Contains(kw, "-") {
			if strings.Contains(t, kw) {
				return true
			}
		} else {
			for _, w := range words {
				if w == kw || (len(kw) >= 4 && strings.HasPrefix(w, kw)) {
					return true
				}
			}
		}
	}
	return false
}

// ResolveCargoTraits calculates effective traits of the cargo after the adjective is applied,
// ensuring every adjective explicitly removes conflicting base traits.
func ResolveCargoTraits(cargo Cargo, adj string) map[string]bool {
	traits := make(map[string]bool)
	for _, tr := range cargo.Traits {
		traits[strings.ToLower(strings.TrimSpace(tr))] = true
	}

	adjLower := strings.ToLower(strings.TrimSpace(adj))
	isTheoretical := strings.Contains(adjLower, "theoretical") || strings.Contains(adjLower, "theory")
	if isTheoretical {
		// "theoretical" explicitly DOES NOT remove physical traits!
		return traits
	}

	if IsMetaphorical(adjLower) {
		// Metaphorical eliminates all physical constraints!
		for k := range traits {
			traits[k] = false
		}
		traits["metaphorical"] = true
		traits["abstract"] = true
		return traits
	}

	// 1. Cremated / Ash / Vaporized / Incinerated
	if matchAnyKeyword(adjLower, []string{"cremated", "ash", "ashes", "incinerated", "vaporized", "reduced to ash"}) {
		traits["living"] = false
		traits["wet"] = false
		traits["liquid"] = false
		traits["gas"] = false
		traits["oversized"] = false
		traits["large"] = false
		traits["tall"] = false
		traits["sentient"] = false
		traits["cursed"] = false
		traits["flammable"] = false
		traits["predator"] = false
		traits["aggressive"] = false
		traits["dry"] = true
		traits["ash"] = true
		traits["dead"] = true
	}

	// 2. Dehydrated / Dried / Freeze-Dried / Desiccated / Powdered
	if matchAnyKeyword(adjLower, DryAdjectives) {
		traits["wet"] = false
		traits["liquid"] = false
		traits["dry"] = true
		if matchAnyKeyword(adjLower, []string{"powdered", "powder", "freeze-dried"}) {
			traits["oversized"] = false
			traits["large"] = false
			traits["tall"] = false
		}
	}

	// 3. Frozen / Solidified / Crystallized
	if matchAnyKeyword(adjLower, []string{"frozen", "solidified", "crystallized", "solid", "ice"}) {
		traits["liquid"] = false
		traits["wet"] = false
		traits["gas"] = false
		traits["solid"] = true
		traits["dry"] = true
	}

	// 4. Inflatable / Deflated / Balloon
	if matchAnyKeyword(adjLower, []string{"inflatable", "deflated", "balloon"}) {
		traits["living"] = false
		traits["sentient"] = false
		traits["predator"] = false
		traits["aggressive"] = false
		traits["toy"] = true
		traits["fake"] = true
		traits["decoy"] = true
		// Deflated / foldable inflatable removes oversized constraint
		traits["oversized"] = false
		traits["large"] = false
		traits["tall"] = false
	}

	// 5. Compressing / Shrinking / Folding / Miniature
	if matchAnyKeyword(adjLower, CompressingAdjectives) {
		traits["oversized"] = false
		traits["large"] = false
		traits["tall"] = false
		traits["tiny"] = true
		traits["compressed"] = true
		if matchAnyKeyword(adjLower, []string{"origami", "paper"}) {
			traits["living"] = false
			traits["sentient"] = false
			traits["predator"] = false
		}
	}

	// 6. Toy / Decoy / Fake / Replica / Taxidermied / Robotic / Mechanical / Digital
	if matchAnyKeyword(adjLower, NonLivingAdjectives) {
		traits["living"] = false
		traits["sentient"] = false
		traits["cursed"] = false
		traits["predator"] = false
		traits["aggressive"] = false
		traits["fake"] = true
		traits["inert"] = true
	}

	// 7. Pacifying / Domesticated / Sedated / Sleeping
	if matchAnyKeyword(adjLower, PacifyingAdjectives) {
		traits["predator"] = false
		traits["aggressive"] = false
		traits["pacified"] = true
	}

	// 8. Lead / Shielded / Insulated
	if matchAnyKeyword(adjLower, ShieldingAdjectives) {
		traits["radiation"] = false
		traits["ionized"] = false
		traits["shielded"] = true
	}

	return traits
}

// EvaluateRule evaluates a single law against the cargo and player's disguise.
func EvaluateRule(law Law, cargo Cargo, d Disguise) bool {
	adj := strings.ToLower(strings.TrimSpace(d.Adjective))
	cnt := strings.ToLower(strings.TrimSpace(d.Container))

	// "theoretical" explicitly does not work as a get-out-of-jail card
	isTheoretical := strings.Contains(adj, "theoretical") || strings.Contains(adj, "theory")
	isMeta := IsMetaphorical(adj) && !isTheoretical

	effectiveTraits := ResolveCargoTraits(cargo, adj)

	switch law.RuleType {
	case "no_living":
		if effectiveTraits["living"] {
			return false
		}
		return true

	case "no_glass":
		if matchAnyKeyword(cnt, GlassKeywords) {
			return false
		}
		return true

	case "pocket_only":
		if isTheoretical {
			return false // theoretical doesn't shrink or fit
		}
		if isMeta {
			return true // metaphorical fits anywhere
		}
		// Container must be pocket-sized:
		if !matchAnyKeyword(cnt, PocketContainers) {
			return false
		}
		// Cargo must NOT still have oversized/large/tall trait:
		if effectiveTraits["oversized"] || effectiveTraits["large"] || effectiveTraits["tall"] {
			return false
		}
		return true

	case "no_wet":
		if effectiveTraits["wet"] || effectiveTraits["liquid"] {
			if matchAnyKeyword(cnt, WaterproofContainers) {
				return true
			}
			return false
		}
		return true

	case "no_metal":
		if matchAnyKeyword(cnt, MetalKeywords) {
			return false
		}
		return true

	// Regime 2 Rules
	case "no_liquid_gas":
		if effectiveTraits["liquid"] || effectiveTraits["gas"] || effectiveTraits["wet"] {
			if matchAnyKeyword(cnt, WaterproofContainers) {
				return true
			}
			return false
		}
		return true

	case "no_sentient_cursed":
		if effectiveTraits["sentient"] || effectiveTraits["cursed"] || effectiveTraits["living"] {
			return false
		}
		return true

	case "rigid_container":
		// Must be a rigid container; cannot be soft or flexible
		if matchAnyKeyword(cnt, RigidContainers) {
			return true
		}
		if matchAnyKeyword(cnt, SoftContainers) {
			return false
		}
		return true

	case "no_flammable":
		if effectiveTraits["flammable"] || matchAnyKeyword(cnt, FlammableMaterials) || matchAnyKeyword(adj, []string{"burning", "flaming", "fire", "combustible", "flammable"}) {
			if isMeta {
				return true
			}
			return false
		}
		return true

	// Regime 3 Rules
	case "no_organic_container":
		if matchAnyKeyword(cnt, OrganicContainers) {
			return false
		}
		return true

	case "transparent_container":
		if matchAnyKeyword(cnt, TransparentContainers) {
			return true
		}
		return false

	case "no_predators":
		if effectiveTraits["predator"] || effectiveTraits["aggressive"] {
			return false
		}
		return true

	case "no_radiation":
		if effectiveTraits["radiation"] || effectiveTraits["ionized"] || effectiveTraits["cosmic"] {
			if matchAnyKeyword(cnt, []string{"lead", "lead crate", "shielded"}) {
				return true
			}
			return false
		}
		return true

	case "must_use_container":
		if cnt == "" || matchAnyKeyword(cnt, []string{"bare hands", "hands", "none", "uncontained"}) {
			return false
		}
		return true

	// Precedent Law Categories & Loophole Bans
	case "ban_all_metaphors":
		if isMeta || strings.Contains(adj, "metaphor") || strings.Contains(adj, "imaginar") {
			return false
		}
		return true

	case "ban_shrinking_adj":
		if matchAnyKeyword(adj, CompressingAdjectives) {
			return false
		}
		return true

	case "ban_drying_adj":
		if matchAnyKeyword(adj, DryAdjectives) {
			return false
		}
		return true

	case "ban_decoy_adj":
		if matchAnyKeyword(adj, []string{"plush", "inflatable", "fake", "toy", "replica", "synthetic", "origami", "taxidermied", "taxidermy", "dummy", "puppet", "mannequin", "statue", "drawn", "painted"}) {
			return false
		}
		return true

	case "ban_synthetic_adj":
		if matchAnyKeyword(adj, []string{"robotic", "mechanical", "digital", "cybernetic", "automated", "electronic"}) {
			return false
		}
		return true

	case "ban_paper_containers":
		if matchAnyKeyword(cnt, []string{"envelope", "folder", "paper", "cardboard", "packet", "wrapper", "journal", "book", "origami", "sleeve"}) {
			return false
		}
		return true

	case "ban_soft_pockets":
		if matchAnyKeyword(cnt, []string{"pocket", "velvet pouch", "pouch", "sock", "glove", "baggie", "coin purse", "pencil case", "silicone pouch"}) {
			return false
		}
		return true

	case "ban_flexible_bags":
		if matchAnyKeyword(cnt, []string{"plastic bag", "backpack", "tote bag", "duffel", "burlap sack", "mesh bag", "grocery bag", "dry bag", "trash bag", "sack", "bag"}) {
			return false
		}
		return true

	case "ban_glass_containers":
		if matchAnyKeyword(cnt, GlassKeywords) {
			return false
		}
		return true

	case "ban_metal_containers":
		if matchAnyKeyword(cnt, MetalKeywords) {
			return false
		}
		return true

	case "banned_adjective":
		target := strings.ToLower(strings.TrimSpace(law.TargetAdj))
		if target != "" && (strings.Contains(adj, target) || strings.Contains(target, adj)) {
			return false
		}
		return true

	case "banned_container":
		target := strings.ToLower(strings.TrimSpace(law.TargetCargo))
		if target != "" && (strings.Contains(cnt, target) || strings.Contains(target, cnt)) {
			return false
		}
		return true

	default:
		return true
	}
}

// Evaluate evaluates all active laws against cargo and disguise (MaxDailyAttempts limit).
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
		MaxAttempts: MaxDailyAttempts, // Synced to constant
		EmojiRow:    emojiBuilder.String(),
	}
}

// CheckSolvability verifies whether at least 1 winning disguise combo exists for the given cargo under the laws.
func CheckSolvability(laws []Law, cargo Cargo) bool {
	for _, adj := range SolvabilityAdjectivePool {
		for _, cnt := range SolvabilityContainerPool {
			v := Evaluate(laws, cargo, Disguise{Adjective: adj, Container: cnt}, 1)
			if v.Passed {
				return true
			}
		}
	}
	return false
}

// CompileLoopholeRules creates TWO mutated rules from the winning disguise.
// Includes a solvability check against tomorrow's cargo to guarantee at least 1 winning combo remains!
func CompileLoopholeRules(nextDay int, d Disguise, existingLaws []Law, tomorrowCargo ...Cargo) []Law {
	existingCount := len(existingLaws)
	adj := strings.TrimSpace(d.Adjective)
	cnt := strings.TrimSpace(d.Container)
	adjLower := strings.ToLower(adj)
	cntLower := strings.ToLower(cnt)

	var law1 Law
	if IsMetaphorical(adj) {
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  "NO METAPHORS OR IMAGINARY CONCEPTS",
			RuleType:  "ban_all_metaphors",
			TargetAdj: "metaphorical",
			CreatedAt: time.Now(),
		}
	} else if matchAnyKeyword(adjLower, CompressingAdjectives) {
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  "NO SHRINKING OR SCALING ADJECTIVES",
			RuleType:  "ban_shrinking_adj",
			TargetAdj: adjLower,
			CreatedAt: time.Now(),
		}
	} else if matchAnyKeyword(adjLower, []string{"dehydrated", "freeze-dried", "frozen", "powdered", "evaporated", "dry", "dried", "cured", "desiccated"}) {
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  "NO DEHYDRATING OR DRYING ADJECTIVES",
			RuleType:  "ban_drying_adj",
			TargetAdj: adjLower,
			CreatedAt: time.Now(),
		}
	} else if matchAnyKeyword(adjLower, []string{"plush", "inflatable", "fake", "toy", "replica", "synthetic", "origami", "taxidermied", "taxidermy", "dummy", "puppet", "mannequin", "statue"}) {
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  "NO FAKE, TOY, OR DECOY ADJECTIVES",
			RuleType:  "ban_decoy_adj",
			TargetAdj: adjLower,
			CreatedAt: time.Now(),
		}
	} else if matchAnyKeyword(adjLower, []string{"robotic", "mechanical", "digital", "cybernetic", "automated", "electronic"}) {
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  "NO SYNTHETIC OR ROBOTIC ADJECTIVES",
			RuleType:  "ban_synthetic_adj",
			TargetAdj: adjLower,
			CreatedAt: time.Now(),
		}
	} else {
		law1 = Law{
			ID:        int64(existingCount + 1),
			Day:       nextDay,
			Signpost:  fmt.Sprintf("NOTHING %s", strings.ToUpper(adj)),
			RuleType:  "banned_adjective",
			TargetAdj: adjLower,
			CreatedAt: time.Now(),
		}
	}

	var law2 Law
	if strings.Contains(cntLower, "bare hands") || strings.Contains(cntLower, "hands") || cntLower == "none" || cntLower == "uncontained" {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "MUST DECLARE A REAL CONTAINER",
			RuleType:    "must_use_container",
			TargetCargo: "bare hands",
			CreatedAt:   time.Now(),
		}
	} else if matchAnyKeyword(cntLower, []string{"envelope", "folder", "paper bag", "packet", "wrapper", "journal", "book", "hollow book", "origami", "sleeve"}) {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "NO PAPER OR ENVELOPE CONTAINERS",
			RuleType:    "ban_paper_containers",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	} else if matchAnyKeyword(cntLower, []string{"pocket", "velvet pouch", "pouch", "sock", "glove", "coin purse", "baggie", "silicone pouch", "pencil case"}) {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "NO SOFT POCKETS OR POUCHES",
			RuleType:    "ban_soft_pockets",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	} else if matchAnyKeyword(cntLower, []string{"plastic bag", "backpack", "tote bag", "duffel", "burlap sack", "mesh bag", "trash bag", "grocery bag"}) {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "NO BAGS OR FLEXIBLE SACKS",
			RuleType:    "ban_flexible_bags",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	} else if matchAnyKeyword(cntLower, GlassKeywords) {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "NO GLASS OR CRYSTAL CONTAINERS",
			RuleType:    "ban_glass_containers",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	} else if matchAnyKeyword(cntLower, MetalKeywords) {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "NO METAL OR ARMORED CONTAINERS",
			RuleType:    "ban_metal_containers",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	} else if matchAnyKeyword(cntLower, RigidContainers) {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    "CONTAINERS MUST BE TRANSPARENT",
			RuleType:    "transparent_container",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	} else {
		law2 = Law{
			ID:          int64(existingCount + 2),
			Day:         nextDay,
			Signpost:    formatContainerSignpost(cnt),
			RuleType:    "banned_container",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
	}

	// Solvability check: If tomorrow's cargo is provided, ensure at least 1 winning combo remains!
	if len(tomorrowCargo) > 0 && tomorrowCargo[0].ID != "" {
		cargo := tomorrowCargo[0]
		combined := append(append([]Law{}, existingLaws...), law1, law2)

		// 1. Check if both proposed laws keep the game solvable
		if CheckSolvability(combined, cargo) {
			return []Law{law1, law2}
		}

		// 2. Fallback: Soften law1 to specific adjective ban
		law1Soft := Law{
			ID:        law1.ID,
			Day:       nextDay,
			Signpost:  fmt.Sprintf("NOTHING %s", strings.ToUpper(adj)),
			RuleType:  "banned_adjective",
			TargetAdj: adjLower,
			CreatedAt: time.Now(),
		}
		if CheckSolvability(append(append([]Law{}, existingLaws...), law1Soft, law2), cargo) {
			return []Law{law1Soft, law2}
		}

		// 3. Fallback: Also soften law2 to specific container ban
		law2Soft := Law{
			ID:          law2.ID,
			Day:         nextDay,
			Signpost:    formatContainerSignpost(cnt),
			RuleType:    "banned_container",
			TargetCargo: cntLower,
			CreatedAt:   time.Now(),
		}
		if CheckSolvability(append(append([]Law{}, existingLaws...), law1Soft, law2Soft), cargo) {
			return []Law{law1Soft, law2Soft}
		}

		// 4. Fallback: If still unsolvable, keep whichever single law preserves solvability
		if CheckSolvability(append(append([]Law{}, existingLaws...), law1Soft), cargo) {
			return []Law{law1Soft}
		}
		if CheckSolvability(append(append([]Law{}, existingLaws...), law2Soft), cargo) {
			return []Law{law2Soft}
		}
	}

	return []Law{law1, law2}
}

// toTitleWords converts a string to Title Case words.
func toTitleWords(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		lower := strings.ToLower(w)
		if len(lower) > 0 {
			words[i] = strings.ToUpper(lower[:1]) + lower[1:]
		}
	}
	return strings.Join(words, " ")
}

// FormatDisguiseSentence builds a grammatically correct sentence describing the smuggling attempt.
// Cleans up manifest display names so adjectives replace "LIVE" instead of printing "Cremated LIVE GIRAFFE".
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
		// When modified by an adjective, replace "LIVE " in cargo display
		// e.g. "Cremated Giraffe" instead of "Cremated LIVE GIRAFFE"
		cargoWithoutLive := cleanCargo
		if strings.HasPrefix(strings.ToUpper(cargoWithoutLive), "LIVE ") {
			cargoWithoutLive = strings.TrimSpace(cargoWithoutLive[5:])
		}
		subject = fmt.Sprintf("%s %s", cleanAdj, toTitleWords(cargoWithoutLive))
	} else {
		subject = toTitleWords(cleanCargo)
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
