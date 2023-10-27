package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var updaterUpdate = boil.Whitelist(
	models.RoleUpdaterColumns.DoUpdates,
	models.RoleUpdaterColumns.LastTimestamp,
	models.RoleUpdaterColumns.UpdatedAt,
)

// UpsertRoleUpdater inserts a role updater or updates it on conflict.
func UpsertRoleUpdater(ctx context.Context, exec boil.ContextExecutor, s *models.RoleUpdater) error {
	return s.Upsert(ctx, exec, true, []string{"guild_id"}, updaterUpdate, boil.Infer())
}

// FetchRoleUpdater fetches a role updater from the database.
func FetchRoleUpdater(ctx context.Context, exec boil.ContextExecutor, id string) (*models.RoleUpdater, error) {
	return models.RoleUpdaters(qm.Where("guild_id = ?", id)).One(ctx, exec)
}
