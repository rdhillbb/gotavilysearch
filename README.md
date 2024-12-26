# Function Documentation for Tavily Search Package

## Overview

This package provides a concurrent search implementation using the Tavily API. It supports multiple search modes and uses Go's errgroup package for efficient concurrent execution.

## Core Components

### searchTavily.go Functions

#### DeepSearch
```go
func DeepSearch(searchQuery string) (string, error)
```
Performs an expanded search operation by breaking down the query into multiple sub-queries and aggregating results.

**Parameters:**
- `searchQuery` (string): The original search query to be expanded and searched

**Returns:**
- `string`: Combined search results from all sub-queries
- `error`: Error if the search fails

**Error cases:**
- Empty query provided
- Query expansion failure
- All searches failed
- No results found

#### SearchInternet
```go
func SearchInternet(searchQuery string) (string, error)
```
Performs a single search query against the Tavily API.

**Parameters:**
- `searchQuery` (string): The search query to execute

**Returns:**
- `string`: Search results in JSON format
- `error`: Error if the search fails

**Error cases:**
- Empty query provided
- Search execution failed
- No results returned
- Empty results returned

#### SearchAtentInternet
```go
func SearchAtentInternet(searchQueries []string, searchFunc func(string, int) (string, error), envVar string)
```
Manages concurrent execution of multiple search queries using errgroup.Group.

**Parameters:**
- `searchQueries` ([]string): Array of search queries to execute
- `searchFunc` (func(string, int) (string, error)): Function to perform individual searches
- `envVar` (string): Environment variable name for max results configuration

**Returns:**
- `[]SearchResult`: Array of search results with timing and error information

**Concurrency Features:**
- Uses errgroup.Group for goroutine management
- Mutex-based synchronization for result collection
- Context-based timeout (30 seconds)
- Automatic cleanup and cancellation

### tavilyclient.go Functions

#### TavilySearch
```go
func TavilySearch(searchQuery string, maxResults int) (string, error)
```
Performs a general internet search using the Tavily API.

**Parameters:**
- `searchQuery` (string): The search query to execute
- `maxResults` (int): Maximum number of results to return

**Returns:**
- `string`: JSON-formatted search results
- `error`: Error if the search fails

**Configuration:**
- Topic: General
- IncludeAnswer: true
- IncludeImages: false
- IncludeRawContent: false

#### TavilyMaxSearch
```go
func TavilyMaxSearch(searchQuery string) (string, error)
```
Performs a deep search with advanced depth configuration.

**Configuration:**
- Topic: General
- SearchDepth: advance
- IncludeAnswer: true
- IncludeRawContent: true
- MaxResults: From MAXDeepRESULTS environment variable

#### TavilyNewsSearch
```go
func TavilyNewsSearch(searchQuery string) (string, error)
```
Performs a basic news search.

**Configuration:**
- Topic: News
- SearchDepth: basic
- IncludeAnswer: true
- MaxResults: From MAXRESULTS environment variable

#### TavilyRawNewsSearch
```go
func TavilyRawNewsSearch(searchQuery string) (string, error)
```
Performs an advanced news search with expanded results.

**Configuration:**
- Topic: News
- SearchDepth: advance
- IncludeAnswer: true
- MaxResults: From MAXDeepRESULTS environment variable
- Days: 5

#### API Key and Configuration Management
```go
func getAPIKey() string
func getMaxResults(envVar string) int
```

**API Key Management:**
1. Checks environment variable TAVILY_API_KEY
2. Falls back to default API key if environment variable is not set

**Max Results Configuration:**
1. Reads from specified environment variable (MAXRESULTS or MAXDeepRESULTS)
2. Falls back to default value of 3 if not set

## Data Structures

### SearchResult
```go
type SearchResult struct {
    Query      string        `json:"query"`
    Results    string        `json:"results"`
    WorkerID   int          `json:"worker_id"`
    TimeSpent  time.Duration `json:"time_spent"`
    Error      error        `json:"error"`
}
```
Structure for holding search results and metadata.

## Error Handling Patterns

1. Input Validation
```go
if searchQuery == "" {
    return "", fmt.Errorf("empty search query provided")
}
```

2. Error Wrapping
```go
if err != nil {
    return "", fmt.Errorf("search failed: %v", err)
}
```

3. Concurrent Error Handling
```go
if err := g.Wait(); err != nil {
    logging.WriteLogs(fmt.Sprintf("Search group error: %v", err))
}
```

## Logging

The package implements comprehensive logging throughout:
- Search initialization and completion
- Worker status updates and timing
- Error conditions
- API interactions

## Dependencies

- golang.org/x/sync/errgroup
- github.com/hekmon/tavily
- Internal packages:
  - anthropicfunc/logging
  - anthropicfunc/util

## Configuration

### Environment Variables
- TAVILY_API_KEY: API key for Tavily service
- MAXRESULTS: Maximum results for regular searches
- MAXDeepRESULTS: Maximum results for deep searches

### Default Settings
- Default API key provided as fallback
- Default max results: 3
- Search timeout: 30 seconds
- News search default days: 5

### Concurrency Settings
- Worker count matches query count
- Mutex-protected result collection
- Context-based timeout management
