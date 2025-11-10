package slicer

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

	slicerpb "github.com/godev90/slicer/pb"
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
			Op:    string(c.Op),
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

	var searchAnd *slicerpb.SearchQueryAnd
	if q.SearchAnd != nil && len(q.SearchAnd.Fields) > 0 {
		searchFields := make([]*slicerpb.SearchField, 0, len(q.SearchAnd.Fields))
		for _, field := range q.SearchAnd.Fields {
			searchFields = append(searchFields, &slicerpb.SearchField{
				Field:   field.Field,
				Keyword: field.Keyword,
			})
		}
		searchAnd = &slicerpb.SearchQueryAnd{
			Fields: searchFields,
		}
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	if q.Page == 0 {
		q.Page = 1
	}

	return &slicerpb.QueryOptions{
		Page:        uint32(q.Page),
		Limit:       uint32(q.Limit),
		Sort:        sort,
		Search:      search,
		SearchAnd:   searchAnd,
		Select:      q.Select,
		Filters:     q.Filters,
		GroupBy:     q.GroupBy,
		Comparisons: comparisons,
	}
}

func QueryFromProto(pb *slicerpb.QueryOptions) QueryOptions {
	if pb == nil {
		return QueryOptions{
			Limit: 10,
			Page:  1,
		}
	}

	if pb.Limit == 0 {
		pb.Limit = 10
	}

	if pb.Page == 0 {
		pb.Page = 1
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
			Op:    ComparisonOp(c.Op),
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

	var searchAnd *SearchQueryAnd
	if pb.SearchAnd != nil && len(pb.SearchAnd.Fields) > 0 {
		searchFields := make([]*SearchField, 0, len(pb.SearchAnd.Fields))
		for _, field := range pb.SearchAnd.Fields {
			searchFields = append(searchFields, &SearchField{
				Field:   field.Field,
				Keyword: field.Keyword,
			})
		}
		searchAnd = &SearchQueryAnd{
			Fields: searchFields,
		}
	}

	return QueryOptions{
		Page:        int(pb.Page),
		Limit:       int(pb.Limit),
		Offset:      int(pb.Page-1) * int(pb.Limit),
		Sort:        sort,
		Search:      search,
		SearchAnd:   searchAnd,
		Select:      pb.Select,
		GroupBy:     pb.GroupBy,
		Filters:     pb.Filters,
		Comparisons: comparisons,
	}
}

// func (data PageData) ToProto() (*slicerpb.PageData, error) {
// 	jbytes, err := json.Marshal(data.Items)

// 	if err != nil {
// 		return nil, err
// 	}

// 	page := int32(data.Page)
// 	limit := int32(data.Limit)

// 	if page == 0 {
// 		page = 1
// 	}

// 	if limit == 0 {
// 		limit = 10
// 	}

// 	return &slicerpb.PageData{
// 		Page:  page,
// 		Limit: limit,
// 		Total: data.Total,
// 		Items: jbytes,
// 	}, nil
// }

func (data PageData) ToProto() (*slicerpb.PageData, error) {
	jbytes, err := json.Marshal(data.Items)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err = gw.Write(jbytes)
	if err != nil {
		_ = gw.Close()
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	compressed := buf.Bytes()

	page := int32(data.Page)
	limit := int32(data.Limit)

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	return &slicerpb.PageData{
		Page:  page,
		Limit: limit,
		Total: data.Total,
		Items: compressed,
	}, nil
}

// func PageFromProto(protoData *slicerpb.PageData, destSchema any) (*PageData, error) {
// 	if protoData == nil {
// 		return nil, nil
// 	}

// 	err := json.Unmarshal(protoData.Items, destSchema)
// 	if err != nil {
// 		return nil, err
// 	}

// 	page := protoData.Page
// 	limit := protoData.Limit

// 	if page == 0 {
// 		page = 1
// 	}

// 	if limit == 0 {
// 		limit = 10
// 	}

// 	return &PageData{
// 		Page:  int(page),
// 		Limit: int(limit),
// 		Total: protoData.Total,
// 		Items: destSchema,
// 	}, nil
// }

func PageFromProto(protoData *slicerpb.PageData, destSchema any) (*PageData, error) {
	if protoData == nil {
		return nil, nil
	}

	items := protoData.Items
	// detect gzip by magic bytes 0x1f 0x8b
	if len(items) >= 2 && items[0] == 0x1f && items[1] == 0x8b {
		gr, err := gzip.NewReader(bytes.NewReader(items))
		if err != nil {
			return nil, err
		}
		decompressed, err := io.ReadAll(gr)
		_ = gr.Close()
		if err != nil {
			return nil, err
		}
		items = decompressed
	}

	err := json.Unmarshal(items, destSchema)
	if err != nil {
		return nil, err
	}

	page := protoData.Page
	limit := protoData.Limit

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	return &PageData{
		Page:  int(page),
		Limit: int(limit),
		Total: protoData.Total,
		Items: destSchema,
	}, nil
}
