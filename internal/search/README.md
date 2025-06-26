# Search Package

The search package provides functionality to search through product families using various criteria and filters. It implements a faceted search system backed by Elasticsearch.

- **product family** - product group aggregated by some properties: color, condition, size etc

## Features

### Product Family Search
- Full-text search across product families
- Faceted filtering
- Pagination support
- Sorting options

### Build facets
- fetches filters (facets) with possible options for next filtering 

## API

### gRPC Interface

```protobuf
service SearchService {
    // Search for product families with filters
    rpc Items(ItemsRequest) returns (ItemsResult);
    
    // Get available filters
    rpc Filters(FiltersRequest) returns (FiltersResult);
}
```

### Request/Response Examples

See:
- `tests/filters/**/fixture` or `tests/filters/**/expected` for Fiters API 
- `tests/items/**/fixture` or `tests/items/**/expected` for Items API 

## Usage

Before use your should init package dependency.

```go
import (
    "gitlab.trgdev.com/gotrg/white-label/services/catalog/internal/search/config"
)

func main() {
    config.Init(<EsClientInstance>)
}
```

## Future Improvements

- [ ] Add fuzzy search support
- [ ] Implement result highlighting
- [ ] Add aggregation caching
- [ ] Support multi-language search
- [ ] Support multi-curency context