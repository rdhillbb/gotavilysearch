package gotavilysearch

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
    "golang.org/x/sync/errgroup"
    "github.com/rdhillbb/logging"
    "github.com/rdhillbb/util"
)

type SearchResult struct {
    Query      string        `json:"query"`
    Results    string        `json:"results"`
    WorkerID   int          `json:"worker_id"`
    TimeSpent  time.Duration `json:"time_spent"`
    Error      error        `json:"error"`
}

func DeepSearch(searchQuery string) (string, error) {
    if searchQuery == "" {
        return "", fmt.Errorf("empty search query provided")
    }

    expandedQueries, err := util.ReWriteQR(searchQuery)
    if err != nil {
        return "", fmt.Errorf("failed to expand query: %v", err)
    }

    searchResults := SearchAtentInternet(expandedQueries, TavilySearch, "MAXDeepRESULTS")
    
    var searchErrors []string
    var combinedResults string
    
    for _, result := range searchResults {
        if result.Error != nil {
            searchErrors = append(searchErrors, fmt.Sprintf("worker %d: %v", result.WorkerID, result.Error))
        }
        if result.Results != "" {
            combinedResults += result.Results + " "
        }
    }
    
    if len(searchErrors) == len(searchResults) {
        return "", fmt.Errorf("all searches failed: %v", searchErrors)
    }
    
    if len(searchErrors) > 0 {
        logging.WriteLogs(fmt.Sprintf("Some searches failed: %v", searchErrors))
    }
    
    if combinedResults == "" {
        return "", fmt.Errorf("no results found")
    }
    
    fmt.Println(combinedResults)
    return combinedResults, nil
}

func SearchInternet(searchQuery string) (string, error) {
    if searchQuery == "" {
        return "", fmt.Errorf("empty search query provided")
    }

    singleQuery := []string{searchQuery}
    searchResults := SearchAtentInternet(singleQuery, TavilySearch, "MAXRESULTS")
    
    if len(searchResults) == 0 {
        return "", fmt.Errorf("no results returned from search")
    }
    
    if searchResults[0].Error != nil {
        return "", fmt.Errorf("search failed: %v", searchResults[0].Error)
    }
    
    if searchResults[0].Results == "" {
        return "", fmt.Errorf("empty result returned from search")
    }
    
    return searchResults[0].Results, nil
}

func SearchAtentInternet(searchQueries []string, searchFunc func(string, int) (string, error), envVar string) []SearchResult {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    g, gctx := errgroup.WithContext(ctx)
    
    var mu sync.Mutex
    var collectedResults []SearchResult
    searchStartTime := time.Now()
    
    logging.WriteLogs(fmt.Sprintf("\nStarting search with %d workers at %v", 
        len(searchQueries), searchStartTime.Format("15:04:05")))
    logging.WriteLogs("Timeout set to 30 seconds\n")
    
    for workerID, query := range searchQueries {
        workerID, query := workerID, query // Create new variables for goroutine closure
        
        g.Go(func() error {
            workerStartTime := time.Now()
            logging.WriteLogs(fmt.Sprintf("Worker %d starting search with query: %s", workerID, query))
            
            result := SearchResult{
                Query:    query,
                WorkerID: workerID,
            }
            
            // Check if context is already done before starting search
            select {
            case <-gctx.Done():
                result.Error = fmt.Errorf("search canceled before start: %v", gctx.Err())
                result.TimeSpent = time.Since(workerStartTime)
            default:
                searchResultText, searchErr := searchFunc(query, getMaxResults(envVar))
                result.TimeSpent = time.Since(workerStartTime)
                result.Results = searchResultText
                result.Error = searchErr
                
                if searchErr != nil {
                    logging.WriteLogs(fmt.Sprintf("Worker %d recorded error: %v", workerID, searchErr))
                } else {
                    logging.WriteLogs(fmt.Sprintf("Worker %d completed search in %v", workerID, result.TimeSpent))
                }
            }
            
            mu.Lock()
            collectedResults = append(collectedResults, result)
            mu.Unlock()
            
            return nil // Don't propagate search errors to allow other searches to continue
        })
    }
    
    if err := g.Wait(); err != nil {
        logging.WriteLogs(fmt.Sprintf("Search group error: %v", err))
    }
    
    totalSearchTime := time.Since(searchStartTime)
    logging.WriteLogs(fmt.Sprintf("\nAll searches completed in %v", totalSearchTime))
    
    return collectedResults
}

func TESTmain() {
    rand.Seed(time.Now().UnixNano())
    
    testQueries := []string{
        "What is concurrent programming?",
        "How does concurrent programming differ from parallel programming, and what are its real-world applications?",
        "What are the fundamental principles and challenges of concurrent programming in modern software development?",
        "How do different programming languages handle concurrent programming, and what are their unique approaches?",
        "What are the design patterns and best practices for implementing concurrent systems in software architecture?",
        "How has concurrent programming evolved with multi-core processors, and what are its performance implications?",
        "What are the common pitfalls and synchronization issues in concurrent programming, and how are they resolved?",
        "How do concurrent programming models like actor-based and CSP compare in solving modern computing challenges?",
        "What role does concurrent programming play in distributed systems and cloud computing architectures?",
    }
    
    logging.WriteLogs("Starting concurrent programming research...")
    testResults := SearchAtentInternet(testQueries, TavilySearch, "MAXRESULTS")
    
    logging.WriteLogs("\nFinal Results Summary:")
    logging.WriteLogs("----------------------")
    for _, testResult := range testResults {
        if testResult.Error != nil {
            logging.WriteLogs(fmt.Sprintf("\nWorker %d: Failed - Error: %v\nQuery: %s", 
                testResult.WorkerID, testResult.Error, testResult.Query))
        } else {
            logging.WriteLogs(fmt.Sprintf("\nWorker %d: Completed in %v\nQuery: %s\nResults: %s", 
                testResult.WorkerID, 
                testResult.TimeSpent,
                testResult.Query,
                testResult.Results))
        }
    }
}
