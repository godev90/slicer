package slicer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type SlicePaginator[T any] struct {
	source []T
	items  []T
	fields map[string]string
}

func NewSlicePaginator[T any](source []T, allowedFields map[string]string) *SlicePaginator[T] {
	return &SlicePaginator[T]{
		source: source,
		fields: allowedFields,
		items:  []T{}, // Initialize with empty slice instead of nil
	}
}

// NewSlicePaginator creates and returns a new SlicePaginator for the provided
// source slice. `allowedFields` is a map from logical field names to column
// names that will be used when filtering, searching and sorting.
// The paginator initializes with an empty items slice.

func (p *SlicePaginator[T]) Items() []T {
	return p.items
}

// Items returns the current page items stored in the paginator. It returns an
// empty slice when no items have been set.

func (p *SlicePaginator[T]) SetItems(items []T) {
	p.items = items
}

// SetItems sets the paginator's items to the provided slice. This is used by
// pagination routines to store the resulting page.

func SlicePage[T any](p *SlicePaginator[T], opts QueryOptions) (PageData, error) {
	var filtered []T

	opts.Offset = (opts.Page - 1) * opts.Limit

	// 1. Apply ComparisonFilters
	for _, item := range p.source {
		match := true
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for _, cmp := range opts.Comparisons {
			if _, ok := p.fields[cmp.Field]; !ok {
				continue
			}
			field := findFieldByColumn(v, cmp.Field)
			if !field.IsValid() || !compare(field.Interface(), cmp.Value, cmp.Op) {
				match = false
				break
			}
		}

		if match {
			filtered = append(filtered, item)
		}
	}

	if len(opts.Filters) > 0 {
		var temp []T
		for _, item := range filtered {
			v := reflect.ValueOf(item)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			match := true
			for field, val := range opts.Filters {
				if _, ok := p.fields[field]; !ok {
					continue
				}
				f := findFieldByColumn(v, field)
				if !f.IsValid() {
					match = false
					break
				}
				actual := fmt.Sprintf("%v", f.Interface())
				if actual != val {
					match = false
					break
				}
			}
			if match {
				temp = append(temp, item)
			}
		}
		filtered = temp
	}

	// 3. Apply search
	if opts.Search != nil {
		var searched []T
		for _, item := range filtered {
			v := reflect.ValueOf(item)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			matched := false
			for _, key := range opts.Search.Fields {
				if _, ok := p.fields[key]; !ok {
					continue
				}
				field := findFieldByColumn(v, key)
				if field.IsValid() {
					val := fmt.Sprintf("%v", field.Interface())
					if strings.Contains(strings.ToLower(val), strings.ToLower(opts.Search.Keyword)) {
						matched = true
						break
					}
				}
			}
			if matched {
				searched = append(searched, item)
			}
		}
		filtered = searched
	}

	// 3.1. Apply search_and
	if opts.SearchAnd != nil && len(opts.SearchAnd.Fields) > 0 {
		var searchedAnd []T
		for _, item := range filtered {
			v := reflect.ValueOf(item)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			matchedAll := true
			for _, searchField := range opts.SearchAnd.Fields {
				if _, ok := p.fields[searchField.Field]; !ok {
					continue
				}
				field := findFieldByColumn(v, searchField.Field)
				if field.IsValid() {
					val := fmt.Sprintf("%v", field.Interface())
					if !strings.Contains(strings.ToLower(val), strings.ToLower(searchField.Keyword)) {
						matchedAll = false
						break
					}
				} else {
					matchedAll = false
					break
				}
			}
			if matchedAll {
				searchedAnd = append(searchedAnd, item)
			}
		}
		filtered = searchedAnd
	}

	// 4. Sorting
	for i := len(opts.Sort) - 1; i >= 0; i-- {
		sortField := opts.Sort[i]
		if _, ok := p.fields[sortField.Field]; !ok {
			continue
		}

		sort.SliceStable(filtered, func(i, j int) bool {
			vi := reflect.ValueOf(filtered[i])
			vj := reflect.ValueOf(filtered[j])
			if vi.Kind() == reflect.Ptr {
				vi = vi.Elem()
				vj = vj.Elem()
			}
			fi := findFieldByColumn(vi, sortField.Field)
			fj := findFieldByColumn(vj, sortField.Field)

			return compareSort(fi.Interface(), fj.Interface(), sortField.Desc)
		})
	}

	// 5. Pagination
	total := len(filtered)
	start := opts.Offset
	end := opts.Offset + opts.Limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	pageItems := filtered[start:end]
	p.SetItems(pageItems)

	if p.Items() == nil {
		p.SetItems([]T{})
	}

	return PageData{
		Items: p.Items(),
		Total: int64(total),
		Page:  opts.Page,
		Limit: opts.Limit,
	}, nil
}

// SlicePage applies the provided QueryOptions to the paginator's source data
// and returns a PageData containing the resulting page slice, total count,
// and pagination metadata. The function performs comparisons, filters,
// search (including search AND), sorting and pagination in that order. The
// paginator's items are updated with the selected page slice.
