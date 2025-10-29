package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "sync"
    "time"
)

type Post struct {
    ID        int64     `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Tags      []string  `json:"tags"`
    Version   int32     `json:"version"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type PostResponse struct {
    Data Post `json:"data"`
}

func main() {
    const (
        baseURL    = "http://localhost:4000/v1"
        postID     = "3"
        numWorkers = 5
    )

    fmt.Printf("Starting concurrency test with %d workers on post ID %s\n\n", numWorkers, postID)
    
    // Get initial post
    initialPost, err := getPost(baseURL, postID)
    if err != nil {
        fmt.Printf("Error getting initial post: %v\n", err)
        return
    }
    
    fmt.Println("Initial post state:")
    printPost(initialPost)
    fmt.Println()

    var wg sync.WaitGroup
    results := make(chan string, numWorkers)

    // All workers start with the same initial version
    initialPostForWorkers, err := getPost(baseURL, postID)
    if err != nil {
        fmt.Printf("Error getting initial post for workers: %v\n", err)
        return
    }

    // Start workers with the same initial version
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            // Add a small delay to increase chance of conflicts
            time.Sleep(time.Duration(workerID) * 50 * time.Millisecond)
            
            // Prepare update with the initial version
            update := map[string]interface{}{
                "title":   fmt.Sprintf("Updated by worker %d", workerID),
                "content": fmt.Sprintf("Content updated by worker %d at %v", workerID, time.Now().Format(time.RFC3339Nano)),
                "tags":    []string{"concurrency", "test", fmt.Sprintf("worker-%d", workerID)},
                "version": initialPostForWorkers.Version, // Use the initial version
            }

            results <- fmt.Sprintf("Worker %d: Attempting update (version: %d)...", workerID, initialPostForWorkers.Version)

            // Perform update
            err = updatePost(baseURL, postID, update)
            if err != nil {
                if strings.Contains(err.Error(), "409") {
                    results <- fmt.Sprintf("Worker %d: CONFLICT - Version mismatch (tried: %d)", workerID, initialPostForWorkers.Version)
                } else {
                    results <- fmt.Sprintf("Worker %d: ERROR - %v", workerID, err)
                }
                return
            }

            // If we get here, the update was successful
            updatedPost, err := getPost(baseURL, postID)
            if err != nil {
                results <- fmt.Sprintf("Worker %d: Error verifying update: %v", workerID, err)
                return
            }

            results <- fmt.Sprintf("Worker %d: SUCCESS - New version: %d, Title: %q", 
                workerID, updatedPost.Version, updatedPost.Title)
        }(i + 1)
    }

    // Start a goroutine to print results as they come in
    go func() {
        for result := range results {
            fmt.Println(result)
        }
    }()

    wg.Wait()
    close(results)

    // Show final state
    fmt.Println("\nFinal post state:")
    finalPost, err := getPost(baseURL, postID)
    if err != nil {
        fmt.Printf("Error getting final post: %v\n", err)
        return
    }
    printPost(finalPost)
}

func getPost(baseURL, id string) (*Post, error) {
    resp, err := http.Get(fmt.Sprintf("%s/posts/%s", baseURL, id))
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, string(body))
    }

    var response struct {
        Data Post `json:"data"` 
    }
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    return &response.Data, nil
}

func updatePost(baseURL, id string, update map[string]interface{}) error {
    jsonData, err := json.Marshal(update)
    if err != nil {
        return fmt.Errorf("error marshaling update: %w", err)
    }

    url := fmt.Sprintf("%s/posts/%s", baseURL, id)
    req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("error creating request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("update failed with status %s: %s", resp.Status, string(body))
    }

    return nil
}

func printPost(p *Post) {
    fmt.Printf("ID: %d\nTitle: %s\nVersion: %d\nUpdated: %v\nTags: %v\n", 
        p.ID, p.Title, p.Version, p.UpdatedAt, p.Tags)
}