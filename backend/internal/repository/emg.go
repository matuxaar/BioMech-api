package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/model"
)

type EMGRepository struct {
	db *pgxpool.Pool
}

func NewEMGRepository(db *pgxpool.Pool) *EMGRepository {
	return &EMGRepository{db: db}
}

func (r *EMGRepository) CreateSession(ctx context.Context, userID string, req *model.CreateEMGSessionRequest) (*model.EMGSession, error) {
	session := &model.EMGSession{
		ID:        uuid.New().String(),
		UserID:    userID,
		DeviceID:  req.DeviceID,
		Label:     req.Label,
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO emg_sessions (id, user_id, device_id, label, started_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		session.ID, session.UserID, session.DeviceID, session.Label, session.StartedAt, session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *EMGRepository) EndSession(ctx context.Context, sessionID string) error {
	now := time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE emg_sessions SET ended_at = $1 WHERE id = $2`, now, sessionID,
	)
	return err
}

func (r *EMGRepository) FindSessionsByUserID(ctx context.Context, userID string) ([]model.EMGSession, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, device_id, COALESCE(label, ''), started_at, ended_at, created_at
		 FROM emg_sessions WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.EMGSession
	for rows.Next() {
		var s model.EMGSession
		if err := rows.Scan(&s.ID, &s.UserID, &s.DeviceID, &s.Label, &s.StartedAt, &s.EndedAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (r *EMGRepository) FindSessionByID(ctx context.Context, id string) (*model.EMGSession, error) {
	s := &model.EMGSession{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, device_id, COALESCE(label, ''), started_at, ended_at, created_at
		 FROM emg_sessions WHERE id = $1`, id,
	).Scan(&s.ID, &s.UserID, &s.DeviceID, &s.Label, &s.StartedAt, &s.EndedAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *EMGRepository) AddSample(ctx context.Context, sessionID string, req *model.AddSampleRequest) (*model.EMGSample, error) {
	sample := &model.EMGSample{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Timestamp: req.Timestamp,
		Channel1:  req.Channel1,
		Channel2:  req.Channel2,
		Channel3:  req.Channel3,
		Channel4:  req.Channel4,
		Channel5:  req.Channel5,
		Channel6:  req.Channel6,
		Channel7:  req.Channel7,
		Channel8:  req.Channel8,
		Metadata:  req.Metadata,
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO emg_samples (id, session_id, timestamp, channel_1, channel_2, channel_3, channel_4, channel_5, channel_6, channel_7, channel_8, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		sample.ID, sample.SessionID, sample.Timestamp,
		sample.Channel1, sample.Channel2, sample.Channel3, sample.Channel4,
		sample.Channel5, sample.Channel6, sample.Channel7, sample.Channel8,
		sample.Metadata,
	)
	if err != nil {
		return nil, err
	}

	return sample, nil
}

func (r *EMGRepository) AddSamplesBatch(ctx context.Context, sessionID string, samples []model.AddSampleRequest) error {
	batch := &pgx.Batch{}
	for _, s := range samples {
		batch.Queue(
			`INSERT INTO emg_samples (id, session_id, timestamp, channel_1, channel_2, channel_3, channel_4, channel_5, channel_6, channel_7, channel_8, metadata)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			uuid.New().String(), sessionID, s.Timestamp,
			s.Channel1, s.Channel2, s.Channel3, s.Channel4,
			s.Channel5, s.Channel6, s.Channel7, s.Channel8,
			s.Metadata,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(samples); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (r *EMGRepository) FindSamplesBySessionID(ctx context.Context, sessionID string) ([]model.EMGSample, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, session_id, timestamp, channel_1, channel_2, channel_3, channel_4, channel_5, channel_6, channel_7, channel_8, COALESCE(metadata, '')
		 FROM emg_samples WHERE session_id = $1 ORDER BY timestamp ASC`, sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var samples []model.EMGSample
	for rows.Next() {
		var s model.EMGSample
		if err := rows.Scan(&s.ID, &s.SessionID, &s.Timestamp,
			&s.Channel1, &s.Channel2, &s.Channel3, &s.Channel4,
			&s.Channel5, &s.Channel6, &s.Channel7, &s.Channel8,
			&s.Metadata,
		); err != nil {
			return nil, err
		}
		samples = append(samples, s)
	}

	return samples, nil
}
