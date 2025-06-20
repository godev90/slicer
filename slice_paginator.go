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
	}
}

func (p *SlicePaginator[T]) Items() []T {
	return p.items
}

func (p *SlicePaginator[T]) SetItems(items []T) {
	p.items = items
}

func SlicePage[T any](p *SlicePaginator[T], opts QueryOptions) (PageData, error) {
	var filtered []T

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

	// 2. Apply search
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

	// 3. Sorting
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

	// 4. Pagination
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

	return PageData{
		Items: p.Items(),
		Total: int64(total),
		Page:  opts.Page,
		Limit: opts.Limit,
	}, nil
}
