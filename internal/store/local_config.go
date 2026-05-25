package store

import (
	"database/sql"

	"github.com/Mininglamp-OSS/octo-speech/internal/config"
)

type LocalASRConfig struct {
	Enabled       bool
	TimeoutMs     int
	ProbeURL      string
	TranscribeURL string
}

type LocalConfigStore struct {
	db  *sql.DB
	cfg *config.Config
}

func NewLocalConfigStore(db *sql.DB, cfg *config.Config) *LocalConfigStore {
	return &LocalConfigStore{db: db, cfg: cfg}
}

func (s *LocalConfigStore) Upsert(appID, subjectID, scopeType, scopeID string, enabled bool, timeoutMs *int, probeURL, transcribeURL *string) error {
	_, err := s.db.Exec(
		`INSERT INTO local_asr_config (app_id, subject_id, scope_type, scope_id, enabled, timeout_ms, probe_url, transcribe_url)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?) AS new
		 ON DUPLICATE KEY UPDATE enabled=new.enabled, timeout_ms=new.timeout_ms, probe_url=new.probe_url, transcribe_url=new.transcribe_url`,
		appID, subjectID, scopeType, scopeID, enabled, timeoutMs, probeURL, transcribeURL,
	)
	return err
}

func (s *LocalConfigStore) Delete(appID, subjectID, scopeType, scopeID string) error {
	_, err := s.db.Exec(
		`DELETE FROM local_asr_config
		 WHERE app_id = ? AND subject_id = ? AND scope_type = ? AND scope_id = ?`,
		appID, subjectID, scopeType, scopeID,
	)
	return err
}

func (s *LocalConfigStore) Query(appID, subjectID, scopeType, scopeID string) *LocalASRConfig {
	result := &LocalASRConfig{
		Enabled:       s.cfg.LocalEnabled,
		TimeoutMs:     s.cfg.LocalTimeoutMs,
		ProbeURL:      s.cfg.LocalProbeURL,
		TranscribeURL: s.cfg.LocalTranscribeURL,
	}

	if appID == "" || subjectID == "" || scopeType == "" || scopeID == "" {
		return result
	}

	var enabled sql.NullInt64
	var timeoutMs sql.NullInt64
	var probeURL sql.NullString
	var transcribeURL sql.NullString

	err := s.db.QueryRow(
		`SELECT enabled, timeout_ms, probe_url, transcribe_url
		 FROM local_asr_config
		 WHERE app_id = ? AND subject_id = ? AND scope_type = ? AND scope_id = ?`,
		appID, subjectID, scopeType, scopeID,
	).Scan(&enabled, &timeoutMs, &probeURL, &transcribeURL)

	if err != nil {
		return result
	}

	if enabled.Valid {
		result.Enabled = enabled.Int64 == 1
	}
	if timeoutMs.Valid {
		result.TimeoutMs = int(timeoutMs.Int64)
	}
	if probeURL.Valid {
		result.ProbeURL = probeURL.String
	}
	if transcribeURL.Valid {
		result.TranscribeURL = transcribeURL.String
	}

	return result
}
