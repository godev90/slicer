# slicer

**slicer** is a versatile and performant Go library for working with slices. It provides a rich set of functional-style methods to simplify operations like filtering, mapping, folding, finding, and more—without sacrificing speed or readability.

---

## ✨ Features

- Filter, map, reduce, any, all, find, exists, chunk, unique, reverse, and more
- Works with generic slice types (`[]T`)
- Immutable operations (returns new slices)
- Extensible: Chain multiple calls for functional-style pipelines
- Zero dependencies — no reliance on reflect or third-party libraries
- Advanced pagination for both in-memory slices and database queries
- Support for OR and AND search operations
- Protocol Buffer integration for serialization

---

## 🔍 Search Operations

### Regular Search (OR Logic)
The existing `search` parameter applies OR logic between fields using a single keyword:

```go
// URL: ?search=name,description&keyword=golang
// SQL: WHERE (name ILIKE '%golang%' OR description ILIKE '%golang%')

opts := QueryOptions{
    Search: &SearchQuery{
        Fields:  []string{"name", "description"},
        Keyword: "golang",
    },
}
```

### AND Search with Field-Specific Keywords
The new `search_and` parameter allows you to specify different keywords for different fields with AND logic:

```go
// URL: ?search_and.status=active&search_and.category=backend
// SQL: WHERE status ILIKE '%active%' AND category ILIKE '%backend%'

opts := QueryOptions{
    SearchAnd: &SearchQueryAnd{
        Fields: []*SearchField{
            {Field: "status", Keyword: "active"},
            {Field: "category", Keyword: "backend"},
        },
    },
}
```

### Combined Search
You can use both search types together:

```go
// URL: ?search=name,description&keyword=golang&search_and.status=active&search_and.category=backend
// SQL: WHERE (name ILIKE '%golang%' OR description ILIKE '%golang%') 
//       AND status ILIKE '%active%' 
//       AND category ILIKE '%backend%'

opts := QueryOptions{
    Search: &SearchQuery{
        Fields:  []string{"name", "description"},
        Keyword: "golang",
    },
    SearchAnd: &SearchQueryAnd{
        Fields: []*SearchField{
            {Field: "status", Keyword: "active"},
            {Field: "category", Keyword: "backend"},
        },
    },
}
```

---