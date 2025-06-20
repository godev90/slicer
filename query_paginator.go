package paginator

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/godev90/validator/errors"

	"gorm.io/gorm"
)

type (
	Tabler interface {
		TableName() string
	}

	Paginator[T Tabler] interface {
		AllowedFields() map[string]string
		DB() *gorm.DB
		UseDB(*gorm.DB)
		Model() T
		Items() []T
		SetItems([]T)
	}
)

func DefaultTablerAllowedFields(model Tabler) map[string]string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonName := f.Tag.Get("json")
		if jsonName == "" || jsonName == "-" {
			jsonName = strings.ToLower(f.Name)
		}
		gormTag := f.Tag.Get("gorm")
		var columnName string
		if gormTag != "" {
			for _, part := range strings.Split(gormTag, ";") {
				if strings.HasPrefix(part, "column:") {
					columnName = strings.TrimPrefix(part, "column:")
					break
				}
			}
		}
		if columnName != "" && regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(columnName) {
			fields[jsonName] = columnName
		}
	}
	return fields
}

func Page[T Tabler](paginator Paginator[T], opts QueryOptions) (PageData, error) {
	var (
		model   = paginator.Model()
		db      = paginator.DB().Model(model)
		allowed = paginator.AllowedFields()

		modelType = reflect.TypeOf(model)
	)

	for key, val := range opts.Filters {
		if col, ok := allowed[key]; ok {
			db = db.Where(fmt.Sprintf("%s = ?", col), val)
		}
	}

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for _, cmp := range opts.Comparisons {
		if col, ok := allowed[cmp.Field]; ok {
			parsed := cmp.Value

			var foundField reflect.StructField
			found := false

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
		first := true
		keyword := strings.ToLower(opts.Search.Keyword)
		for _, field := range opts.Search.Fields {
			if col, ok := allowed[field]; ok {
				cond := fmt.Sprintf("%s LIKE ?", col)
				if first {
					db = db.Where(cond, "%"+keyword+"%")
					first = false
				} else {
					db = db.Or(cond, "%"+keyword+"%")
				}
			}
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
	countErr := db.WithContext(ctx).Session(&gorm.Session{}).Count(&total).Error
	if countErr != nil {
		if errors.Is(countErr, context.DeadlineExceeded) {
			total = -1
		} else {
			return PageData{
				Items: []string{},
				Total: 0,
				Page:  opts.Page,
				Limit: opts.Limit,
				LastError: errors.New(countErr, &errors.ErrAttr{
					Code: http.StatusInternalServerError,
				}),
			}, countErr
		}
	}

	items := paginator.Items()
	db = db.Offset(opts.Offset).Limit(opts.Limit)
	if err := db.Find(&items).Error; err != nil {
		return PageData{Items: paginator.Items(),
			Total: total,
			Page:  opts.Page,
			Limit: opts.Limit,
			LastError: errors.New(err, &errors.ErrAttr{
				Code: http.StatusInternalServerError,
			})}, err
	}

	paginator.SetItems(items)

	return PageData{Items: paginator.Items(), Total: total, Page: opts.Page, Limit: opts.Limit}, nil
}
