package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"precedent/internal/db"
	"precedent/internal/rules"
)

// GameHandler coordinates HTTP requests, game state, and templating.
type GameHandler struct {
	store     *db.Store
	templates *template.Template
}

// Word suggestions for datalists and quick-pick chips
var (
	AdjectiveSuggestions = []string{
		"Inflatable", "Dehydrated", "Folded", "Metaphorical", "Frozen",
		"Microscopic", "Taxidermied", "Robotic", "Ceremonial", "Rental",
		"Origami", "Plastic", "Wooden", "Ghostly", "Digital",
		"Holographic", "Miniature", "Stuffed", "Decoy", "Paper",
	}

	ContainerSuggestions = []string{
		"Envelope", "Pocket", "Matchbox", "Thermos", "Velvet Pouch",
		"Cardboard Box", "Glass Jar", "Lead Crate", "Mesh Bag", "Sock",
		"Coin Purse", "Thimble", "Ziploc Bag", "Tin Can", "Briefcase",
		"Backpack", "Shoebox", "Pillowcase", "Lunchbox", "Ammo Can",
	}
)

// DeskViewData holds data for the main play view.
type DeskViewData struct {
	CurrentDay           int
	AttemptCount         int
	MaxAttempts          int
	Manifests            []db.Manifest
	SelectedManifest     *db.Manifest
	ActiveLaws           []rules.Law
	AdjectiveSuggestions []string
	ContainerSuggestions []string
	WeeklyDebrief        db.WeeklyDebrief
}

// VerdictPartialData holds data passed to verdict.html.
type VerdictPartialData struct {
	Verdict        rules.AuditVerdict
	Manifest       db.Manifest
	Disguise       rules.Disguise
	ActiveLaws     []rules.Law
	CurrentDay     int
	OutOfTries     bool
	DailyShareText string
}

// NewGameHandler initializes handler with template parsing.
func NewGameHandler(store *db.Store, templateDir string) (*GameHandler, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"join": strings.Join,
		"sub": func(a, b int) int {
			diff := a - b
			if diff < 0 {
				return 0
			}
			return diff
		},
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := 0; i < n; i++ {
				s[i] = i + 1
			}
			return s
		},
	}).ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &GameHandler{
		store:     store,
		templates: tmpl,
	}, nil
}

// RegisterRoutes sets up Chi router endpoints.
func (h *GameHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.HandleDeskView)
	r.Post("/declare", h.HandleSubmitDisguise)
	r.Get("/select-manifest", h.HandleSelectManifest)
	r.Get("/weekly-debrief", h.HandleWeeklyDebrief)
	r.Post("/retry", h.HandleRetryAttempts)
	r.Post("/dev/next-day", h.HandleSimulateNextDay)
	r.Post("/dev/reset", h.HandleResetDay1)
}

// HandleDeskView renders the streamlined customs checkpoint.
func (h *GameHandler) HandleDeskView(w http.ResponseWriter, r *http.Request) {
	currentDay, _ := h.store.GetCurrentDay()
	attemptCount, _ := h.store.GetAttemptCount()
	allManifests, _ := h.store.GetAllManifests()

	dayManifests := make([]db.Manifest, 0)
	for _, m := range allManifests {
		if m.Day <= currentDay {
			dayManifests = append(dayManifests, m)
		}
	}
	if len(dayManifests) == 0 {
		dayManifests = allManifests
	}

	selectedID := r.URL.Query().Get("manifest_id")
	if selectedID == "" {
		selectedID, _ = h.store.GetSelectedManifestID()
	}
	if selectedID == "" && len(dayManifests) > 0 {
		selectedID = dayManifests[0].ID
	}

	selectedManifest, _ := h.store.GetManifestByID(selectedID)
	if selectedManifest == nil && len(dayManifests) > 0 {
		selectedManifest = &dayManifests[0]
	}

	activeLaws, _ := h.store.GetActiveLaws()
	weeklyDebrief, _ := h.store.GetWeeklyDebrief()

	data := DeskViewData{
		CurrentDay:           currentDay,
		AttemptCount:         attemptCount,
		MaxAttempts:          3, // 3 attempts per day
		Manifests:            dayManifests,
		SelectedManifest:     selectedManifest,
		ActiveLaws:           activeLaws,
		AdjectiveSuggestions: AdjectiveSuggestions,
		ContainerSuggestions: ContainerSuggestions,
		WeeklyDebrief:        weeklyDebrief,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

// HandleSelectManifest switches the active parcel inspected on desk.
func (h *GameHandler) HandleSelectManifest(w http.ResponseWriter, r *http.Request) {
	manifestID := r.URL.Query().Get("manifest_id")
	if manifestID != "" {
		_ = h.store.SetSelectedManifestID(manifestID)
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/?manifest_id="+manifestID)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/?manifest_id="+manifestID, http.StatusSeeOther)
}

// HandleSubmitDisguise evaluates user-entered disguise against active laws.
func (h *GameHandler) HandleSubmitDisguise(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Malformed form", http.StatusBadRequest)
		return
	}

	manifestID := r.FormValue("manifest_id")
	adjective := strings.TrimSpace(r.FormValue("adjective"))
	container := strings.TrimSpace(r.FormValue("container"))

	if adjective == "" {
		adjective = "Unmodified"
	}
	if container == "" {
		container = "Bare Hands"
	}

	manifest, err := h.store.GetManifestByID(manifestID)
	if err != nil {
		http.Error(w, "Manifest parcel not found", http.StatusNotFound)
		return
	}

	currentDay, _ := h.store.GetCurrentDay()
	attemptCount, _ := h.store.GetAttemptCount()
	activeLaws, _ := h.store.GetActiveLaws()

	disguise := rules.Disguise{
		Adjective: adjective,
		Container: container,
	}

	cargo := manifest.ToCargo()
	verdict := rules.Evaluate(activeLaws, cargo, disguise, attemptCount)

	// If Cleared: Compile punchline law
	if verdict.Passed {
		newLaw := rules.CompilePunchlineLaw(currentDay+1, cargo, disguise, len(activeLaws))
		_, err := h.store.SaveLaw(newLaw)
		if err == nil {
			verdict.Punchline = &newLaw
		}
	} else {
		// Increment attempts (up to 3)
		if attemptCount < 3 {
			_, _ = h.store.IncrementAttempt()
		}
	}

	_ = h.store.SaveDeclaration(manifestID, currentDay, disguise, attemptCount, verdict)
	updatedLaws, _ := h.store.GetActiveLaws()

	outOfTries := (!verdict.Passed && attemptCount >= 3)

	// Fetch all attempts today for this manifest to construct the Wordle-like share text
	records, _ := h.store.GetDeclarationsForDay(manifestID, currentDay)
	shareText := formatDailyShare(currentDay, manifest.Name, verdict.Passed, attemptCount, records, verdict.Punchline)

	partialData := VerdictPartialData{
		Verdict:        verdict,
		Manifest:       *manifest,
		Disguise:       disguise,
		ActiveLaws:     updatedLaws,
		CurrentDay:     currentDay,
		OutOfTries:     outOfTries,
		DailyShareText: shareText,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		_ = h.templates.ExecuteTemplate(w, "verdict.html", partialData)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func formatDailyShare(day int, cargoName string, cleared bool, attempts int, records []db.DeclarationRecord, punchline *rules.Law) string {
	score := fmt.Sprintf("%d/3", attempts)
	if !cleared && attempts >= 3 {
		score = "X/3"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PRECEDENT Day %d %s\n", day, score))
	sb.WriteString(fmt.Sprintf("Contraband: %s\n\n", cargoName))

	for _, rec := range records {
		sb.WriteString(fmt.Sprintf("%s (%s + %s)\n", rec.EmojiRow, rec.Adjective, rec.Container))
	}

	if punchline != nil {
		sb.WriteString(fmt.Sprintf("\n⚖️ New Precedent: %s\n", punchline.Signpost))
	}
	sb.WriteString("\nSmuggle your cargo: https://github.com/jacobh0521/precedent")
	return sb.String()
}

// HandleWeeklyDebrief displays the week's end customs report modal.
func (h *GameHandler) HandleWeeklyDebrief(w http.ResponseWriter, r *http.Request) {
	debrief, _ := h.store.GetWeeklyDebrief()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.templates.ExecuteTemplate(w, "weekly_debrief.html", debrief)
}

// HandleRetryAttempts resets attempt count to 1 for the current parcel.
func (h *GameHandler) HandleRetryAttempts(w http.ResponseWriter, r *http.Request) {
	_ = h.store.ResetAttempt()
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleSimulateNextDay simulates the next day to play with new punchline laws.
func (h *GameHandler) HandleSimulateNextDay(w http.ResponseWriter, r *http.Request) {
	_, _ = h.store.AdvanceToNextDay()
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleResetDay1 resets customs game to Day 1 starter state.
func (h *GameHandler) HandleResetDay1(w http.ResponseWriter, r *http.Request) {
	_ = h.store.ResetToDay1()
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
