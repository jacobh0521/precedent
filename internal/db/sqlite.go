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
	AuditResults []rules.RuleAudit `json:"audit_results"`
	CreatedAt    time.Time         `json:"created_at"`
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
		selected_manifest_id TEXT NOT NULL DEFAULT 'cargo-giraffe',
		attempt_count INTEGER NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS manifests (
		id TEXT PRIMARY KEY,
		day INTEGER NOT NULL,
		name TEXT NOT NULL,
		traits TEXT NOT NULL,
		description TEXT NOT NULL,
		flavor_text TEXT NOT NULL
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
		day INTEGER NOT NULL,
		adjective TEXT NOT NULL,
		container TEXT NOT NULL,
		attempt INTEGER NOT NULL DEFAULT 1,
		status TEXT NOT NULL,
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
		_, err = s.db.Exec(`INSERT INTO game_state (id, current_day, selected_manifest_id, attempt_count) VALUES (1, 1, 'cargo-giraffe', 1)`)
		if err != nil {
			return err
		}
	}

	// Starter Manifests: 3 contraband items
	var mCount int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM manifests`).Scan(&mCount)
	if err != nil {
		return err
	}
	if mCount == 0 {
		starterManifests := []Manifest{
			{
				ID:          "cargo-giraffe",
				Day:         1,
				Name:        "A LIVE GIRAFFE",
				Traits:      []string{"living", "oversized"},
				Description: "A 5.5-meter tall long-necked ungulate.",
				FlavorText:  "Staring at customs agents with complete nonchalance.",
			},
			{
				ID:          "cargo-weather",
				Day:         1,
				Name:        "TOMORROW'S WEATHER",
				Traits:      []string{"wet", "oversized"},
				Description: "A high-pressure storm pocket with 90% rain probability.",
				FlavorText:  "Dripping through standard baggage carousels.",
			},
			{
				ID:          "cargo-regrets",
				Day:         1,
				Name:        "THREE REGRETS",
				Traits:      []string{"wet", "oversized"},
				Description: "Heavy psychic residue soaked in cold tears.",
				FlavorText:  "Deeply awkward to explain at the x-ray machine.",
			},
		}

		for _, m := range starterManifests {
			traitsJSON, _ := json.Marshal(m.Traits)
			_, err = s.db.Exec(
				`INSERT INTO manifests (id, day, name, traits, description, flavor_text) VALUES (?, ?, ?, ?, ?, ?)`,
				m.ID, m.Day, m.Name, string(traitsJSON), m.Description, m.FlavorText,
			)
			if err != nil {
				return err
			}
		}
	}

	// Starter 4-word Signpost Laws
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
	return day, nil
}

// GetAttemptCount returns current attempt number (1-5).
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

// IncrementAttempt increments attempt number up to 5.
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

// GetSelectedManifestID returns the currently selected manifest ID.
func (s *Store) GetSelectedManifestID() (string, error) {
	var id string
	err := s.db.QueryRow(`SELECT selected_manifest_id FROM game_state WHERE id = 1`).Scan(&id)
	if err != nil {
		return "cargo-giraffe", err
	}
	if id == "" {
		return "cargo-giraffe", nil
	}
	return id, nil
}

// SetSelectedManifestID updates the currently selected manifest and resets attempts.
func (s *Store) SetSelectedManifestID(manifestID string) error {
	_, err := s.db.Exec(`UPDATE game_state SET selected_manifest_id = ?, attempt_count = 1 WHERE id = 1`, manifestID)
	return err
}

// GetAllManifests returns all manifests.
func (s *Store) GetAllManifests() ([]Manifest, error) {
	rows, err := s.db.Query(`SELECT id, day, name, traits, description, flavor_text FROM manifests ORDER BY day ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Manifest
	for rows.Next() {
		var m Manifest
		var traitsJSON string
		if err := rows.Scan(&m.ID, &m.Day, &m.Name, &traitsJSON, &m.Description, &m.FlavorText); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(traitsJSON), &m.Traits)
		list = append(list, m)
	}
	return list, nil
}

// GetManifestByID retrieves a specific manifest.
func (s *Store) GetManifestByID(id string) (*Manifest, error) {
	var m Manifest
	var traitsJSON string
	err := s.db.QueryRow(
		`SELECT id, day, name, traits, description, flavor_text FROM manifests WHERE id = ?`,
		id,
	).Scan(&m.ID, &m.Day, &m.Name, &traitsJSON, &m.Description, &m.FlavorText)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal([]byte(traitsJSON), &m.Traits)
	return &m, nil
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

// SaveLaw inserts a new signpost punchline law.
func (s *Store) SaveLaw(l rules.Law) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO laws (day, signpost, rule_type, target_adj, target_cargo, is_active, created_at)
		 VALUES (?, ?, ?, ?, ?, 1, ?)`,
		l.Day, l.Signpost, l.RuleType, l.TargetAdj, l.TargetCargo, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// SaveDeclaration logs an attempt.
func (s *Store) SaveDeclaration(manifestID string, day int, d rules.Disguise, attempt int, verdict rules.AuditVerdict) error {
	auditJSON, _ := json.Marshal(verdict.Audits)

	_, err := s.db.Exec(
		`INSERT INTO declarations (manifest_id, day, adjective, container, attempt, status, audit_results, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		manifestID, day, d.Adjective, d.Container, attempt, verdict.Status, string(auditJSON), time.Now(),
	)
	return err
}

// AdvanceToNextDay increments day, resets attempt count, and generates new day parcels if needed.
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

	nextDayManifests := []Manifest{
		{
			ID:          fmt.Sprintf("cargo-day%d-feline", nextDay),
			Day:         nextDay,
			Name:        "A QUANTUM CAT",
			Traits:      []string{"living", "oversized"},
			Description: "Simultaneously meowing and purring in super-position.",
			FlavorText:  "Station policy: Do not look directly into the carrier.",
		},
		{
			ID:          fmt.Sprintf("cargo-day%d-nebula", nextDay),
			Day:         nextDay,
			Name:        "A BOTTLED NEBULA",
			Traits:      []string{"wet", "oversized"},
			Description: "Stellar dust cloud condensing into ionized precipitation.",
			FlavorText:  "Declared value: One small galaxy.",
		},
	}

	for _, m := range nextDayManifests {
		var exists int
		_ = s.db.QueryRow(`SELECT COUNT(*) FROM manifests WHERE id = ?`, m.ID).Scan(&exists)
		if exists == 0 {
			traitsJSON, _ := json.Marshal(m.Traits)
			_, _ = s.db.Exec(
				`INSERT INTO manifests (id, day, name, traits, description, flavor_text) VALUES (?, ?, ?, ?, ?, ?)`,
				m.ID, m.Day, m.Name, string(traitsJSON), m.Description, m.FlavorText,
			)
		}
	}

	if len(nextDayManifests) > 0 {
		_ = s.SetSelectedManifestID(nextDayManifests[0].ID)
	}

	return nextDay, nil
}

// ResetToDay1 resets game state and laws back to starter 5 rules.
func (s *Store) ResetToDay1() error {
	_, err := s.db.Exec(`UPDATE game_state SET current_day = 1, selected_manifest_id = 'cargo-giraffe', attempt_count = 1 WHERE id = 1`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM laws WHERE day > 1`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM declarations`)
	return err
}
