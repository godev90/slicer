package slicer

import (
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/godev90/validator/typedef"
)

type (
	QueryOptions struct {
		Page        int
		Limit       int
		Offset      int
		Sort        []SortField
		Search      *SearchQuery
		SearchAnd   *SearchQueryAnd
		Filters     map[string]string
		Select      []string
		GroupBy     []string
		Comparisons []ComparisonFilter
	}

	SortField struct {
		Field string
		Desc  bool
	}

	SearchQuery struct {
		Fields  []string
		Keyword string
	}

	SearchField struct {
		Field   string
		Keyword string
	}

	SearchQueryAnd struct {
		Fields []*SearchField
	}

	PageData struct {
		LastError error `json:"error,omitempty"`
		Items     any   `json:"items"`
		Total     int64 `json:"total"`
		Page      int   `json:"page"`
		Limit     int   `json:"limit"`
	}

	ComparisonOp string

	ComparisonFilter struct {
		Field string
		Op    ComparisonOp
		Value string
	}
)

const (
	GT  ComparisonOp = "gt"
	GTE ComparisonOp = "gte"
	LT  ComparisonOp = "lt"
	LTE ComparisonOp = "lte"
	EQ  ComparisonOp = "eq"
)

var valueSeparator = ","

func SetValueSeparator(separator string) {
	valueSeparator = separator
}

func ParseOpts(values url.Values) QueryOptions {
	opts := QueryOptions{
		Page:    1,
		Limit:   10,
		Filters: map[string]string{},
	}

	if p := values.Get("page"); p != "" {
		if page, _ := strconv.Atoi(p); page > 0 {
			opts.Page = page
		}
	}
	if l := values.Get("limit"); l != "" {
		if limit, _ := strconv.Atoi(l); limit > 0 {
			opts.Limit = limit
		}
	}
	opts.Offset = (opts.Page - 1) * opts.Limit

	if sort := values.Get("sort"); sort != "" {
		fields := strings.Split(sort, valueSeparator)
		for _, f := range fields {
			desc := strings.HasPrefix(f, "-")
			field := strings.TrimPrefix(f, "-")
			opts.Sort = append(opts.Sort, SortField{Field: field, Desc: desc})
		}
	}
	if fields := values.Get("search"); fields != "" {
		if keyword := values.Get("keyword"); keyword != "" {
			opts.Search = &SearchQuery{
				Fields:  strings.Split(fields, valueSeparator),
				Keyword: keyword,
			}
		}
	}
	if sel := values.Get("select"); sel != "" {
		opts.Select = strings.Split(sel, valueSeparator)
	}
	if group := values.Get("group"); group != "" {
		opts.GroupBy = strings.Split(group, valueSeparator)
		// force selection equal to group. this prevent full group by only
		opts.Select = strings.Split(group, valueSeparator)

		for _, sort := range opts.Sort {
			opts.GroupBy = append(opts.GroupBy, sort.Field)
		}
	}

	for key, val := range values {
		if key == "page" || key == "limit" || key == "sort" || key == "search" || key == "keyword" || key == "select" || key == "group" {
			continue
		}
		// Handle both searchAnd.field=keyword and search_and.field=keyword formats
		if (strings.HasPrefix(key, "searchAnd.") || strings.HasPrefix(key, "search_and.")) && len(val) > 0 {
			var fieldName string
			if strings.HasPrefix(key, "searchAnd.") {
				fieldName = strings.TrimPrefix(key, "searchAnd.")
			} else {
				fieldName = strings.TrimPrefix(key, "search_and.")
			}
			if fieldName != "" && val[0] != "" {
				if opts.SearchAnd == nil {
					opts.SearchAnd = &SearchQueryAnd{Fields: []*SearchField{}}
				}
				opts.SearchAnd.Fields = append(opts.SearchAnd.Fields, &SearchField{
					Field:   fieldName,
					Keyword: val[0],
				})
			}
			continue
		}
		if matches := regexp.MustCompile(`^([a-zA-Z0-9_]+)\[(gt|gte|lt|lte|eq)\]$`).FindStringSubmatch(key); len(matches) == 3 {
			opts.Comparisons = append(opts.Comparisons, ComparisonFilter{
				Field: matches[1],
				Op:    ComparisonOp(matches[2]),
				Value: val[0],
			})
			continue
		}
		opts.Filters[key] = val[0]
	}
	return opts
}

func ErrorPage(err error, opts QueryOptions) PageData {
	return PageData{
		LastError: err,
		Items:     []string{},
		Total:     0,
		Page:      1,
		Limit:     opts.Limit,
	}
}

// findFieldByColumn returns the struct field by matching the JSON tag or field name
func findFieldByColumn(v reflect.Value, column string) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName == column {
				return v.Field(i)
			}
		}
		if strings.EqualFold(field.Name, column) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

// compare is used for filtering values (uses ComparisonOp from external file)
func compare(fieldVal interface{}, strVal string, op ComparisonOp) bool {
	switch v := fieldVal.(type) {
	case string:
		return compareString(v, strVal, op)
	case typedef.Integer:
		return compareInt64(v.String(), strVal, op)
	case typedef.Float:
		return compareFloat64(v.String(), strVal, op)
	case typedef.Date:
		return compareTime(v.Time(), strVal, op)
	case typedef.Datetime:
		return compareTime(v.Time(), strVal, op)
	case int, int64:
		return compareInt64(toString(v), strVal, op)
	case float32, float64:
		return compareFloat64(toString(v), strVal, op)
	case time.Time:
		return compareTime(v, strVal, op)
	default:
		return false
	}
}

// Sorting comparator (used in sort.SliceStable)
func compareSort(a, b interface{}, desc bool) bool {
	switch va := a.(type) {
	case typedef.Integer:
		bi := b.(typedef.Integer)
		return sortInt64(va.Int64(), bi.Int64(), desc)
	case typedef.Float:
		bf := b.(typedef.Float)
		return sortFloat64(va.Float64(), bf.Float64(), desc)
	case typedef.Date:
		return sortTime(va.Time(), b.(typedef.Date).Time(), desc)
	case typedef.Datetime:
		return sortTime(va.Time(), b.(typedef.Datetime).Time(), desc)
	case string:
		return sortString(va, b.(string), desc)
	case int, int64:
		return sortInt64(toInt64(va), toInt64(b), desc)
	case float32, float64:
		return sortFloat64(toFloat64(va), toFloat64(b), desc)
	case time.Time:
		return sortTime(va, b.(time.Time), desc)
	default:
		return false
	}
}

// === Comparison helpers ===

func compareString(a, b string, op ComparisonOp) bool {
	switch op {
	case EQ:
		return a == b
	case GT:
		return a > b
	case GTE:
		return a >= b
	case LT:
		return a < b
	case LTE:
		return a <= b
	default:
		return false
	}
}

func compareInt64(aStr, bStr string, op ComparisonOp) bool {
	a, err1 := strconv.ParseInt(aStr, 10, 64)
	b, err2 := strconv.ParseInt(bStr, 10, 64)
	if err1 != nil || err2 != nil {
		return false
	}
	switch op {
	case EQ:
		return a == b
	case GT:
		return a > b
	case GTE:
		return a >= b
	case LT:
		return a < b
	case LTE:
		return a <= b
	default:
		return false
	}
}

func compareFloat64(aStr, bStr string, op ComparisonOp) bool {
	a, err1 := strconv.ParseFloat(aStr, 64)
	b, err2 := strconv.ParseFloat(bStr, 64)
	if err1 != nil || err2 != nil {
		return false
	}
	switch op {
	case EQ:
		return a == b
	case GT:
		return a > b
	case GTE:
		return a >= b
	case LT:
		return a < b
	case LTE:
		return a <= b
	default:
		return false
	}
}

func compareTime(a time.Time, bStr string, op ComparisonOp) bool {
	layout := "2006-01-02"
	if len(bStr) > 10 {
		layout = "2006-01-02 15:04:05"
	}
	b, err := time.Parse(layout, bStr)
	if err != nil {
		return false
	}
	switch op {
	case EQ:
		return a.Equal(b)
	case GT:
		return a.After(b)
	case GTE:
		return a.After(b) || a.Equal(b)
	case LT:
		return a.Before(b)
	case LTE:
		return a.Before(b) || a.Equal(b)
	default:
		return false
	}
}

// === Sorting helpers ===

func sortString(a, b string, desc bool) bool {
	if desc {
		return a > b
	}
	return a < b
}

func sortInt64(a, b int64, desc bool) bool {
	if desc {
		return a > b
	}
	return a < b
}

func sortFloat64(a, b float64, desc bool) bool {
	if desc {
		return a > b
	}
	return a < b
}

func sortTime(a, b time.Time, desc bool) bool {
	if desc {
		return a.After(b)
	}
	return a.Before(b)
}

// === Conversion helpers ===

func toInt64(v interface{}) int64 {
	switch i := v.(type) {
	case int:
		return int64(i)
	case int64:
		return i
	case string:
		n, _ := strconv.ParseInt(i, 10, 64)
		return n
	default:
		return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch f := v.(type) {
	case float32:
		return float64(f)
	case float64:
		return f
	case string:
		n, _ := strconv.ParseFloat(f, 64)
		return n
	default:
		return 0
	}
}

func toString(v interface{}) string {
	switch x := v.(type) {
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 64)
	default:
		return ""
	}
}

// === JSON tag extractor ===

func DefaultFilterByJson[T any]() map[string]string {
	var model T
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonTag := f.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			name := strings.Split(jsonTag, ",")[0]
			fields[name] = name
		} else {
			fields[strings.ToLower(f.Name)] = strings.ToLower(f.Name)
		}
	}
	return fields
}
