package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var (
	noProfileError = &api.InteractionResponseData{
		Content: option.NewNullableString("You need to initialize a profile with `/profile init` before running this command"),
		Flags:   discord.EphemeralMessage,
	}
)

func (b *Bot) cmdProfileInit(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	id := data.Event.SenderID()
	if id == 0 {
		ctxlog.Error(ctx, "sender id is 0")
		return respondError("An unexpected error has occurred, oops!")
	}

	exists, err := models.UserProfiles(qm.Where("user_id = ?", id.String())).Exists(ctx, b.db)
	if err != nil {
		ctxlog.Error(ctx, "error fetching user profile", zap.Error(err))
		return dbError
	}

	if exists {
		return respond("You have already created a profile")
	}

	p := &models.UserProfile{
		UserID: id.String(),
	}

	if err := modelsx.UpsertUserProfile(ctx, b.db, p); err != nil {
		ctxlog.Error(ctx, "error upserting user profile", zap.Error(err))
		return dbError
	}

	return respond("Your profile has been created")
}

func (b *Bot) cmdProfileDelete(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	confirmation, err := data.Options.Find("are-you-sure").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing value as boolean", zap.Error(err))
		return respondError("Error parsing confirmation as a boolean.. hehe")
	}

	if !confirmation {
		return respond("Stop wasting my time...")
	}

	id := data.Event.SenderID()
	if id == 0 {
		ctxlog.Error(ctx, "sender id is 0")
		return respondError("An unexpected error has occurred, oops!")
	}

	p, err := modelsx.FetchUserProfile(ctx, b.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respond("Can't delete a profile if you don't have one...")
		}

		ctxlog.Error(ctx, "error fetching profile", zap.Error(err))
		return dbError
	}

	if err := p.Delete(ctx, b.db); err != nil {
		ctxlog.Error(ctx, "error deleting profile", zap.Error(err))
		return dbError
	}

	return respond("Your profile has been deleted")
}

func (b *Bot) cmdProfileGet(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	id := data.Event.SenderID()
	if id == 0 {
		ctxlog.Error(ctx, "sender id is 0")
		return respondError("Unexpected error has occurred, oops!")
	}

	exists, err := modelsx.UserProfileExists(ctx, b.db, id.String())
	if err != nil {
		ctxlog.Error(ctx, "error checking if profile exists", zap.Error(err))
		return dbError
	}

	if !exists {
		return noProfileError
	}

	addr := fmt.Sprintf("%s/profile", b.siteAddress)
	return respondf("You can find your profile at <%s>", addr)
}

func (b *Bot) cmdProfileSetTimezone(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	tz := data.Options.Find("timezone").String()
	if tz == "" {
		return respondError("Timezone cannot be blank!!")
	}

	if !validTimezone(tz) {
		return respondError("Timezone must be in **Area/Location** format e.g. **America/Phoenix**")
	}

	_, err := time.LoadLocation(tz)
	if err != nil {
		return respondError("Unable to parse given timezone. Are you sure it's real?")
	}

	id := data.Event.SenderID()
	if id == 0 {
		ctxlog.Error(ctx, "sender id is 0")
		return respondError("Unexpected error has occurred, oops!")
	}

	p, err := modelsx.FetchUserProfile(ctx, b.db, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return noProfileError
		}

		ctxlog.Error(ctx, "error fetching profile", zap.Error(err))
		return dbError
	}

	p.Timezone = null.StringFrom(tz)
	if err := modelsx.UpsertUserProfile(ctx, b.db, p); err != nil {
		return dbError
	}

	return respondf("Your timezone has been set to `%s`", tz)
}
