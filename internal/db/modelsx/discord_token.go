package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/oauth2"
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

// FetchToken fetches a Discord OAuth2 token by its associated user ID.
func FetchToken(ctx context.Context, exec boil.ContextExecutor, id string) (*models.DiscordToken, error) {
	return models.DiscordTokens(qm.Where("user_id = ?", id)).One(ctx, exec)
}

// ModelToToken converts a stored Discord Token into an *oauth2.Token.
func ModelToToken(dt *models.DiscordToken) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  dt.AccessToken,
		TokenType:    dt.TokenType,
		RefreshToken: dt.RefreshToken,
		Expiry:       dt.Expiry,
	}
}

// TokenToModel converts an OAuth2 token and user id to a Discord token.
func TokenToModel(id string, tok *oauth2.Token) *models.DiscordToken {
	return &models.DiscordToken{
		UserID:       id,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
		TokenType:    tok.TokenType,
	}
}
