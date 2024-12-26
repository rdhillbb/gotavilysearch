package gotavilysearch

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "github.com/hekmon/tavily"
    "anthropicfunc/logging"
)

const defaultAPIKey = "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"

func getAPIKey() string {
    if apiKey := os.Getenv("TAVILY_API_KEY"); apiKey != "" {
        return apiKey
    }
    return defaultAPIKey
}

func getMaxResults(envVar string) int {
    if maxResults := os.Getenv(envVar); maxResults != "" {
        if val, err := strconv.Atoi(maxResults); err == nil {
            return val
        }
    }
    return 3 // Default value
}

func convertToString(searchAnswer tavily.SearchAnswer) (string, error) {
    serializedJson, err := json.Marshal(searchAnswer)
    if err != nil {
        logging.WriteLogs("Error converting search answer to string: " + err.Error())
        return "", fmt.Errorf("failed to serialize search results: %v", err)
    }
    return string(serializedJson), nil
}

func TavilySearch(searchQuery string, maxResults int) (string, error) {
    if searchQuery == "" {
        return "", fmt.Errorf("empty search query provided")
    }

    logging.WriteLogs("Starting Tavily search with query: " + searchQuery)
    tavilyClient := tavily.NewClient(getAPIKey(), nil)
    
    searchAnswer, err := tavilyClient.Search(context.TODO(), tavily.SearchQuery{
        Query:                    searchQuery,
        Topic:                    tavily.SearchQueryTopicGeneral,
        IncludeAnswer:            true,
        IncludeImages:            false,
        IncludeImageDescriptions: false,
        IncludeRawContent:        false,
        MaxResults:               maxResults,
    })
    if err != nil {
        logging.WriteLogs("Tavily search error: " + err.Error())
        return "", fmt.Errorf("search failed: %v", err)
    }

    jsonResponse, err := convertToString(searchAnswer)
    if err != nil {
        logging.WriteLogs("Error converting answer to string: " + err.Error())
        return "", err
    }
    
    logging.WriteLogs("Search completed successfully")
    return jsonResponse, nil
}

func TavilyMaxSearch(searchQuery string) (string, error) {
    if searchQuery == "" {
        return "", fmt.Errorf("empty search query provided")
    }

    logging.WriteLogs("Starting deep Tavily search with query: " + searchQuery)
    tavilyClient := tavily.NewClient(getAPIKey(), nil)
    
    searchAnswer, err := tavilyClient.Search(context.TODO(), tavily.SearchQuery{
        Query:                    searchQuery,
        Topic:                    tavily.SearchQueryTopicGeneral,
        SearchDepth:              "advance",
        IncludeAnswer:            true,
        IncludeImages:            false,
        IncludeImageDescriptions: false,
        IncludeRawContent:        true,
        MaxResults:               getMaxResults("MAXDeepRESULTS"),
    })
    if err != nil {
        logging.WriteLogs("Tavily search error: " + err.Error())
        return "", fmt.Errorf("deep search failed: %v", err)
    }
    
    jsonResponse, err := convertToString(searchAnswer)
    if err != nil {
        logging.WriteLogs("Error converting answer to string: " + err.Error())
        return "", err
    }
    
    logging.WriteLogs("Deep search completed successfully")
    return jsonResponse, nil
}

func TavilyNewsSearch(searchQuery string) (string, error) {
    if searchQuery == "" {
        return "", fmt.Errorf("empty search query provided")
    }

    logging.WriteLogs("Starting Tavily news search with query: " + searchQuery)
    tavilyClient := tavily.NewClient(getAPIKey(), nil)
    
    searchAnswer, err := tavilyClient.Search(context.TODO(), tavily.SearchQuery{
        Query:                    searchQuery,
        Topic:                    tavily.SearchQueryTopicNews,
        SearchDepth:              "basic",
        IncludeAnswer:            true,
        IncludeImages:            false,
        IncludeImageDescriptions: false,
        IncludeRawContent:        false,
        MaxResults:               getMaxResults("MAXRESULTS"),
    })
    if err != nil {
        logging.WriteLogs("Tavily search error: " + err.Error())
        return "", fmt.Errorf("news search failed: %v", err)
    }

    jsonResponse, err := convertToString(searchAnswer)
    if err != nil {
        logging.WriteLogs("Error converting answer to string: " + err.Error())
        return "", err
    }

    logging.WriteLogs("News search completed successfully")
    return jsonResponse, nil
}

func TavilyRawNewsSearch(searchQuery string) (string, error) {
    if searchQuery == "" {
        return "", fmt.Errorf("empty search query provided")
    }

    logging.WriteLogs("Starting raw Tavily news search with query: " + searchQuery)
    tavilyClient := tavily.NewClient(getAPIKey(), nil)
    
    searchAnswer, err := tavilyClient.Search(context.TODO(), tavily.SearchQuery{
        Query:                    searchQuery,
        Topic:                    tavily.SearchQueryTopicNews,
        SearchDepth:              "advance",
        IncludeAnswer:            true,
        IncludeImages:            false,
        IncludeImageDescriptions: false,
        IncludeRawContent:        false,
        MaxResults:               getMaxResults("MAXDeepRESULTS"),
        Days:                     5,
    })
    if err != nil {
        logging.WriteLogs("Tavily search error: " + err.Error())
        return "", fmt.Errorf("raw news search failed: %v", err)
    }

    jsonResponse, err := convertToString(searchAnswer)
    if err != nil {
        logging.WriteLogs("Error converting answer to string: " + err.Error())
        return "", err
    }

    logging.WriteLogs("Raw news search completed successfully")
    return jsonResponse, nil
}
