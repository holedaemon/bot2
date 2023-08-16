package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var tokenUpdate = boil.Whitelist(
	models.DiscordTokenColumns.UpdatedAt,
	models.DiscordTokenColumns.AccessToken,
	models.DiscordTokenColumns.RefreshToken,
	models.DiscordTokenColumns.TokenType,
	models.DiscordTokenColumns.Expiry,
)

// UpsertDiscordToken inserts a Discord OAuth2 token, or upserts it on conflict.
func UpsertDiscordToken(ctx context.Context, exec boil.ContextExecutor, t *models.DiscordToken) error {
	return t.Upsert(ctx, exec, true, []string{models.DiscordTokenColumns.UserID}, tokenUpdate, boil.Infer())
}
