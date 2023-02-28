package bot

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) onGuildRoleDelete(e *gateway.GuildRoleDeleteEvent) {
	ctx := context.Background()
	log := b.l.With(zap.String("guild_id", e.GuildID.String()))

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("error creating transaction", zap.Error(err))
		return
	}

	defer func() {
		if err := tx.Commit(); err != nil {
			log.Error("error committing transaction", zap.Error(err))
			return
		}
	}()

	exists, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", e.GuildID.String(), e.RoleID.String())).Exists(ctx, tx)
	if err != nil {
		log.Error("error checking for role", zap.Error(err))
		return
	}

	if !exists {
		return
	}

	log.Info("tracked role has been deleted on Discord, deleting from database...", zap.String("role_id", e.RoleID.String()))

	if err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", e.GuildID.String(), e.RoleID.String())).DeleteAll(ctx, tx); err != nil {
		log.Error("error deleting role", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			log.Error("error rolling back transaction", zap.Error(err))
		}
		return
	}
}

func (b *Bot) cmdRoleCreate(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	name := data.Options.Find("name").String()
	if name == "" {
		return respondError("You've got to provide a name for the role")
	}

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error creating transaction", zap.Error(err))
		return respondError("Unable to create DB transaction, try again later")
	}

	var color int32
	rawColor := data.Options.Find("color").String()
	if rawColor != "" {
		rawColor = strings.TrimPrefix(rawColor, "#")
		rc, err := strconv.ParseInt(rawColor, 16, 32)
		if err != nil {
			ctxlog.Error(ctx, "error converting hex to decimal", zap.Error(err))
			return respondError("The hex you gave me is invalid")
		}

		color = int32(rc)
	}

	hoisted, err := data.Options.Find("hoisted").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error converting option to bool", zap.Error(err))
		return respondError("Unable to parse `hoisted` argument as boolean")
	}

	r, err := b.s.CreateRole(data.Event.GuildID, api.CreateRoleData{
		Name:        name,
		Hoist:       hoisted,
		Mentionable: true,
		Color:       discord.Color(color),
	})
	if err != nil {
		ctxlog.Error(ctx, "error creating role", zap.Error(err))

		return respondError("Error creating role, oops")
	}

	role := models.Role{
		GuildID: data.Event.GuildID.String(),
		RoleID:  r.ID.String(),
	}

	if err := role.Insert(ctx, tx, boil.Infer()); err != nil {
		if err := tx.Rollback(); err != nil {
			ctxlog.Error(ctx, "error rolling back tx", zap.Error(err))
		}

		if err := b.s.DeleteRole(data.Event.GuildID, r.ID, "rolling back"); err != nil {
			ctxlog.Error(ctx, "error deleting role from Discord", zap.Error(err))
		}

		return respondError("A database error has occurred, try again later or something")
	}

	if err := tx.Commit(); err != nil {
		ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		return respondError("A database error has occurred, try again later or something")
	}

	return respondf("Role `%s` has been created", r.Name)
}

func (b *Bot) cmdRoleDelete(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	sf, err := data.Options.Find("role").SnowflakeValue()
	if err != nil {
		return respondError("Unable to convert the given argument into a Snowflake xD")
	}

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error beginning tx", zap.Error(err))
		return respondError("A database error occurred :/")
	}

	role, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", data.Event.GuildID.String(), sf.String())).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondError("I'm not tracking a role by that ID, oop")
		}

		ctxlog.Error(ctx, "error querying for role", zap.Error(err))
		return dbError
	}

	if err := role.Delete(ctx, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			ctxlog.Error(ctx, "error rolling back transaction", zap.Error(err))
		}

		ctxlog.Error(ctx, "error deleting role", zap.Error(err))
		return dbError
	}

	if err := b.s.DeleteRole(data.Event.GuildID, discord.RoleID(sf), "Deletion requested"); err != nil {
		ctxlog.Error(ctx, "error deleting role", zap.Error(err))
		return respondError("The API got mad at me when I tried deleting the role :/")
	}

	if err := tx.Commit(); err != nil {
		ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		return dbError
	}

	return respond("Role has been deleted, woop woop")
}

func (b *Bot) cmdRoleRelinquish(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var sb strings.Builder
	sb.WriteString("Relinquished ")

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error beginning transaction", zap.Error(err))
		return dbError
	}

	for i, opt := range data.Options {
		sf, err := opt.SnowflakeValue()
		if err != nil {
			ctxlog.Error(ctx, "error converting to snowflake", zap.Error(err))
			return respondError("Unable to convert role to Snowflake")
		}

		if !sf.IsValid() {
			continue
		}

		role, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", data.Event.GuildID.String(), sf.String())).One(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			ctxlog.Error(ctx, "error querying for role", zap.Error(err))
			return dbError
		}

		if err := role.Delete(ctx, tx); err != nil {
			ctxlog.Error(ctx, "error deleting role", zap.Error(err))
			return respondErrorf("There was an issue deleting <@&%s> from the database", sf.String())
		}

		if i == len(data.Options)-1 {
			sb.WriteString("<@&" + sf.String() + ">")
			break
		}

		sb.WriteString("<@&" + sf.String() + ">, ")
	}

	if err := tx.Commit(); err != nil {
		if !errors.Is(err, sql.ErrTxDone) {
			ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
			return dbError
		}
	}

	return respond(sb.String())
}

func (b *Bot) cmdRoleAdd(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	sf, err := data.Options.Find("role").SnowflakeValue()
	if err != nil {
		return respondError("The role you gave me sucks")
	}

	exists, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", data.Event.GuildID.String(), sf.String())).Exists(ctx, b.db)
	if err != nil {
		ctxlog.Error(ctx, "error querying role", zap.Error(err))
		return dbError
	}

	if !exists {
		return respondError("I'm not tracking that role!!!")
	}

	if err := b.s.AddRole(data.Event.GuildID, data.Event.SenderID(), discord.RoleID(sf), api.AddRoleData{
		AuditLogReason: "Vanity role requested",
	}); err != nil {
		ctxlog.Error(ctx, "error adding role to user", zap.Error(err))
		return respondError("Something happened when I tried giving you that role")
	}

	return respondf("Role has been granted")
}

func (b *Bot) cmdRoleRemove(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	sf, err := data.Options.Find("role").SnowflakeValue()
	if err != nil {
		return respondError("The role you gave me sucks")
	}

	exists, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", data.Event.GuildID.String(), sf.String())).Exists(ctx, b.db)
	if err != nil {
		ctxlog.Error(ctx, "error querying role", zap.Error(err))
		return dbError
	}

	if !exists {
		return respondError("I'm not tracking that role in my database")
	}

	if err := b.s.RemoveRole(data.Event.GuildID, data.Event.SenderID(), discord.RoleID(sf), "Requested vanity role removal"); err != nil {
		ctxlog.Error(ctx, "error removing role from user", zap.Error(err))
		return respondError("Something happened when I tried removing that role from you")
	}

	return respondf("Role has been removed")
}

func (b *Bot) cmdRoleRename(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	sf, err := data.Options.Find("role").SnowflakeValue()
	if err != nil {
		return respondError("The role you gave me sucks")
	}

	newName := data.Options.Find("new-name").String()
	if newName == "" {
		return respondError("You gotta give me a new name to use!!")
	}

	exists, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", data.Event.GuildID.String(), sf.String())).Exists(ctx, b.db)
	if err != nil {
		ctxlog.Error(ctx, "error checking for role", zap.Error(err))
		return dbError
	}

	if !exists {
		return respondError("I'm not tracking a role by that ID!!")
	}

	if _, err := b.s.ModifyRole(data.Event.GuildID, discord.RoleID(sf), api.ModifyRoleData{
		Name: option.NewNullableString(newName),
	}); err != nil {
		ctxlog.Error(ctx, "error updating role", zap.Error(err))
		return respondError("The API got mad at me when I tried updating the role")
	}

	return respond("Role name has been updated")
}

func (b *Bot) cmdRoleSetColor(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	sf, err := data.Options.Find("role").SnowflakeValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing snowflake", zap.Error(err))
		return respondError("The role you provided is invalid")
	}

	role := discord.RoleID(sf)

	exists, err := models.Roles(qm.Where("guild_id = ? AND role_id = ?", data.Event.GuildID.String(), role)).Exists(ctx, b.db)
	if err != nil {
		ctxlog.Error(ctx, "error querying role", zap.Error(err))
		return dbError
	}

	if !exists {
		return respondError("I'm not managing that role!")
	}

	rawColor := data.Options.Find("color").String()
	rawColor = strings.TrimPrefix(rawColor, "#")

	i, err := strconv.ParseInt(rawColor, 16, 32)
	if err != nil {
		return respondError("The \"color\" you provided is invalid")
	}

	if _, err := b.s.ModifyRole(data.Event.GuildID, role, api.ModifyRoleData{
		Color: discord.Color(int32(i)),
	}); err != nil {
		ctxlog.Error(ctx, "error modifying role", zap.Error(err))
		return respondError("Error occurred when modifying the role's color")
	}

	return respond("Role color has been updated")
}

func (b *Bot) cmdRoleImport(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var sb strings.Builder
	sb.WriteString("Roles ")

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error beginning transaction", zap.Error(err))
		return dbError
	}

	for i, opt := range data.Options {
		sf, err := opt.SnowflakeValue()
		if err != nil {
			ctxlog.Error(ctx, "error converting option to snowflake")
			return respondError("Unable to convert option into Snowflake")
		}

		if !sf.IsValid() {
			continue
		}

		role := models.Role{
			RoleID:  sf.String(),
			GuildID: data.Event.GuildID.String(),
		}

		if err := role.Insert(ctx, tx, boil.Infer()); err != nil {
			ctxlog.Error(ctx, "error inserting role", zap.Error(err))
			return respondErrorf("There was an issue inserting <@&%s> into the database", sf.String())
		}

		if i == len(data.Options)-1 {
			sb.WriteString("<@&" + sf.String() + ">")
			break
		}

		sb.WriteString("<@&" + sf.String() + ">, ")
	}

	if err := tx.Commit(); err != nil {
		if !errors.Is(err, sql.ErrTxDone) {
			ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
			return dbError
		}
	}

	sb.WriteString(" have been added to the database")

	return respond(sb.String())
}

func (b *Bot) cmdRoleList(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	roles, err := models.Roles(qm.Where("guild_id = ?", data.Event.GuildID.String())).All(ctx, b.db)
	if err != nil {
		ctxlog.Error(ctx, "error querying roles", zap.Error(err))
		return dbError
	}

	var sb strings.Builder

	for i, r := range roles {
		sf, err := discord.ParseSnowflake(r.RoleID)
		if err != nil {
			ctxlog.Error(ctx, "error parsing role ID into snowflake", zap.Error(err))
			return respondError("OOPSIE WOOPSIE, unable to parse role ID into snowflake xD")
		}

		role, err := b.s.Role(data.Event.GuildID, discord.RoleID(sf))
		if err != nil {
			ctxlog.Error(ctx, "error retrieving role from cabinet", zap.Error(err))
			return respondError("Unable to retrieve role from cabinet :3")
		}

		if i == len(roles)-1 {
			sb.WriteString(role.Mention())
			continue
		}

		sb.WriteString(role.Mention() + ", ")
	}

	return respond(sb.String())
}
