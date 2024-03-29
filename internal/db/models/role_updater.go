// Code generated by SQLBoiler 4.16.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// RoleUpdater is an object representing the database table.
type RoleUpdater struct {
	ID            int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	GuildID       string    `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	DoUpdates     bool      `boil:"do_updates" json:"do_updates" toml:"do_updates" yaml:"do_updates"`
	LastTimestamp time.Time `boil:"last_timestamp" json:"last_timestamp" toml:"last_timestamp" yaml:"last_timestamp"`
	CreatedAt     time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *roleUpdaterR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L roleUpdaterL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var RoleUpdaterColumns = struct {
	ID            string
	GuildID       string
	DoUpdates     string
	LastTimestamp string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "id",
	GuildID:       "guild_id",
	DoUpdates:     "do_updates",
	LastTimestamp: "last_timestamp",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

var RoleUpdaterTableColumns = struct {
	ID            string
	GuildID       string
	DoUpdates     string
	LastTimestamp string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "role_updater.id",
	GuildID:       "role_updater.guild_id",
	DoUpdates:     "role_updater.do_updates",
	LastTimestamp: "role_updater.last_timestamp",
	CreatedAt:     "role_updater.created_at",
	UpdatedAt:     "role_updater.updated_at",
}

// Generated where

var RoleUpdaterWhere = struct {
	ID            whereHelperint64
	GuildID       whereHelperstring
	DoUpdates     whereHelperbool
	LastTimestamp whereHelpertime_Time
	CreatedAt     whereHelpertime_Time
	UpdatedAt     whereHelpertime_Time
}{
	ID:            whereHelperint64{field: "\"role_updater\".\"id\""},
	GuildID:       whereHelperstring{field: "\"role_updater\".\"guild_id\""},
	DoUpdates:     whereHelperbool{field: "\"role_updater\".\"do_updates\""},
	LastTimestamp: whereHelpertime_Time{field: "\"role_updater\".\"last_timestamp\""},
	CreatedAt:     whereHelpertime_Time{field: "\"role_updater\".\"created_at\""},
	UpdatedAt:     whereHelpertime_Time{field: "\"role_updater\".\"updated_at\""},
}

// RoleUpdaterRels is where relationship names are stored.
var RoleUpdaterRels = struct {
}{}

// roleUpdaterR is where relationships are stored.
type roleUpdaterR struct {
}

// NewStruct creates a new relationship struct
func (*roleUpdaterR) NewStruct() *roleUpdaterR {
	return &roleUpdaterR{}
}

// roleUpdaterL is where Load methods for each relationship are stored.
type roleUpdaterL struct{}

var (
	roleUpdaterAllColumns            = []string{"id", "guild_id", "do_updates", "last_timestamp", "created_at", "updated_at"}
	roleUpdaterColumnsWithoutDefault = []string{"guild_id"}
	roleUpdaterColumnsWithDefault    = []string{"id", "do_updates", "last_timestamp", "created_at", "updated_at"}
	roleUpdaterPrimaryKeyColumns     = []string{"id"}
	roleUpdaterGeneratedColumns      = []string{}
)

type (
	// RoleUpdaterSlice is an alias for a slice of pointers to RoleUpdater.
	// This should almost always be used instead of []RoleUpdater.
	RoleUpdaterSlice []*RoleUpdater

	roleUpdaterQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	roleUpdaterType                 = reflect.TypeOf(&RoleUpdater{})
	roleUpdaterMapping              = queries.MakeStructMapping(roleUpdaterType)
	roleUpdaterPrimaryKeyMapping, _ = queries.BindMapping(roleUpdaterType, roleUpdaterMapping, roleUpdaterPrimaryKeyColumns)
	roleUpdaterInsertCacheMut       sync.RWMutex
	roleUpdaterInsertCache          = make(map[string]insertCache)
	roleUpdaterUpdateCacheMut       sync.RWMutex
	roleUpdaterUpdateCache          = make(map[string]updateCache)
	roleUpdaterUpsertCacheMut       sync.RWMutex
	roleUpdaterUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single roleUpdater record from the query.
func (q roleUpdaterQuery) One(ctx context.Context, exec boil.ContextExecutor) (*RoleUpdater, error) {
	o := &RoleUpdater{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for role_updater")
	}

	return o, nil
}

// All returns all RoleUpdater records from the query.
func (q roleUpdaterQuery) All(ctx context.Context, exec boil.ContextExecutor) (RoleUpdaterSlice, error) {
	var o []*RoleUpdater

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to RoleUpdater slice")
	}

	return o, nil
}

// Count returns the count of all RoleUpdater records in the query.
func (q roleUpdaterQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count role_updater rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q roleUpdaterQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if role_updater exists")
	}

	return count > 0, nil
}

// RoleUpdaters retrieves all the records using an executor.
func RoleUpdaters(mods ...qm.QueryMod) roleUpdaterQuery {
	mods = append(mods, qm.From("\"role_updater\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"role_updater\".*"})
	}

	return roleUpdaterQuery{q}
}

// FindRoleUpdater retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRoleUpdater(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*RoleUpdater, error) {
	roleUpdaterObj := &RoleUpdater{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"role_updater\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, roleUpdaterObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from role_updater")
	}

	return roleUpdaterObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *RoleUpdater) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no role_updater provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(roleUpdaterColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	roleUpdaterInsertCacheMut.RLock()
	cache, cached := roleUpdaterInsertCache[key]
	roleUpdaterInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			roleUpdaterAllColumns,
			roleUpdaterColumnsWithDefault,
			roleUpdaterColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(roleUpdaterType, roleUpdaterMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(roleUpdaterType, roleUpdaterMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"role_updater\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"role_updater\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into role_updater")
	}

	if !cached {
		roleUpdaterInsertCacheMut.Lock()
		roleUpdaterInsertCache[key] = cache
		roleUpdaterInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the RoleUpdater.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *RoleUpdater) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	key := makeCacheKey(columns, nil)
	roleUpdaterUpdateCacheMut.RLock()
	cache, cached := roleUpdaterUpdateCache[key]
	roleUpdaterUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			roleUpdaterAllColumns,
			roleUpdaterPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update role_updater, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"role_updater\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, roleUpdaterPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(roleUpdaterType, roleUpdaterMapping, append(wl, roleUpdaterPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update role_updater row")
	}

	if !cached {
		roleUpdaterUpdateCacheMut.Lock()
		roleUpdaterUpdateCache[key] = cache
		roleUpdaterUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q roleUpdaterQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for role_updater")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RoleUpdaterSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roleUpdaterPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"role_updater\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, roleUpdaterPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in roleUpdater slice")
	}

	return nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *RoleUpdater) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no role_updater provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	nzDefaults := queries.NonZeroDefaultSet(roleUpdaterColumnsWithDefault, o)

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

	roleUpdaterUpsertCacheMut.RLock()
	cache, cached := roleUpdaterUpsertCache[key]
	roleUpdaterUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			roleUpdaterAllColumns,
			roleUpdaterColumnsWithDefault,
			roleUpdaterColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			roleUpdaterAllColumns,
			roleUpdaterPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert role_updater, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(roleUpdaterPrimaryKeyColumns))
			copy(conflict, roleUpdaterPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"role_updater\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(roleUpdaterType, roleUpdaterMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(roleUpdaterType, roleUpdaterMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert role_updater")
	}

	if !cached {
		roleUpdaterUpsertCacheMut.Lock()
		roleUpdaterUpsertCache[key] = cache
		roleUpdaterUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single RoleUpdater record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *RoleUpdater) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no RoleUpdater provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), roleUpdaterPrimaryKeyMapping)
	sql := "DELETE FROM \"role_updater\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from role_updater")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q roleUpdaterQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no roleUpdaterQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from role_updater")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RoleUpdaterSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roleUpdaterPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"role_updater\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, roleUpdaterPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from roleUpdater slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *RoleUpdater) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindRoleUpdater(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RoleUpdaterSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := RoleUpdaterSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), roleUpdaterPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"role_updater\".* FROM \"role_updater\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, roleUpdaterPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RoleUpdaterSlice")
	}

	*o = slice

	return nil
}

// RoleUpdaterExists checks if the RoleUpdater row exists.
func RoleUpdaterExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"role_updater\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if role_updater exists")
	}

	return exists, nil
}

// Exists checks if the RoleUpdater row exists.
func (o *RoleUpdater) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return RoleUpdaterExists(ctx, exec, o.ID)
}
