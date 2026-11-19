package db

import (
	"fmt"
)

// CuratedIconicCargo contains the opening iconic daily items.
var CuratedIconicCargo = []Manifest{
	{
		ID:          "cargo-day-1",
		Day:         1,
		Name:        "A LIVE GIRAFFE",
		Traits:      []string{"living", "oversized"},
		Description: "A 5.5-meter long-necked herbivorous ungulate with an inquisitive stare.",
		FlavorText:  "Staring at customs inspectors with absolute nonchalance.",
	},
	{
		ID:          "cargo-day-2",
		Day:         2,
		Name:        "TOMORROW'S WEATHER",
		Traits:      []string{"wet", "oversized"},
		Description: "A pressurized atmospheric cold front with 90% rain probability.",
		FlavorText:  "Already dripping onto the conveyer belt.",
	},
	{
		ID:          "cargo-day-3",
		Day:         3,
		Name:        "THREE REGRETS",
		Traits:      []string{"wet", "oversized"},
		Description: "Heavy psychic residue soaked in cold tears from a misspent youth.",
		FlavorText:  "Deeply awkward to declare at the customs terminal.",
	},
	{
		ID:          "cargo-day-4",
		Day:         4,
		Name:        "A QUANTUM CAT",
		Traits:      []string{"living", "oversized"},
		Description: "Simultaneously purring and quantum-decohering in super-position.",
		FlavorText:  "Customs warning: Do not open carrier or waveform collapses.",
	},
	{
		ID:          "cargo-day-5",
		Day:         5,
		Name:        "A BOTTLED NEBULA",
		Traits:      []string{"wet", "oversized"},
		Description: "Stellar dust cloud condensing into ionized cosmic precipitation.",
		FlavorText:  "Shipper declared contents as: 'ambient pocket warmer'.",
	},
	{
		ID:          "cargo-day-6",
		Day:         6,
		Name:        "AN ANGRY GOOSE",
		Traits:      []string{"living", "oversized"},
		Description: "A vicious avian terror with zero regard for interstellar law.",
		FlavorText:  "Hissing aggressively at the radiation scanner.",
	},
	{
		ID:          "cargo-day-7",
		Day:         7,
		Name:        "A SLEEPING VOLCANO",
		Traits:      []string{"wet", "oversized"},
		Description: "Molten magma chamber snoring quietly beneath tectonic plates.",
		FlavorText:  "Slight seismic tremors detected in security line.",
	},
	{
		ID:          "cargo-day-8",
		Day:         8,
		Name:        "EXISTENTIAL DREAD",
		Traits:      []string{"wet", "oversized"},
		Description: "Vast metaphysical weight that makes inspectors contemplate retirement.",
		FlavorText:  "Declared value: Nothing matters anyway.",
	},
	{
		ID:          "cargo-day-9",
		Day:         9,
		Name:        "A RUNAWAY GLACIER",
		Traits:      []string{"wet", "oversized"},
		Description: "Four million cubic meters of ancient blue pack-ice on the move.",
		FlavorText:  "Currently traveling at 0.0002 miles per hour.",
	},
	{
		ID:          "cargo-day-10",
		Day:         10,
		Name:        "A GALLON OF TEARS",
		Traits:      []string{"wet", "oversized"},
		Description: "Harvested from heartbroken poets across seven star systems.",
		FlavorText:  "Smells faintly of lavender and broken promises.",
	},
	{
		ID:          "cargo-day-11",
		Day:         11,
		Name:        "THE CONCEPT OF MONDAY",
		Traits:      []string{"wet", "oversized"},
		Description: "A collective damp depression draining morale across the orbital ring.",
		FlavorText:  "Station clocks mysteriously lose 15 minutes in its presence.",
	},
	{
		ID:          "cargo-day-12",
		Day:         12,
		Name:        "A BLUE WHALE CALF",
		Traits:      []string{"living", "wet", "oversized"},
		Description: "Twenty tons of mammalian marine royalty looking for krill.",
		FlavorText:  "Echolocating loudly into customs intercom.",
	},
	{
		ID:          "cargo-day-13",
		Day:         13,
		Name:        "AN ACTIVE WATERFALL",
		Traits:      []string{"wet", "oversized"},
		Description: "Roaring cascade of pristine mountain river water falling endlessly.",
		FlavorText:  "Inspector boots are completely soaked.",
	},
	{
		ID:          "cargo-day-14",
		Day:         14,
		Name:        "AN ELEPHANT IN THE ROOM",
		Traits:      []string{"living", "oversized"},
		Description: "Enormous pachyderm that everyone is strictly pretending not to notice.",
		FlavorText:  "Customs officer: 'I see no contraband on this conveyor.'",
	},
}

// Procedural building blocks for Day 15 to Day 500+ (guaranteeing unique daily contraband through 2028)
var (
	Descriptors = []string{
		"A FERAL", "A LUMINESCENT", "A CARNIVOROUS", "A TAX-EXEMPT", "AN UNSTABLE",
		"A HAUNTED", "A SYNTHETIC", "AN HYPERACTIVE", "A FROZEN", "A COSMIC",
		"A SENTIENT", "A DIPLOMATIC", "A REBELLIOUS", "AN OBSCURE", "AN EXTINCT",
		"A DISORIENTED", "A TITANIC", "AN INTERSTELLAR", "A SACRED", "A BOILING",
		"AN ANARCHIST", "A SUBATOMIC", "A CRYPTIC", "AN UNREGISTERED", "A CELESTIAL",
	}

	Nouns = []struct {
		Name   string
		Traits []string
		Desc   string
		Flavor string
	}{
		{"OCTOPUS", []string{"living", "wet", "oversized"}, "Eight tentacled cephalopod picking customs locks.", "Currently wearing a fake passport."},
		{"GRIZZLY BEAR", []string{"living", "oversized"}, "A 900-pound apex predator in transit.", "Smells like salmon and territorial hostility."},
		{"GEYSER", []string{"wet", "oversized"}, "Pressurized boiling geothermal water fountain.", "Erupts every 4 minutes on the dot."},
		{"SWARM OF BEES", []string{"living", "oversized"}, "50,000 angry pollinators inside a burlap sac.", "Customs officers standing very, very still."},
		{"PIRATE GALLEON", []string{"wet", "oversized"}, "Full-sized 17th-century wooden warship with sails raised.", "Cannons currently loaded with citrus fruits."},
		{"TYRANNOSAURUS REX", []string{"living", "oversized"}, "Cloned Cretaceous carnivore seeking a mid-day snack.", "Has tiny arms but very large boarding grievances."},
		{"RAINCLOUD", []string{"wet", "oversized"}, "Localized precipitation cell following individual travelers.", "Drizzling melancholy on the scanner."},
		{"MOOSE", []string{"living", "oversized"}, "Six-foot antlered forest monarch with no ticket.", "Chewing on customs declaration paperwork."},
		{"SOUP TSUNAMI", []string{"wet", "oversized"}, "Five thousand gallons of bubbling chicken noodle broth.", "Requires emergency customs crackers."},
		{"HIPPOPOTAMUS", []string{"living", "wet", "oversized"}, "Two-ton aquatic territory defender with bad attitude.", "Submerged in the terminal drinking fountain."},
		{"TORNADO", []string{"oversized"}, "F5 vortex of spinning debris and cosmic dust.", "Currently rearranging the duty-free shop."},
		{"GIANT PANDA", []string{"living", "oversized"}, "Bamboo-obsessed bear refusing to leave terminal floor.", "Declared by consignee as 'heavy throw pillow'."},
		{"BREWERY VAT", []string{"wet", "oversized"}, "Overflowing wooden reservoir of fermented dwarven ale.", "Customs agents suddenly volunteering for inspection."},
		{"OSTRICH", []string{"living", "oversized"}, "Seven-foot flightless bird capable of 45 mph sprints.", "Stuck its head into the x-ray machine."},
		{"ICEBERG", []string{"wet", "oversized"}, "Titanic block of glacial ice with submerged hazards.", "90% of its volume hidden beneath customs counter."},
		{"GORILLA", []string{"living", "oversized"}, "Silverback primate inspecting customs badges with curiosity.", "Taught sign language: 'Let me through.'"},
		{"KELP FOREST", []string{"wet", "oversized"}, "Dense tangled marine grove teeming with bioluminescent crabs.", "Leaking ocean brine across terminal C."},
		{"FLAMINGO FLOCK", []string{"living", "oversized"}, "Sixty bright pink wading birds standing on one leg.", "Refusing to provide customs fingerprints."},
		{"SWAMP", []string{"wet", "oversized"}, "One acre of humid boggy peat and croaking toads.", "Mud stains on station carpet are permanent."},
		{"WALRUS", []string{"living", "wet", "oversized"}, "Tusked blubber colossus taking a nap on terminal scale.", "Weight capacity severely exceeded."},
	}
)

// GetManifestForDay returns the single deterministic contraband parcel for any calendar day (1 to 500+).
func GetManifestForDay(day int) Manifest {
	if day < 1 {
		day = 1
	}

	// First 14 days use hand-crafted iconic items
	if day <= len(CuratedIconicCargo) {
		m := CuratedIconicCargo[day-1]
		m.Day = day
		return m
	}

	// Days 15 to 500+ deterministically combine 25 Descriptors x 20 Nouns = 500 unique combinations
	offset := day - 15
	descIdx := offset % len(Descriptors)
	nounIdx := (offset / len(Descriptors)) % len(Nouns)

	d := Descriptors[descIdx]
	n := Nouns[nounIdx]

	name := fmt.Sprintf("%s %s", d, n.Name)
	return Manifest{
		ID:          fmt.Sprintf("cargo-day-%d", day),
		Day:         day,
		Name:        name,
		Traits:      n.Traits,
		Description: n.Desc,
		FlavorText:  fmt.Sprintf("Day %d manifest: %s", day, n.Flavor),
	}
}
