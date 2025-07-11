package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// driverFlavor distinguishes parameter placeholder style between databases.
// e.g. Postgres uses $1, $2 ... while MySQL uses ?
// The Prepare() phase will map generic "?" placeholders to the correct style.
type driverFlavor int

const (
	flavorMySQL driverFlavor = iota
	flavorPostgres
)

// SqlQueryAdapter is a lightweight ORM-like query builder that satisfies QueryAdapter.
// It uses standard database/sql underneath and supports both MySQL and PostgreSQL
// by swapping placeholder styles at execution time.
// All building state is immutable — chain calls clone the adapter to guarantee
// thread-safety.
type SqlQueryAdapter struct {
	db     *sql.DB
	flavor driverFlavor
	ctx    context.Context

	// query parts
	table     string
	fields    []string
	wheres    []string
	whereArgs []any
	orWheres  []string
	orArgs    []any
	orderBy   string
	limit     *int
	offset    *int
}

// NewSqlAdapter creates a new adapter based on *sql.DB.
// The flavor is detected from the driver name, but can be overridden.
func NewSqlAdapter(db *sql.DB) *SqlQueryAdapter {
	return &SqlQueryAdapter{
		db:       db,
		flavor:   detectFlavor(db),
		ctx:      context.Background(),
		fields:   []string{"*"},
		wheres:   []string{},
		orWheres: []string{},
	}
}

// clone makes a deep copy so chain calls do not mutate the original.
func (q *SqlQueryAdapter) clone() *SqlQueryAdapter {
	cp := *q
	// copy slices to maintain immutability
	cp.fields = append([]string(nil), q.fields...)
	cp.wheres = append([]string(nil), q.wheres...)
	cp.whereArgs = append([]any(nil), q.whereArgs...)
	cp.orWheres = append([]string(nil), q.orWheres...)
	cp.orArgs = append([]any(nil), q.orArgs...)
	return &cp
}

// --- QueryAdapter interface implementation ---

func (q *SqlQueryAdapter) WithContext(ctx context.Context) QueryAdapter {
	cp := q.clone()
	cp.ctx = ctx
	return cp
}

func (q *SqlQueryAdapter) UseModel(m Tabler) QueryAdapter {
	cp := q.clone()
	cp.table = m.TableName()
	return cp
}

func (q *SqlQueryAdapter) Model() Tabler {
	// Empty implementation — caller generally already has the model
	return nil
}

func (q *SqlQueryAdapter) Select(s []string) QueryAdapter {
	cp := q.clone()
	if len(s) > 0 {
		cp.fields = s
	}
	return cp
}

func (q *SqlQueryAdapter) Where(cond any, args ...any) QueryAdapter {
	cp := q.clone()
	cp.wheres = append(cp.wheres, toString(cond))
	cp.whereArgs = append(cp.whereArgs, args...)
	return cp
}

func (q *SqlQueryAdapter) Or(cond any, args ...any) QueryAdapter {
	cp := q.clone()
	cp.orWheres = append(cp.orWheres, toString(cond))
	cp.orArgs = append(cp.orArgs, args...)
	return cp
}

func (q *SqlQueryAdapter) Order(order string) QueryAdapter {
	cp := q.clone()
	cp.orderBy = order
	return cp
}

func (q *SqlQueryAdapter) Limit(l int) QueryAdapter {
	cp := q.clone()
	cp.limit = &l
	return cp
}

func (q *SqlQueryAdapter) Offset(o int) QueryAdapter {
	cp := q.clone()
	cp.offset = &o
	return cp
}

func (q *SqlQueryAdapter) Count(target *int64) error {
	query, args := q.build(true)
	row := q.db.QueryRowContext(q.ctx, query, args...)
	return row.Scan(target)
}

func (q *SqlQueryAdapter) Scan(dest any) error {
	query, args := q.build(false)
	rows, err := q.db.QueryContext(q.ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	// sql.Rows can scan into structs with sqlx, but here we keep it simple.
	// Caller is responsible for rows.Scan in a custom manner.
	switch d := dest.(type) {
	case *[]map[string]any:
		cols, _ := rows.Columns()
		for rows.Next() {
			vals := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				return err
			}
			m := map[string]any{}
			for i, c := range cols {
				m[c] = vals[i]
			}
			*d = append(*d, m)
		}
		return rows.Err()
	default:
		return fmt.Errorf("Scan: unsupported destination type %T", dest)
	}
}

func (q *SqlQueryAdapter) Clone() QueryAdapter {
	return q.clone()
}

// --- helpers ---

// detectFlavor inspects the database driver to determine if it is
func detectFlavor(db *sql.DB) driverFlavor {
	t := reflect.TypeOf(db.Driver()).String() // includes leading '*'
	t = strings.TrimPrefix(t, "*")
	switch {
	case strings.Contains(t, "pq"), strings.Contains(t, "pgx"), strings.Contains(t, "postgres"):
		return flavorPostgres
	default:
		return flavorMySQL
	}
}

func (q *SqlQueryAdapter) build(countMode bool) (string, []any) {
	var sb strings.Builder

	if countMode {
		sb.WriteString("SELECT COUNT(1) FROM ")
	} else {
		sb.WriteString("SELECT ")
		sb.WriteString(strings.Join(q.fields, ", "))
		sb.WriteString(" FROM ")
	}
	sb.WriteString(q.table)

	args := make([]any, 0, len(q.whereArgs)+len(q.orArgs))

	if len(q.wheres) > 0 || len(q.orWheres) > 0 {
		sb.WriteString(" WHERE ")
		if len(q.wheres) > 0 {
			sb.WriteString(strings.Join(q.wheres, " AND "))
			args = append(args, q.whereArgs...)
		}
		if len(q.orWheres) > 0 {
			if len(q.wheres) > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString("(")
			sb.WriteString(strings.Join(q.orWheres, " OR "))
			sb.WriteString(")")
			args = append(args, q.orArgs...)
		}
	}

	if q.orderBy != "" && !countMode {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(q.orderBy)
	}

	if q.limit != nil && !countMode {
		sb.WriteString(" LIMIT ")
		sb.WriteString(fmt.Sprint(*q.limit))
	}

	if q.offset != nil && !countMode {
		sb.WriteString(" OFFSET ")
		sb.WriteString(fmt.Sprint(*q.offset))
	}

	// translate placeholders if needed
	query := sb.String()
	if q.flavor == flavorPostgres {
		// replace ? with $n
		var idx int
		var b strings.Builder
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				idx++
				b.WriteString("$")
				b.WriteString(fmt.Sprint(idx))
			} else {
				b.WriteByte(query[i])
			}
		}
		query = b.String()
	}

	return query, args
}

func toString(cond any) string {
	switch v := cond.(type) {
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}
