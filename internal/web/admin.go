package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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

type quoteTime struct {
	time.Time
}

func (t *quoteTime) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	tme, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = tme
	return nil
}

func (s *Server) adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(s.admins) == 0 {
			s.notAuthorized(w, r, false)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			s.notAuthorized(w, r, true)
			return
		}

		expected := s.admins[user]
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
	r.Get("/quotes/update", s.adminUpdateQuotes)
	r.Post("/quotes/update", s.postAdminUpdateQuotes)
	r.Get("/quotes/delete", s.adminDeleteQuotes)
	r.Post("/quotes/delete", s.postAdminDeleteQuotes)
}

func (s *Server) adminImportQuotes(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.AdminQuotesPage{
		BasePage: s.basePage(r),
	})
}

func (s *Server) adminUpdateQuotes(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.AdminQuotesUpdatePage{
		BasePage: s.basePage(r),
	})
}

func (s *Server) adminDeleteQuotes(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.AdminQuotesDeletePage{
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
	Quote          string    `json:"quote"`
	QuoterID       string    `json:"quoter_id"`
	QuotedID       string    `json:"quoted_id"`
	QuotedUsername string    `json:"quoted_username"`
	GuildID        string    `json:"guild_id"`
	ChannelID      string    `json:"channel_id"`
	MessageID      string    `json:"message_id"`
	CreatedAt      quoteTime `json:"created_at"`
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

	guild, err := modelsx.FetchGuild(ctx, s.db, f.GuildID)
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

	if err := modelsx.ToggleGuildQuotes(ctx, s.db, guild, false); err != nil {
		ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
		respondf(w, "error toggling quotes off on guild")
		return
	}

	defer func() {
		if reenable {
			if err := modelsx.ToggleGuildQuotes(ctx, s.db, guild, true); err != nil {
				ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
				respondf(w, "error toggling quotes off on guild")
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
	).Bind(ctx, s.db, &row)
	if err != nil {
		ctxlog.Error(ctx, "error getting latest quote number from database", zap.Error(err))
		respondf(w, "error getting latest quote number from database")
		return
	}

	nextNum := row.MaxNum.Int + 1

	tx, err := s.db.BeginTx(ctx, nil)
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
			CreatedAt:      q.CreatedAt.Time,
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

type updatedQuote struct {
	MessageID string `json:"message_id"` // This cannot be updated

	CreatedAt quoteTime `json:"created_at"`
}

type updatedQuotesFile struct {
	GuildID string          `json:"guild_id"`
	Quotes  []*updatedQuote `json:"quotes"`
}

func (s *Server) postAdminUpdateQuotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, _, err := r.FormFile("quotes")
	if err != nil {
		ctxlog.Error(ctx, "error retrieving form file", zap.Error(err))
		respondf(w, "error retrieving form file")
		return
	}

	defer file.Close()

	var f *updatedQuotesFile
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

	guild, err := modelsx.FetchGuild(ctx, s.db, f.GuildID)
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

	if err := modelsx.ToggleGuildQuotes(ctx, s.db, guild, false); err != nil {
		ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
		respondf(w, "error toggling quotes off on guild")
		return
	}

	defer func() {
		if reenable {
			if err := modelsx.ToggleGuildQuotes(ctx, s.db, guild, true); err != nil {
				ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
				respondf(w, "error toggling quotes off on guild")
				return
			}
		}
	}()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error starting transaction", zap.Error(err))
		respondf(w, "error starting transaction")
		return
	}

	defer tx.Rollback()

	quotes, err := models.Quotes(qm.Where("guild_id = ?", f.GuildID)).All(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondf(w, "guild has no quotes")
			return
		}

		ctxlog.Error(ctx, "error fetching quotes", zap.Error(err))
		respondf(w, "error fetching quotes")
		return
	}

	for _, q := range quotes {
		for _, fq := range f.Quotes {
			if strings.EqualFold(q.MessageID, fq.MessageID) {
				q.CreatedAt = fq.CreatedAt.Time

				if err := q.Update(ctx, tx, boil.Whitelist(
					models.QuoteColumns.CreatedAt,
					models.QuoteColumns.UpdatedAt,
				)); err != nil {
					ctxlog.Error(ctx, "error updating quote", zap.Error(err))
					respondf(w, "error updating quote: %s", fq.MessageID)
					return
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		respondf(w, "error committing transaction", zap.Error(err))
	}
}

type deletedQuotesFile struct {
	GuildID string `json:"guild_id"`
	Quotes  []int  `json:"quotes"`
}

func (s *Server) postAdminDeleteQuotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, _, err := r.FormFile("quotes")
	if err != nil {
		ctxlog.Error(ctx, "error retrieving form file", zap.Error(err))
		respondf(w, "error retrieving form file")
		return
	}

	defer file.Close()

	var f *deletedQuotesFile
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

	guild, err := modelsx.FetchGuild(ctx, s.db, f.GuildID)
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

	if err := modelsx.ToggleGuildQuotes(ctx, s.db, guild, false); err != nil {
		ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
		respondf(w, "error toggling quotes off on guild")
		return
	}

	defer func() {
		if reenable {
			if err := modelsx.ToggleGuildQuotes(ctx, s.db, guild, true); err != nil {
				ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
				respondf(w, "error toggling quotes off on guild")
				return
			}
		}
	}()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		ctxlog.Error(ctx, "error starting transaction", zap.Error(err))
		respondf(w, "error starting transaction")
		return
	}

	defer tx.Rollback()

	quotes, err := models.Quotes(qm.Where("guild_id = ?", f.GuildID)).All(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondf(w, "guild has no quotes")
			return
		}

		ctxlog.Error(ctx, "error fetching quotes", zap.Error(err))
		respondf(w, "error fetching quotes")
		return
	}

	for _, q := range quotes {
		for _, fq := range f.Quotes {
			if q.Num == fq {
				if err := q.Delete(ctx, tx); err != nil {
					ctxlog.Error(ctx, "error deleting quote", zap.Error(err))
					respondf(w, "error deleting quote %d", fq)
					return
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		ctxlog.Error(ctx, "error committing transaction", zap.Error(err))
		respondf(w, "error committing transaction", zap.Error(err))
	}
}
