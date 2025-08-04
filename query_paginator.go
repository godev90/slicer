package slicer

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/godev90/orm"
	"github.com/godev90/validator/faults"
)

type (
	Paginator[T orm.Tabler] interface {
		AllowedFields() map[string]string
		Adapter() orm.QueryAdapter
		Model() T
		Items() []T
		SetItems([]T)
	}
)

func QueryPage[T orm.Tabler](paginator Paginator[T], opts QueryOptions) (PageData, error) {
	var (
		model   = paginator.Model()
		db      = paginator.Adapter().UseModel(model)
		allowed = paginator.AllowedFields()

		modelType = reflect.TypeOf(model)
	)

	for key, val := range opts.Filters {
		if col, ok := allowed[key]; ok {
			parts := strings.Split(val, ",")
			if len(parts) == 1 {
				db = db.Where(fmt.Sprintf("%s = ?", col), val)
			} else {
				args := make([]any, len(parts))
				for i, v := range parts {
					args[i] = v
				}
				placeholders := strings.Repeat("?,", len(parts))
				placeholders = strings.TrimRight(placeholders, ",")
				db = db.Where(fmt.Sprintf("%s IN (%s)", col, placeholders), args...)
			}
		}
	}

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for _, cmp := range opts.Comparisons {
		if col, ok := allowed[cmp.Field]; ok {
			var (
				parsed     = cmp.Value
				found      = false
				foundField reflect.StructField
			)

			for i := 0; i < modelType.NumField(); i++ {
				field := modelType.Field(i)
				jsonTag := field.Tag.Get("json")
				if jsonTag == cmp.Field {
					foundField = field
					found = true
					break
				}
				if jsonTag == "" && strings.EqualFold(field.Name, cmp.Field) {
					foundField = field
					found = true
					break
				}
			}

			if found && (foundField.Type.String() == "types.Date" || foundField.Type.String() == "types.Datetime") {
				if t, err := time.Parse("2006-01-02", cmp.Value); err == nil {
					switch cmp.Op {
					case GT:
						t = t.Add(24 * time.Hour)
						parsed = t.Format("2006-01-02 15:04:05")
					case LTE:
						t = t.Add(24 * time.Hour).Add(-time.Nanosecond)
						parsed = t.Format("2006-01-02 15:04:05")
					default:
						parsed = t.Format("2006-01-02 15:04:05")
					}
				}
			}

			symbol := map[ComparisonOp]string{
				GT:  ">",
				GTE: ">=",
				LT:  "<",
				LTE: "<=",
				EQ:  "=",
			}[cmp.Op]

			db = db.Where(fmt.Sprintf("%s %s ?", col, symbol), parsed)
		}
	}

	if opts.Search != nil {
		var (
			useKeyword = false
			keyword    = strings.ToLower(opts.Search.Keyword)
			clone      = db.Clone()
		)

		for _, field := range opts.Search.Fields {
			if col, ok := allowed[field]; ok {
				cond := fmt.Sprintf("%s LIKE ?", col)
				clone = clone.Or(cond, "%"+keyword+"%")
				useKeyword = true
			}
		}

		if useKeyword {
			db = db.Where(clone)
		}
	}

	if len(opts.Select) > 0 {
		columns := []string{}
		for _, field := range opts.Select {
			if col, ok := allowed[field]; ok {
				columns = append(columns, col)
			}
		}
		if len(columns) > 0 {
			db = db.Select(columns)
		}
	}

	for _, s := range opts.Sort {
		if col, ok := allowed[s.Field]; ok {
			if s.Desc {
				db = db.Order(col + " DESC")
			} else {
				db = db.Order(col + " ASC")
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var total int64
	countErr := db.WithContext(ctx).Count(&total)
	if countErr != nil {
		if faults.Is(countErr, context.DeadlineExceeded) {
			total = -1
		} else {
			return PageData{
				Items: []string{},
				Total: 0,
				Page:  opts.Page,
				Limit: opts.Limit,
				LastError: faults.New(countErr, &faults.ErrAttr{
					Code: http.StatusInternalServerError,
				}),
			}, countErr
		}
	}

	items := paginator.Items()
	db = db.Offset(opts.Offset).Limit(opts.Limit)
	if err := db.Scan(&items); err != nil {
		return PageData{Items: paginator.Items(),
			Total: total,
			Page:  opts.Page,
			Limit: opts.Limit,
			LastError: faults.New(err, &faults.ErrAttr{
				Code: http.StatusInternalServerError,
			})}, err
	}

	if items == nil {
		items = []T{}
	}

	paginator.SetItems(items)

	return PageData{Items: paginator.Items(), Total: total, Page: opts.Page, Limit: opts.Limit}, nil
}
