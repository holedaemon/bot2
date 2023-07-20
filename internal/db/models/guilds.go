// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Guild is an object representing the database table.
type Guild struct {
	ID        int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	GuildID   string    `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	DoQuotes  bool      `boil:"do_quotes" json:"do_quotes" toml:"do_quotes" yaml:"do_quotes"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *guildR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L guildL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GuildColumns = struct {
	ID        string
	GuildID   string
	DoQuotes  string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	GuildID:   "guild_id",
	DoQuotes:  "do_quotes",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var GuildTableColumns = struct {
	ID        string
	GuildID   string
	DoQuotes  string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "guilds.id",
	GuildID:   "guilds.guild_id",
	DoQuotes:  "guilds.do_quotes",
	CreatedAt: "guilds.created_at",
	UpdatedAt: "guilds.updated_at",
}

// Generated where

var GuildWhere = struct {
	ID        whereHelperint64
	GuildID   whereHelperstring
	DoQuotes  whereHelperbool
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperint64{field: "\"guilds\".\"id\""},
	GuildID:   whereHelperstring{field: "\"guilds\".\"guild_id\""},
	DoQuotes:  whereHelperbool{field: "\"guilds\".\"do_quotes\""},
	CreatedAt: whereHelpertime_Time{field: "\"guilds\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"guilds\".\"updated_at\""},
}

// GuildRels is where relationship names are stored.
var GuildRels = struct {
}{}

// guildR is where relationships are stored.
type guildR struct {
}

// NewStruct creates a new relationship struct
func (*guildR) NewStruct() *guildR {
	return &guildR{}
}

// guildL is where Load methods for each relationship are stored.
type guildL struct{}

var (
	guildAllColumns            = []string{"id", "guild_id", "do_quotes", "created_at", "updated_at"}
	guildColumnsWithoutDefault = []string{"guild_id"}
	guildColumnsWithDefault    = []string{"id", "do_quotes", "created_at", "updated_at"}
	guildPrimaryKeyColumns     = []string{"id"}
	guildGeneratedColumns      = []string{}
)

type (
	// GuildSlice is an alias for a slice of pointers to Guild.
	// This should almost always be used instead of []Guild.
	GuildSlice []*Guild

	guildQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	guildType                 = reflect.TypeOf(&Guild{})
	guildMapping              = queries.MakeStructMapping(guildType)
	guildPrimaryKeyMapping, _ = queries.BindMapping(guildType, guildMapping, guildPrimaryKeyColumns)
	guildInsertCacheMut       sync.RWMutex
	guildInsertCache          = make(map[string]insertCache)
	guildUpdateCacheMut       sync.RWMutex
	guildUpdateCache          = make(map[string]updateCache)
	guildUpsertCacheMut       sync.RWMutex
	guildUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single guild record from the query.
func (q guildQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Guild, error) {
	o := &Guild{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for guilds")
	}

	return o, nil
}

// All returns all Guild records from the query.
func (q guildQuery) All(ctx context.Context, exec boil.ContextExecutor) (GuildSlice, error) {
	var o []*Guild

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Guild slice")
	}

	return o, nil
}

// Count returns the count of all Guild records in the query.
func (q guildQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count guilds rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q guildQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if guilds exists")
	}

	return count > 0, nil
}

// Guilds retrieves all the records using an executor.
func Guilds(mods ...qm.QueryMod) guildQuery {
	mods = append(mods, qm.From("\"guilds\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"guilds\".*"})
	}

	return guildQuery{q}
}

// FindGuild retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGuild(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Guild, error) {
	guildObj := &Guild{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"guilds\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, guildObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from guilds")
	}

	return guildObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Guild) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no guilds provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	nzDefaults := queries.NonZeroDefaultSet(guildColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	guildInsertCacheMut.RLock()
	cache, cached := guildInsertCache[key]
	guildInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			guildAllColumns,
			guildColumnsWithDefault,
			guildColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(guildType, guildMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(guildType, guildMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"guilds\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"guilds\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into guilds")
	}

	if !cached {
		guildInsertCacheMut.Lock()
		guildInsertCache[key] = cache
		guildInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Guild.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Guild) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	key := makeCacheKey(columns, nil)
	guildUpdateCacheMut.RLock()
	cache, cached := guildUpdateCache[key]
	guildUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			guildAllColumns,
			guildPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update guilds, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"guilds\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, guildPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(guildType, guildMapping, append(wl, guildPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	_, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update guilds row")
	}

	if !cached {
		guildUpdateCacheMut.Lock()
		guildUpdateCache[key] = cache
		guildUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q guildQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for guilds")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GuildSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), guildPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"guilds\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, guildPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in guild slice")
	}

	return nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Guild) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no guilds provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(guildColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	guildUpsertCacheMut.RLock()
	cache, cached := guildUpsertCache[key]
	guildUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			guildAllColumns,
			guildColumnsWithDefault,
			guildColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			guildAllColumns,
			guildPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert guilds, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(guildPrimaryKeyColumns))
			copy(conflict, guildPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"guilds\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(guildType, guildMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(guildType, guildMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert guilds")
	}

	if !cached {
		guildUpsertCacheMut.Lock()
		guildUpsertCache[key] = cache
		guildUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Guild record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Guild) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no Guild provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), guildPrimaryKeyMapping)
	sql := "DELETE FROM \"guilds\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from guilds")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q guildQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no guildQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from guilds")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GuildSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), guildPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"guilds\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, guildPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from guild slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Guild) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGuild(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GuildSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GuildSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), guildPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"guilds\".* FROM \"guilds\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, guildPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GuildSlice")
	}

	*o = slice

	return nil
}

// GuildExists checks if the Guild row exists.
func GuildExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"guilds\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if guilds exists")
	}

	return exists, nil
}

// Exists checks if the Guild row exists.
func (o *Guild) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return GuildExists(ctx, exec, o.ID)
}
