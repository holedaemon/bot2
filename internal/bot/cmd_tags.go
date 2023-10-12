package bot

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const tagNameMaxLength = 100

func (b *Bot) cmdTag(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	name := data.Options.Find("name").String()

	if name == "" {
		return respondError("Name cannot be blank!!")
	}

	name = strings.ToLower(name)

	tag, err := modelsx.FetchTag(ctx, b.DB, data.Event.GuildID.String(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondError("There is no tag by that name in this guild!!")
		}

		ctxlog.Error(ctx, "error fetching tag", zap.Error(err))
		return respondError("Error fetching tag!!!")
	}

	return respond(tag.Content)
}

func (b *Bot) cmdTagCreate(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	name := data.Options.Find("name").String()
	content := data.Options.Find("content").String()

	if name == "" || content == "" {
		return respondError("Both name and content must be set!!")
	}

	if len(name) > tagNameMaxLength {
		return respondErrorf("A tag name cannot exceed %d in length", tagNameMaxLength)
	}

	name = strings.ToLower(name)
	ctx = ctxlog.With(ctx, zap.String("trigger", name))

	exists, err := models.Tags(
		qm.Where("guild_id = ? AND trigger = ?",
			data.Event.GuildID.String(),
			name,
		),
	).Exists(ctx, b.DB)
	if err != nil {
		ctxlog.Error(ctx, "error checking if tag exists", zap.Error(err))
		return respondError("Uh uh uh, error checking if tag already exists!!")
	}

	if exists {
		return respond("A tag with that name already exists in this guild!!")
	}

	tag := &models.Tag{
		GuildID:   data.Event.GuildID.String(),
		CreatorID: data.Event.Member.User.ID.String(),
		Editor:    data.Event.Member.User.Tag(),
		Trigger:   name,
		Content:   content,
	}

	if err := tag.Insert(ctx, b.DB, boil.Infer()); err != nil {
		ctxlog.Error(ctx, "error inserting tag into database", zap.Error(err))
		return respondError("Error inserting tag into database!!!")
	}

	return respondf("Tag `%s` has been created", name)
}

func (b *Bot) cmdTagUpdate(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	name := data.Options.Find("name").String()
	content := data.Options.Find("content").String()

	if name == "" || content == "" {
		return respondError("Both name and content must be set!!")
	}

	name = strings.ToLower(name)
	ctx = ctxlog.With(ctx, zap.String("trigger", name))

	tag, err := modelsx.FetchTag(ctx, b.DB, data.Event.GuildID.String(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondError("No tag by that name exists within this guild!!")
		}

		ctxlog.Error(ctx, "error getting tag from database", zap.Error(err))
		return respondError("Error getting tag from database!!!")
	}

	perms, err := b.State.Permissions(data.Event.ChannelID, data.Event.Member.User.ID)
	if err != nil {
		ctxlog.Error(ctx, "error checking permissions for user", zap.String("user_id", data.Event.Member.User.ID.String()))
		return respondError("Error reading your permissions...!")
	}

	if !perms.Has(discord.PermissionManageMessages) && data.Event.Member.User.ID.String() != tag.CreatorID {
		return respondError("You lack permission to update this tag!!")
	}

	if err := modelsx.UpdateTagContent(ctx, b.DB, tag, content, data.Event.Member.User.Tag()); err != nil {
		ctxlog.Error(ctx, "error updating tag content", zap.Error(err))
		return respondError("Error updating tag's content!!!")
	}

	return respondf("Tag `%s` has had its content updated", name)
}

func (b *Bot) cmdTagRename(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	name := data.Options.Find("name").String()
	newName := data.Options.Find("new-name").String()

	if name == "" || newName == "" {
		return respondError("Both name and new-name must be set!!")
	}

	if len(newName) > tagNameMaxLength {
		return respondErrorf("A tag name cannot exceed %d in length", tagNameMaxLength)
	}

	name = strings.ToLower(name)
	newName = strings.ToLower(newName)
	ctx = ctxlog.With(ctx, zap.String("trigger", name))

	tag, err := modelsx.FetchTag(ctx, b.DB, data.Event.GuildID.String(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondError("No tag by that name exists within this guild!!")
		}

		ctxlog.Error(ctx, "error getting tag from database", zap.Error(err))
		return respondError("Error getting tag from database!!!")
	}

	perms, err := b.State.Permissions(data.Event.ChannelID, data.Event.Member.User.ID)
	if err != nil {
		ctxlog.Error(ctx, "error checking permissions for user", zap.String("user_id", data.Event.Member.User.ID.String()))
		return respondError("Error reading your permissions...!")
	}

	if !perms.Has(discord.PermissionManageMessages) && data.Event.Member.User.ID.String() != tag.CreatorID {
		return respondError("You lack permission to update this tag!!")
	}

	exists, err := models.Tags(qm.Where("guild_id = ? AND trigger = ?", data.Event.GuildID.String(), newName)).Exists(ctx, b.DB)
	if err != nil {
		ctxlog.Error(ctx, "error checking if tag exists", zap.Error(err))
		return respondError("Uh uh uh, error checking if tag already exists!!")
	}

	if exists {
		return respondErrorf("There is already a tag called `%s`", newName)
	}

	if err := modelsx.RenameTag(ctx, b.DB, tag, newName, data.Event.Member.User.Tag()); err != nil {
		ctxlog.Error(ctx, "error updating tag name", zap.Error(err))
		return respondError("Error updating tag's name!!!")
	}

	return respondf("Tag `%s` is henceforth known as `%s`", name, newName)
}

func (b *Bot) cmdTagDelete(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	name := data.Options.Find("name").String()
	if name == "" {
		return respondError("Name cannot be blank!!")
	}

	name = strings.ToLower(name)
	ctx = ctxlog.With(ctx, zap.String("trigger", name))

	tag, err := modelsx.FetchTag(ctx, b.DB, data.Event.GuildID.String(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondError("A tag by that name doesn't exist in this guild!!")
		}

		ctxlog.Error(ctx, "error fetching tag", zap.Error(err))
	}

	perms, err := b.State.Permissions(data.Event.ChannelID, data.Event.Member.User.ID)
	if err != nil {
		ctxlog.Error(ctx, "error checking permissions for user", zap.String("user_id", data.Event.Member.User.ID.String()))
		return respondError("Error reading your permissions...!")
	}

	if !perms.Has(discord.PermissionManageMessages) && data.Event.Member.User.ID.String() != tag.CreatorID {
		return respondError("You lack permission to update this tag!!")
	}

	if err := tag.Delete(ctx, b.DB); err != nil {
		ctxlog.Error(ctx, "error deleting tag", zap.Error(err))
		return respondError("Error deleting tag!!")
	}

	return respondf("Tag `%s` has been deleted", name)
}
