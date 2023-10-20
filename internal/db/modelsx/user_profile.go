package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var profileUpdate = boil.Whitelist(
	models.UserProfileColumns.Timezone,
	models.UserProfileColumns.UpdatedAt,
)

// FetchUserProfile fetches a user profile from the database.
func FetchUserProfile(ctx context.Context, exec boil.ContextExecutor, id string) (*models.UserProfile, error) {
	return models.UserProfiles(qm.Where("user_id = ?", id)).One(ctx, exec)
}

// UpsertUserProfile inserts a user profile or updates it on conflict.
func UpsertUserProfile(ctx context.Context, exec boil.ContextExecutor, p *models.UserProfile) error {
	return p.Upsert(ctx, exec, true, []string{"user_id"}, profileUpdate, boil.Infer())
}
