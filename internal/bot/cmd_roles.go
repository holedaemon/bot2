package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

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

	exists, err := models.Roles(qm.Where("guild_id = ? AND role_name ILIKE ?", data.Event.GuildID.String(), name)).Exists(ctx, tx)
	if err != nil {
		ctxlog.Error(ctx, "error checking for role existence", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			ctxlog.Error(ctx, "error rolling back", zap.Error(err))
		}

		return respondError("A database issue has occurred, watch this")
	}

	if exists {
		return respondError("A role by that name already exists in this server!!")
	}

	r, err := b.s.CreateRole(data.Event.GuildID, api.CreateRoleData{
		Name:        name,
		Mentionable: true,
	})
	if err != nil {
		ctxlog.Error(ctx, "error creating role", zap.Error(err))

		return respondError("Error creating role, oops")
	}

	role := models.Role{
		RoleName: name,
		GuildID:  data.Event.GuildID.String(),
		RoleID:   r.ID.String(),
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

	return respondf("Role `%s` has been created", r.Name)
}

// func (b *Bot) cmdRoleAdd(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
// 	name := data.Options.Find("name").String()
// 	if name == "" {
// 		return respondError("I can't apply a role to you if you don't give me a name")
// 	}

// 	role, err := models.Roles(qm.Where("guild_id = ? AND role_name ILIKE ?", data.Event.GuildID.String(), name)).One(ctx, b.db)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return respondError("That role isn't tracked in my database, I don't feel comfortable giving it you")
// 		}

// 		ctxlog.Error(ctx,"error querying role", zap.Error(err))
// 	}
// }
