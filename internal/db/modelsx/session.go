package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var sessionUpdate = boil.Whitelist(
	models.SessionColumns.Data,
	models.SessionColumns.Expiry,
)

// UpsertSession inserts a session, or upserts it on conflict.
func UpsertSession(ctx context.Context, exec boil.ContextExecutor, s *models.Session) error {
	return s.Upsert(ctx, exec, true, []string{models.SessionColumns.Token}, sessionUpdate, boil.Infer())
}
