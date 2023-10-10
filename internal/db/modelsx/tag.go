package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var tagUpdate = boil.Whitelist(models.TagColumns.Trigger, models.TagColumns.Content, models.TagColumns.UpdatedAt)

// FetchTag retrieves a tag from the database.
func FetchTag(ctx context.Context, exec boil.ContextExecutor, guild, name string) (*models.Tag, error) {
	return models.Tags(qm.Where("guild_id = ? AND trigger = ?", guild, name)).One(ctx, exec)
}

// UpdateTagContent updates a tag's content in the database.
func UpdateTagContent(ctx context.Context, exec boil.ContextExecutor, tag *models.Tag, content string) error {
	tag.Content = content
	return tag.Update(ctx, exec, tagUpdate)
}

// RenameTag renames a tag in the database.
func RenameTag(ctx context.Context, exec boil.ContextExecutor, tag *models.Tag, newName string) error {
	tag.Trigger = newName
	return tag.Update(ctx, exec, tagUpdate)
}
