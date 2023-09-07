package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/holedaemon/bot2/internal/web/templates"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (s *Server) adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(s.Admins) == 0 {
			s.notAuthorized(w, r, false)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			s.notAuthorized(w, r, true)
			return
		}

		expected := s.Admins[user]
		if expected == "" || pass != expected {
			s.notAuthorized(w, r, true)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) routeAdmin(r chi.Router) {
	r.Use(s.adminAuth)

	r.Get("/quotes", s.adminImportQuotes)
	r.Post("/quotes", s.postAdminImportQuotes)
}

func (s *Server) adminImportQuotes(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.AdminQuotesPage{
		BasePage: s.basePage(r),
	})
}

func respondf(w http.ResponseWriter, msg string, args ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}

	fmt.Fprintf(w, msg, args...)
}

type importedQuote struct {
	Quote          string `json:"quote"`
	QuoterID       string `json:"quoter_id"`
	QuotedID       string `json:"quoted_id"`
	QuotedUsername string `json:"quoted_username"`
	GuildID        string `json:"guild_id"`
	ChannelID      string `json:"channel_id"`
	MessageID      string `json:"message_id"`
}

type importedQuotesFile struct {
	GuildID string           `json:"guild_id"`
	Quotes  []*importedQuote `json:"quotes"`
}

func (s *Server) postAdminImportQuotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, _, err := r.FormFile("quotes")
	if err != nil {
		ctxlog.Error(ctx, "error retrieving form file", zap.Error(err))
		respondf(w, "error retrieving form file")
		return
	}

	defer file.Close()

	var f *importedQuotesFile
	if err := json.NewDecoder(file).Decode(&f); err != nil {
		ctxlog.Error(ctx, "error unmarshalling json", zap.Error(err))
		respondf(w, "error unmarshalling json")
		return
	}

	if f.GuildID == "" {
		respondf(w, "You didn't set a guild_id, doofus!!")
		return
	}

	if len(f.Quotes) == 0 {
		respondf(w, "You didn't upload any quotes, doofus!!")
		return
	}

	guild, err := modelsx.FetchGuild(ctx, s.DB, f.GuildID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondf(w, "That guild isn't in the database!!! ")
			return
		}

		ctxlog.Error(ctx, "error toggling quotes off on guild", zap.Error(err))
		respondf(w, "error toggling quotes off on guild")
		return
	}

	var reenable = false
	if guild.DoQuotes {
		reenable = true
	}

	guild.DoQuotes = false

	if err := guild.Update(ctx, s.DB, boil.Infer()); err != nil {
		ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
		respondf(w, "error toggling quotes off on guild")
		return
	}

	defer func() {
		if reenable {
			guild.DoQuotes = true

			if err := guild.Update(ctx, s.DB, boil.Infer()); err != nil {
				ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
				respondf(w, "error toggling quotes on on guild")
				return
			}
		}
	}()

	var row struct {
		MaxNum null.Int
	}

	err = models.Quotes(
		qm.Where("guild_id = ?", f.GuildID),
		qm.Select("max("+models.QuoteColumns.Num+") as max_num"),
	).Bind(ctx, s.DB, &row)
	if err != nil {
		ctxlog.Error(ctx, "error getting latest quote number from database", zap.Error(err))
		respondf(w, "error getting latest quote number from database")
		return
	}

	nextNum := row.MaxNum.Int + 1

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error starting transaction", zap.Error(err))
		respondf(w, "error starting transaction")
		return
	}

	defer tx.Rollback()

	for _, q := range f.Quotes {
		quote := &models.Quote{
			Quote:          q.Quote,
			Num:            nextNum,
			QuoterID:       null.StringFrom(q.QuoterID),
			QuotedID:       q.QuotedID,
			QuotedUsername: q.QuotedUsername,
			GuildID:        q.GuildID,
			ChannelID:      q.ChannelID,
			MessageID:      q.MessageID,
		}

		if err := quote.Insert(ctx, tx, boil.Infer()); err != nil {
			ctxlog.Error(ctx, "error inserting quote", zap.String("quote_message_id", q.MessageID), zap.Error(err))
			respondf(w, "error inserting quote %s", q.MessageID)
			return
		}
		nextNum++
	}

	if err := tx.Commit(); err != nil {
		ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		respondf(w, "error committing transaction", zap.Error(err))
	}
}
