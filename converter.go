package slicer

import (
	"fmt"

	slicerpb "github.com/godev90/slicer/pb"
	"google.golang.org/protobuf/types/known/structpb"
)

func (q QueryOptions) ToProto() *slicerpb.QueryOptions {
	sort := make([]*slicerpb.SortField, 0, len(q.Sort))
	for _, s := range q.Sort {
		sort = append(sort, &slicerpb.SortField{
			Field: s.Field,
			Desc:  s.Desc,
		})
	}

	comparisons := make([]*slicerpb.ComparisonFilter, 0, len(q.Comparisons))
	for _, c := range q.Comparisons {
		comparisons = append(comparisons, &slicerpb.ComparisonFilter{
			Field: c.Field,
			Op:    string(c.Op), // 👈 Cast to string
			Value: c.Value,
		})
	}

	var search *slicerpb.SearchQuery
	if q.Search != nil && (len(q.Search.Fields) > 0 || q.Search.Keyword != "") {
		search = &slicerpb.SearchQuery{
			Fields:  q.Search.Fields,
			Keyword: q.Search.Keyword,
		}
	}

	return &slicerpb.QueryOptions{
		Page:        uint32(q.Page),
		Limit:       uint32(q.Limit),
		Sort:        sort,
		Search:      search,
		Select:      q.Select,
		Filters:     q.Filters,
		Comparisons: comparisons,
	}
}

func QueryFromProto(pb *slicerpb.QueryOptions) QueryOptions {
	if pb == nil {
		return QueryOptions{}
	}

	sort := make([]SortField, 0, len(pb.Sort))
	for _, s := range pb.Sort {
		sort = append(sort, SortField{
			Field: s.Field,
			Desc:  s.Desc,
		})
	}

	comparisons := make([]ComparisonFilter, 0, len(pb.Comparisons))
	for _, c := range pb.Comparisons {
		comparisons = append(comparisons, ComparisonFilter{
			Field: c.Field,
			Op:    ComparisonOp(c.Op), // 👈 Cast back to ComparisonOp
			Value: c.Value,
		})
	}

	var search *SearchQuery
	if pb.Search != nil {
		search = &SearchQuery{
			Fields:  pb.Search.Fields,
			Keyword: pb.Search.Keyword,
		}
	}

	return QueryOptions{
		Page:        int(pb.Page),
		Limit:       int(pb.Limit),
		Sort:        sort,
		Search:      search,
		Select:      pb.Select,
		Filters:     pb.Filters,
		Comparisons: comparisons,
	}
}

func (data PageData) DataToProto() (*slicerpb.PageResult, error) {
	items := []map[string]any{}

	switch v := data.Items.(type) {
	case []map[string]any:
		items = v
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				items = append(items, m)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported item type: %T", data.Items)
	}

	structs := make([]*structpb.Struct, 0, len(items))
	for _, item := range items {
		s, err := structpb.NewStruct(item)
		if err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}

	return &slicerpb.PageResult{
		Total: data.Total,
		Page:  int32(data.Page),
		Limit: int32(data.Limit),
		Items: structs,
	}, nil
}

func PageFromProto(pb *slicerpb.PageResult) PageData {
	items := make([]map[string]any, 0, len(pb.Items))
	for _, s := range pb.Items {
		items = append(items, s.AsMap())
	}

	return PageData{
		Items: items,
		Total: pb.Total,
		Page:  int(pb.Page),
		Limit: int(pb.Limit),
	}
}
