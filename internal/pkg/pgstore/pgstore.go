// Package pgstore implements a Postgres-backed session store for alexedwards/scs.
package pgstore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

// Store implements the CtxStore interface provided by scs.
type Store struct {
	db       *sql.DB
	interval time.Duration
}

// New creates a new Store.
func New(db *sql.DB) *Store {
	return NewWithCleanupInterval(db, 5*time.Minute)
}

// NewWithCleanupInterval creates a new Store with the given cleanup interval.
func NewWithCleanupInterval(db *sql.DB, interval time.Duration) *Store {
	return &Store{
		db:       db,
		interval: interval,
	}
}

// Start tells the Store to begin cleaning up expired sessions.
// The Store will stop cleaning up when the given context is cancelled.
func (p *Store) Start(ctx context.Context) {
	go p.startCleanup(ctx)
}

// FindCtx queries the database for an unexpired session.
func (p *Store) FindCtx(ctx context.Context, token string) ([]byte, bool, error) {
	session, err := models.Sessions(qm.Where("token = ?", token), qm.Where("current_timestamp < expiry")).One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return session.Data, true, nil
}

// Find queries the database for an unexpired session.
func (p *Store) Find(token string) ([]byte, bool, error) {
	return p.FindCtx(context.Background(), token)
}

// Commit adds a session to the database.
func (p *Store) CommitCtx(ctx context.Context, token string, b []byte, expiry time.Time) error {
	session := &models.Session{
		Token:  token,
		Data:   b,
		Expiry: expiry,
	}

	return modelsx.UpsertSession(ctx, p.db, session)
}

// Commit adds a session to the database.
func (p *Store) Commit(token string, b []byte, expiry time.Time) error {
	return p.CommitCtx(context.Background(), token, b, expiry)
}

// Delete removes a session from the database.
func (p *Store) DeleteCtx(ctx context.Context, token string) error {
	return models.Sessions(qm.Where("token = ?", token)).DeleteAll(ctx, p.db)
}

// Delete removes a session from the database.
func (p *Store) Delete(token string) error {
	return p.DeleteCtx(context.Background(), token)
}

func (p *Store) deleteExpired(ctx context.Context) error {
	return models.Sessions(qm.Where("expiry < current_timestamp")).DeleteAll(ctx, p.db)
}

func (p *Store) startCleanup(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	for {
		select {
		case <-ticker.C:
			err := p.deleteExpired(ctx)
			if err != nil {
				ctxlog.Error(ctx, "error deleting expired sessions", zap.Error(err))
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
