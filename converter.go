package slicer

import (
	"encoding/json"

	slicerpb "github.com/godev90/slicer/pb"
	"github.com/golang/snappy"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ToProto converts a QueryOptions value from the slicer package into its
// protobuf representation `slicerpb.QueryOptions`.
//
// It maps sorting, comparisons, search fields and other query options into
// the generated protobuf message so the query may be transmitted over RPC.
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

// QueryFromProto converts a protobuf `slicerpb.QueryOptions` into the
// local `QueryOptions` type used by the slicer package.
//
// It performs safe defaults for page and limit and converts nested fields
// such as Sort and Comparisons back into their native types.
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

// ToProto serializes PageData.Items into JSON and returns a protobuf
// `slicerpb.PageData` containing the serialized (and snappy-compressed)
// bytes in the Items field. It also ensures sane defaults for page and limit.
func (data PageData) ToProto() (*slicerpb.PageData, error) {
	jbytes, err := json.Marshal(data.Items)

	if err != nil {
		return nil, err
	}

	compressed := snappy.Encode(nil, jbytes)

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

// PageFromProto converts a protobuf `slicerpb.PageData` back into the
// local PageData type. It will attempt to decode Snappy-compressed bytes
// in the Items field, fall back to raw JSON if decoding fails, and
// unmarshal the resulting JSON into destSchema.
func PageFromProto(protoData *slicerpb.PageData, destSchema any) (*PageData, error) {
	if protoData == nil {
		return nil, nil
	}

	var dataBytes []byte
	if len(protoData.Items) == 0 {
		dataBytes = nil
	} else {
		// try snappy decode; if fails, assume raw JSON
		if dec, err := snappy.Decode(nil, protoData.Items); err == nil {
			dataBytes = dec
		} else {
			dataBytes = protoData.Items
		}
	}

	err := json.Unmarshal(dataBytes, destSchema)
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

// ToProtoBuf packs PageData.Items as JSON bytes into a protobuf Any
// wrapper (wrapperspb.BytesValue) and stores it in the Rows field of the
// returned `slicerpb.PageData`. This is useful for transporting arbitrary
// JSON while keeping compatibility with consumers that expect an Any.
func (data PageData) ToProtoBuf() (*slicerpb.PageDataBuf, error) {
	// marshal Items to JSON bytes and pack into Any as wrapperspb.BytesValue
	jbytes, err := json.Marshal(data.Items)
	if err != nil {
		return nil, err
	}

	bytesMsg := &wrapperspb.BytesValue{Value: jbytes}
	anyVal, err := anypb.New(bytesMsg)
	if err != nil {
		return nil, err
	}

	page := int32(data.Page)
	limit := int32(data.Limit)

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	return &slicerpb.PageDataBuf{
		Page:  page,
		Limit: limit,
		Total: data.Total,
		Items: anyVal,
	}, nil
}

// PageFromProtoBuf unpacks the Rows Any field (expected to contain
// a wrapperspb.BytesValue), decodes the contained JSON bytes and
// unmarshals them into destSchema. Returns a populated PageData with
// defaults applied for page and limit.
func PageFromProtoBuf(protoData *slicerpb.PageDataBuf, destSchema any) (*PageData, error) {
	if protoData == nil {
		return nil, nil
	}

	if protoData.Items == nil {
		return &PageData{
			Page:  int(protoData.Page),
			Limit: int(protoData.Limit),
			Total: protoData.Total,
			Items: []string{},
		}, nil
	}

	// Unpack Any as wrapperspb.BytesValue and decode JSON
	var bytesMsg wrapperspb.BytesValue
	if err := protoData.Items.UnmarshalTo(&bytesMsg); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytesMsg.Value, destSchema); err != nil {
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
