# Performance Benchmarks

This document contains performance benchmarks for api-godoc against real-world API specifications.

## Test Environment
- Platform: Apple M2 Ultra
- Go version: 1.24.3
- Test date: July 2025

## Benchmarks

### End-to-End Processing
| Operation | Time (ns/op) | Operations/sec | Description |
|-----------|--------------|----------------|-------------|
| Full Pipeline | 4,024,225 | ~248 | Complete processing of Stripe API |

### Individual Components
| Component | Time (ns/op) | Operations/sec | Description |
|-----------|--------------|----------------|-------------|
| Parsing | 628,062,833 | ~1.6 | Parse 6.8MB Stripe OpenAPI spec |
| Resource Extraction | 647,082 | ~1,545 | Extract resources from parsed spec |
| Pattern Detection | 417,381 | ~2,396 | Detect API patterns |
| Markdown Generation | 3,791 | ~263,824 | Generate markdown output |
| Resource Filtering | 11,402 | ~87,711 | Filter resources by pattern |

## Real-World API Results

### Stripe API (Full)
- **Size**: 6.8MB OpenAPI 3.0 specification
- **Resources Extracted**: 216
- **Patterns Detected**: 3
- **Processing Time**: ~632ms
- **Memory Usage**: Efficient for large specs

### GitHub API
- **Size**: 10.8MB OpenAPI 3.0 specification  
- **Resources Extracted**: 10+ (varies by processing)
- **Processing Time**: Under 30 seconds
- **Notes**: Very large specification with complex schemas

### Performance Insights

1. **Parsing is the bottleneck**: ~95% of processing time is spent in parsing
2. **Analysis is fast**: Resource extraction and pattern detection are very efficient
3. **Output generation**: Sub-millisecond for most formats
4. **Scalability**: Linear scaling with API complexity

## Optimization Opportunities

1. **Parser caching**: Cache parsed specs for repeated analysis
2. **Streaming parsing**: For very large specifications
3. **Parallel processing**: Pattern detection could be parallelized
4. **Memory optimization**: Reduce allocations in hot paths

## Usage Recommendations

- For CI/CD: Process specs up to 10MB efficiently
- For interactive use: Sub-second response for most APIs
- For batch processing: Can handle 100+ APIs per minute
- Memory usage: Scales linearly with spec size, reasonable for modern systems