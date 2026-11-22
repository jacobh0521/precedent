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

// DeskViewData holds data for the main play view.
type DeskViewData struct {
	CurrentDay        int
	AttemptCount      int
	MaxAttempts       int
	TodayManifest     db.Manifest
	ActiveLaws        []rules.Law
	WeeklyDebrief     db.WeeklyDebrief
	IsRegimeShift     bool
	RegimeShiftBanner string
	RegimeName        string
	RegimeNumber      int
	CycleStartDay     int
}

// VerdictPartialData holds data passed to verdict.html.
type VerdictPartialData struct {
	Verdict          rules.AuditVerdict
	Manifest         db.Manifest
	Disguise         rules.Disguise
	DisguiseSentence string
	ActiveLaws       []rules.Law
	CurrentDay       int
	OutOfTries       bool
	DailyShareText   string
	IsRegimeShift    bool
	RegimeNumber     int
	RegimeName       string
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
	r.Get("/weekly-debrief", h.HandleWeeklyDebrief)
	r.Post("/retry", h.HandleRetryAttempts)
	r.Post("/dev/next-day", h.HandleSimulateNextDay)
	r.Post("/dev/reset", h.HandleResetDay1)
}

// HandleDeskView renders today's single contraband checkpoint.
func (h *GameHandler) HandleDeskView(w http.ResponseWriter, r *http.Request) {
	currentDay, _ := h.store.GetCurrentDay()
	attemptCount, _ := h.store.GetAttemptCount()
	todayManifest, _ := h.store.GetTodayManifest()
	activeLaws, _ := h.store.GetActiveLaws()
	weeklyDebrief, _ := h.store.GetWeeklyDebrief()

	regime := rules.GetRegimeForDay(currentDay)
	isRegimeShift := rules.IsCycleStartDay(currentDay)
	cycleStartDay := ((currentDay - 1) / 7) * 7 + 1

	data := DeskViewData{
		CurrentDay:        currentDay,
		AttemptCount:      attemptCount,
		MaxAttempts:       rules.MaxDailyAttempts, // Synced to constant
		TodayManifest:     todayManifest,
		ActiveLaws:        activeLaws,
		WeeklyDebrief:     weeklyDebrief,
		IsRegimeShift:     isRegimeShift,
		RegimeShiftBanner: "REGIME SHIFT // CODEX OVERHAULED & PRECEDENT LAWS PURGED",
		RegimeName:        regime.Name,
		RegimeNumber:      regime.Number,
		CycleStartDay:     cycleStartDay,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

// HandleSubmitDisguise evaluates user-entered disguise against active laws.
func (h *GameHandler) HandleSubmitDisguise(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Malformed form", http.StatusBadRequest)
		return
	}

	adjective := strings.TrimSpace(r.FormValue("adjective"))
	container := strings.TrimSpace(r.FormValue("container"))

	if adjective == "" {
		adjective = "Unmodified"
	}
	if container == "" {
		container = "Bare Hands"
	}

	currentDay, _ := h.store.GetCurrentDay()
	attemptCount, _ := h.store.GetAttemptCount()
	manifest, _ := h.store.GetTodayManifest()
	activeLaws, _ := h.store.GetActiveLaws()

	disguise := rules.Disguise{
		Adjective: adjective,
		Container: container,
	}

	cargo := manifest.ToCargo()
	verdict := rules.Evaluate(activeLaws, cargo, disguise, attemptCount)

	// If Cleared: Compile TWO new rules with solvability check for tomorrow's cargo
	if verdict.Passed {
		tomorrowManifest := db.GetManifestForDay(currentDay + 1)
		newLaws := rules.CompileLoopholeRules(currentDay+1, disguise, activeLaws, tomorrowManifest.ToCargo())
		if err := h.store.SaveLaws(newLaws); err == nil {
			verdict.CompiledLaws = newLaws
		}
	} else {
		// Increment attempts (up to MaxDailyAttempts)
		if attemptCount < rules.MaxDailyAttempts {
			_, _ = h.store.IncrementAttempt()
		}
	}

	_ = h.store.SaveDeclaration(manifest.ID, manifest.Name, currentDay, disguise, attemptCount, verdict)
	updatedLaws, _ := h.store.GetActiveLaws()

	outOfTries := (!verdict.Passed && attemptCount >= rules.MaxDailyAttempts)

	// Fetch all attempts today for this manifest to construct the Wordle-like share text
	records, _ := h.store.GetDeclarationsForDay(currentDay)
	shareText := formatDailyShare(currentDay, manifest.Name, verdict.Passed, attemptCount, records, verdict.CompiledLaws)
	sentence := rules.FormatDisguiseSentence(disguise.Adjective, manifest.Name, disguise.Container)
	regime := rules.GetRegimeForDay(currentDay)

	partialData := VerdictPartialData{
		Verdict:          verdict,
		Manifest:         manifest,
		Disguise:         disguise,
		DisguiseSentence: sentence,
		ActiveLaws:       updatedLaws,
		CurrentDay:       currentDay,
		OutOfTries:       outOfTries,
		DailyShareText:   shareText,
		IsRegimeShift:    rules.IsCycleStartDay(currentDay),
		RegimeNumber:     regime.Number,
		RegimeName:       regime.Name,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		_ = h.templates.ExecuteTemplate(w, "verdict.html", partialData)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func formatDailyShare(day int, cargoName string, cleared bool, attempts int, records []db.DeclarationRecord, compiledLaws []rules.Law) string {
	score := fmt.Sprintf("%d/%d", attempts, rules.MaxDailyAttempts)
	if !cleared && attempts >= rules.MaxDailyAttempts {
		score = fmt.Sprintf("X/%d", rules.MaxDailyAttempts)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PRECEDENT Day %d (%s)\n", day, score))
	sb.WriteString(fmt.Sprintf("Contraband: %s\n\n", cargoName))

	for i, rec := range records {
		guessSentence := rules.FormatDisguiseSentence(rec.Adjective, cargoName, rec.Container)
		sb.WriteString(fmt.Sprintf("Guess %d/%d: %s (%s)\n", i+1, rules.MaxDailyAttempts, rec.EmojiRow, guessSentence))
	}

	if len(compiledLaws) > 0 {
		signposts := make([]string, 0, len(compiledLaws))
		for _, l := range compiledLaws {
			signposts = append(signposts, l.Signpost)
		}
		sb.WriteString(fmt.Sprintf("\n⚖️ New Laws: %s\n", strings.Join(signposts, " | ")))
	}
	sb.WriteString("\nPlay Precedent: https://precedent.cloud.run")
	return sb.String()
}

// HandleWeeklyDebrief displays the week's end customs report modal.
func (h *GameHandler) HandleWeeklyDebrief(w http.ResponseWriter, r *http.Request) {
	debrief, _ := h.store.GetWeeklyDebrief()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.templates.ExecuteTemplate(w, "weekly_debrief.html", debrief)
}

// HandleRetryAttempts resets attempt count to 1 for the current day.
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
