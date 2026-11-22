package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
	"precedent/internal/rules"
)

// Manifest represents a contraband parcel.
type Manifest struct {
	ID          string   `json:"id"`
	Day         int      `json:"day"`
	Name        string   `json:"name"`
	Traits      []string `json:"traits"`
	Description string   `json:"description"`
	FlavorText  string   `json:"flavor_text"`
}

// ToCargo converts Manifest to rules.Cargo.
func (m Manifest) ToCargo() rules.Cargo {
	return rules.Cargo{
		ID:     m.ID,
		Name:   m.Name,
		Traits: m.Traits,
	}
}

// DeclarationRecord logs past disguise attempts.
type DeclarationRecord struct {
	ID           int64             `json:"id"`
	ManifestID   string            `json:"manifest_id"`
	ManifestName string            `json:"manifest_name"`
	Day          int               `json:"day"`
	Adjective    string            `json:"adjective"`
	Container    string            `json:"container"`
	Attempt      int               `json:"attempt"`
	Status       string            `json:"status"` // "CLEARED" or "CONFISCATED"
	EmojiRow     string            `json:"emoji_row"`
	AuditResults []rules.RuleAudit `json:"audit_results"`
	CreatedAt    time.Time         `json:"created_at"`
}

// DaySummary holds summary of a day's smuggling attempt.
type DaySummary struct {
	Day          int      `json:"day"`
	ManifestName string   `json:"manifest_name"`
	Cleared      bool     `json:"cleared"`
	AttemptsUsed int      `json:"attempts_used"`
	WinningPath  string   `json:"winning_path"` // e.g. "Metaphorical in Envelope"
	Attempts     []string `json:"attempts"`     // List of "🟩🟩🟩🟩🟩 (Metaphorical + Envelope)"
}

// WeeklyDebrief holds aggregate weekly player performance.
type WeeklyDebrief struct {
	TotalDays     int          `json:"total_days"`
	ClearedDays   int          `json:"cleared_days"`
	SuccessRate   int          `json:"success_rate"`
	TotalAttempts int          `json:"total_attempts"`
	DaySummaries  []DaySummary `json:"day_summaries"`
	ShareText     string       `json:"share_text"`
}

// Store handles SQLite operations.
type Store struct {
	db *sql.DB
}

// NewStore initializes database schema and seeds initial data.
func NewStore(filepath string) (*Store, error) {
	conn, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	conn.SetMaxOpenConns(1)

	s := &Store{db: conn}
	if err := s.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := s.seedDay1(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("seeding failed: %w", err)
	}

	if err := s.EnsureRegimeConsistency(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("regime consistency check failed: %w", err)
	}

	return s, nil
}

// Close closes the underlying SQLite DB.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS game_state (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		current_day INTEGER NOT NULL DEFAULT 1,
		attempt_count INTEGER NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS laws (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		day INTEGER NOT NULL,
		signpost TEXT NOT NULL,
		rule_type TEXT NOT NULL,
		target_adj TEXT DEFAULT '',
		target_cargo TEXT DEFAULT '',
		is_active INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS declarations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		manifest_id TEXT NOT NULL,
		manifest_name TEXT NOT NULL,
		day INTEGER NOT NULL,
		adjective TEXT NOT NULL,
		container TEXT NOT NULL,
		attempt INTEGER NOT NULL DEFAULT 1,
		status TEXT NOT NULL,
		emoji_row TEXT NOT NULL DEFAULT '',
		audit_results TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) seedDay1() error {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM game_state`).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = s.db.Exec(`INSERT INTO game_state (id, current_day, attempt_count) VALUES (1, 1, 1)`)
		if err != nil {
			return err
		}
	}

	var lCount int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM laws`).Scan(&lCount)
	if err != nil {
		return err
	}
	if lCount == 0 {
		starterLaws := rules.StarterLaws()
		for _, law := range starterLaws {
			_, err = s.db.Exec(
				`INSERT INTO laws (day, signpost, rule_type, target_adj, target_cargo, is_active, created_at)
				 VALUES (?, ?, ?, '', '', 1, ?)`,
				law.Day, law.Signpost, law.RuleType, time.Now(),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetCurrentDay returns current in-game day.
func (s *Store) GetCurrentDay() (int, error) {
	var day int
	err := s.db.QueryRow(`SELECT current_day FROM game_state WHERE id = 1`).Scan(&day)
	if err != nil {
		return 1, err
	}
	if day < 1 {
		day = 1
	}
	return day, nil
}

// GetTodayManifest returns the single unique manifest for today.
func (s *Store) GetTodayManifest() (Manifest, error) {
	day, err := s.GetCurrentDay()
	if err != nil {
		day = 1
	}
	return GetManifestForDay(day), nil
}

// GetAttemptCount returns current attempt number (1-3).
func (s *Store) GetAttemptCount() (int, error) {
	var attempt int
	err := s.db.QueryRow(`SELECT attempt_count FROM game_state WHERE id = 1`).Scan(&attempt)
	if err != nil {
		return 1, err
	}
	if attempt < 1 {
		attempt = 1
	}
	return attempt, nil
}

// IncrementAttempt increments attempt number up to 3.
func (s *Store) IncrementAttempt() (int, error) {
	att, _ := s.GetAttemptCount()
	att++
	_, err := s.db.Exec(`UPDATE game_state SET attempt_count = ? WHERE id = 1`, att)
	return att, err
}

// ResetAttempt resets attempt count back to 1.
func (s *Store) ResetAttempt() error {
	_, err := s.db.Exec(`UPDATE game_state SET attempt_count = 1 WHERE id = 1`)
	return err
}

// GetActiveLaws retrieves all active signpost rules.
func (s *Store) GetActiveLaws() ([]rules.Law, error) {
	rows, err := s.db.Query(
		`SELECT id, day, signpost, rule_type, target_adj, target_cargo, created_at 
		 FROM laws WHERE is_active = 1 ORDER BY day ASC, id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var laws []rules.Law
	for rows.Next() {
		var l rules.Law
		if err := rows.Scan(
			&l.ID, &l.Day, &l.Signpost, &l.RuleType, &l.TargetAdj, &l.TargetCargo, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		laws = append(laws, l)
	}
	return laws, nil
}

// SaveLaws inserts new signpost laws.
func (s *Store) SaveLaws(laws []rules.Law) error {
	for _, l := range laws {
		_, err := s.db.Exec(
			`INSERT INTO laws (day, signpost, rule_type, target_adj, target_cargo, is_active, created_at)
			 VALUES (?, ?, ?, ?, ?, 1, ?)`,
			l.Day, l.Signpost, l.RuleType, l.TargetAdj, l.TargetCargo, time.Now(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveDeclaration logs an attempt.
func (s *Store) SaveDeclaration(manifestID, manifestName string, day int, d rules.Disguise, attempt int, verdict rules.AuditVerdict) error {
	auditJSON, _ := json.Marshal(verdict.Audits)

	_, err := s.db.Exec(
		`INSERT INTO declarations (manifest_id, manifest_name, day, adjective, container, attempt, status, emoji_row, audit_results, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		manifestID, manifestName, day, d.Adjective, d.Container, attempt, verdict.Status, verdict.EmojiRow, string(auditJSON), time.Now(),
	)
	return err
}

// GetDeclarationsForDay retrieves declarations for a given day.
func (s *Store) GetDeclarationsForDay(day int) ([]DeclarationRecord, error) {
	rows, err := s.db.Query(
		`SELECT id, manifest_id, manifest_name, day, adjective, container, attempt, status, emoji_row, audit_results, created_at
		 FROM declarations
		 WHERE day = ?
		 ORDER BY attempt ASC, id ASC`,
		day,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []DeclarationRecord
	for rows.Next() {
		var r DeclarationRecord
		var auditJSON string
		if err := rows.Scan(
			&r.ID, &r.ManifestID, &r.ManifestName, &r.Day, &r.Adjective, &r.Container, &r.Attempt, &r.Status, &r.EmojiRow, &auditJSON, &r.CreatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(auditJSON), &r.AuditResults)
		records = append(records, r)
	}
	return records, nil
}

// GetWeeklyDebrief aggregates customs attempts across days.
func (s *Store) GetWeeklyDebrief() (WeeklyDebrief, error) {
	rows, err := s.db.Query(
		`SELECT day, manifest_id, manifest_name, adjective, container, attempt, status, emoji_row
		 FROM declarations
		 ORDER BY day ASC, attempt ASC`,
	)
	if err != nil {
		return WeeklyDebrief{}, err
	}
	defer rows.Close()

	dayMap := make(map[int]*DaySummary)
	dayOrder := make([]int, 0)
	totalAttempts := 0

	for rows.Next() {
		var day, attempt int
		var mID, mName, adj, cnt, status, emojiRow string
		if err := rows.Scan(&day, &mID, &mName, &adj, &cnt, &attempt, &status, &emojiRow); err != nil {
			continue
		}

		totalAttempts++
		summary, exists := dayMap[day]
		if !exists {
			summary = &DaySummary{
				Day:          day,
				ManifestName: mName,
				Attempts:     make([]string, 0),
			}
			dayMap[day] = summary
			dayOrder = append(dayOrder, day)
		}

		line := fmt.Sprintf("%s (%s + %s)", emojiRow, adj, cnt)
		summary.Attempts = append(summary.Attempts, line)
		summary.AttemptsUsed = attempt

		if status == "CLEARED" {
			summary.Cleared = true
			summary.WinningPath = fmt.Sprintf("%s in %s", adj, cnt)
		}
	}

	clearedCount := 0
	summaries := make([]DaySummary, 0, len(dayOrder))
	for _, day := range dayOrder {
		ds := *dayMap[day]
		if ds.Cleared {
			clearedCount++
		}
		summaries = append(summaries, ds)
	}

	totalDays := len(summaries)
	if totalDays == 0 {
		totalDays = 1
	}
	rate := (clearedCount * 100) / totalDays

	shareText := fmt.Sprintf("PRECEDENT Weekly Customs Debrief 📑\n")
	shareText += fmt.Sprintf("Success Rate: %d/%d Days Cleared (%d%%)\n\n", clearedCount, totalDays, rate)
	for _, ds := range summaries {
		statusMark := "❌"
		if ds.Cleared {
			statusMark = "✅"
		}
		shareText += fmt.Sprintf("Day %d %s [Guess %d/%d]: %s\n", ds.Day, statusMark, ds.AttemptsUsed, rules.MaxDailyAttempts, ds.ManifestName)
		if ds.Cleared {
			shareText += fmt.Sprintf("  ↳ Slipped past as: %s\n", ds.WinningPath)
		}
		for _, att := range ds.Attempts {
			shareText += fmt.Sprintf("  %s\n", att)
		}
		shareText += "\n"
	}
	shareText += "Play Precedent: https://precedent.cloud.run"

	return WeeklyDebrief{
		TotalDays:     totalDays,
		ClearedDays:   clearedCount,
		SuccessRate:   rate,
		TotalAttempts: totalAttempts,
		DaySummaries:  summaries,
		ShareText:     shareText,
	}, nil
}

// EnsureRegimeConsistency ensures the active laws strictly match the current day's regime,
// purging all past precedent laws whenever a new 7-day regime cycle begins.
func (s *Store) EnsureRegimeConsistency() error {
	day, err := s.GetCurrentDay()
	if err != nil {
		return err
	}
	regime := rules.GetRegimeForDay(day)
	cycleStartDay := ((day - 1) / 7) * 7 + 1

	// If we are on Day 8+ (or any cycle > 1), purge any lingering precedent laws from earlier cycles
	if day > 1 && cycleStartDay > 1 {
		var oldCount int
		_ = s.db.QueryRow(`SELECT COUNT(*) FROM laws WHERE is_active = 1 AND day < ?`, cycleStartDay).Scan(&oldCount)
		if oldCount > 0 {
			_, _ = s.db.Exec(`UPDATE laws SET is_active = 0 WHERE is_active = 1 AND day < ?`, cycleStartDay)
		}

		// Ensure current regime's 5 starter codex laws exist for this cycle
		var codexCount int
		_ = s.db.QueryRow(`SELECT COUNT(*) FROM laws WHERE is_active = 1 AND day >= ?`, cycleStartDay).Scan(&codexCount)
		if codexCount == 0 {
			for _, law := range regime.Codex {
				_, _ = s.db.Exec(
					`INSERT INTO laws (day, signpost, rule_type, target_adj, target_cargo, is_active, created_at)
					 VALUES (?, ?, ?, '', '', 1, ?)`,
					cycleStartDay, law.Signpost, law.RuleType, time.Now(),
				)
			}
		}
	}
	return nil
}

// AdvanceToNextDay increments day, resets attempts to 1, and triggers Regime Shift on Day 8/15/22...
func (s *Store) AdvanceToNextDay() (int, error) {
	currentDay, err := s.GetCurrentDay()
	if err != nil {
		currentDay = 1
	}
	nextDay := currentDay + 1

	_, err = s.db.Exec(`UPDATE game_state SET current_day = ?, attempt_count = 1 WHERE id = 1`, nextDay)
	if err != nil {
		return currentDay, err
	}

	// Day 8 (and every 7-day cycle) Regime Shift:
	// Purge all precedent laws and install 5 fresh starter codex laws for the new regime!
	if rules.IsCycleStartDay(nextDay) {
		_, err = s.db.Exec(`UPDATE laws SET is_active = 0 WHERE is_active = 1`)
		if err != nil {
			return nextDay, err
		}

		regime := rules.GetRegimeForDay(nextDay)
		for _, law := range regime.Codex {
			_, err = s.db.Exec(
				`INSERT INTO laws (day, signpost, rule_type, target_adj, target_cargo, is_active, created_at)
				 VALUES (?, ?, ?, '', '', 1, ?)`,
				nextDay, law.Signpost, law.RuleType, time.Now(),
			)
			if err != nil {
				return nextDay, err
			}
		}
	}

	return nextDay, nil
}

// ResetToDay1 resets game state and laws back to starter 5 rules.
func (s *Store) ResetToDay1() error {
	_, err := s.db.Exec(`UPDATE game_state SET current_day = 1, attempt_count = 1 WHERE id = 1`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM laws`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM declarations`)
	if err != nil {
		return err
	}

	// Re-seed Day 1 starter laws
	starterLaws := rules.StarterLaws()
	for _, law := range starterLaws {
		_, err = s.db.Exec(
			`INSERT INTO laws (day, signpost, rule_type, target_adj, target_cargo, is_active, created_at)
			 VALUES (?, ?, ?, '', '', 1, ?)`,
			1, law.Signpost, law.RuleType, time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
