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

// 8 Adjective Chips & 8 Container Chips
var (
	AdjectiveChips = []string{
		"Inflatable",
		"Dehydrated",
		"Folded",
		"Metaphorical",
		"Frozen",
		"Microscopic",
		"Ceremonial",
		"Rental",
	}

	ContainerChips = []string{
		"Envelope",
		"Pocket",
		"Thermos",
		"Cardboard Box",
		"Velvet Pouch",
		"Glass Jar",
		"Lead Crate",
		"Mesh Bag",
	}
)

// DeskViewData holds data for the main play view.
type DeskViewData struct {
	CurrentDay       int
	AttemptCount     int
	MaxAttempts      int
	Manifests        []db.Manifest
	SelectedManifest *db.Manifest
	ActiveLaws       []rules.Law
	AdjectiveChips   []string
	ContainerChips   []string
	LastVerdict      *rules.AuditVerdict
}

// VerdictPartialData holds data passed to verdict.html.
type VerdictPartialData struct {
	Verdict    rules.AuditVerdict
	Manifest   db.Manifest
	Disguise   rules.Disguise
	ActiveLaws []rules.Law
	CurrentDay int
	OutOfTries bool
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

	data := DeskViewData{
		CurrentDay:       currentDay,
		AttemptCount:     attemptCount,
		MaxAttempts:      5,
		Manifests:        dayManifests,
		SelectedManifest: selectedManifest,
		ActiveLaws:       activeLaws,
		AdjectiveChips:   AdjectiveChips,
		ContainerChips:   ContainerChips,
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

// HandleSubmitDisguise evaluates the 2-click disguise against active laws.
func (h *GameHandler) HandleSubmitDisguise(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Malformed form", http.StatusBadRequest)
		return
	}

	manifestID := r.FormValue("manifest_id")
	adjective := r.FormValue("adjective")
	container := r.FormValue("container")

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

	// If Cleared: Compile punchline law and advance
	if verdict.Passed {
		newLaw := rules.CompilePunchlineLaw(currentDay+1, cargo, disguise, len(activeLaws))
		_, err := h.store.SaveLaw(newLaw)
		if err == nil {
			verdict.Punchline = &newLaw
		}
	} else {
		// Increment attempts if confiscated
		if attemptCount < 5 {
			_, _ = h.store.IncrementAttempt()
		}
	}

	_ = h.store.SaveDeclaration(manifestID, currentDay, disguise, attemptCount, verdict)
	updatedLaws, _ := h.store.GetActiveLaws()

	outOfTries := (!verdict.Passed && attemptCount >= 5)

	partialData := VerdictPartialData{
		Verdict:    verdict,
		Manifest:   *manifest,
		Disguise:   disguise,
		ActiveLaws: updatedLaws,
		CurrentDay: currentDay,
		OutOfTries: outOfTries,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		_ = h.templates.ExecuteTemplate(w, "verdict.html", partialData)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
