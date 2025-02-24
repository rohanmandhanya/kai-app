package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"kai-app/api/models"
	"kai-app/arch/database"

	"gorm.io/gorm"
)

const githubAPIBaseURL = "https://api.github.com"
const maxRetries = 3 // 1 initial + 2 retries

// GitHubFileResponse represents the JSON response for a file's content
type GitHubFileResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// Downloader struct to hold repo info
type Downloader struct {
	RepoURL string
	Branch  string
	Client  *http.Client
	DB      *gorm.DB
}

// NewDownloader initializes a new Downloader with an HTTP client
func NewDownloader(repoURL, branch string) (*Downloader, error) {

	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}

	return &Downloader{
		RepoURL: repoURL,
		Branch:  branch,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		DB: db,
	}, nil
}

// fetchFileWithRetry fetches a file from GitHub, retrying up to 2 times on failure
func (d *Downloader) fetchFileWithRetry(ctx context.Context, filename string) ([]byte, error) {
	parts := strings.Split(strings.TrimSuffix(d.RepoURL, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GitHub repository URL")
	}

	user := parts[len(parts)-2]
	repo := parts[len(parts)-1]
	apiURL := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", githubAPIBaseURL, user, repo, filename, d.Branch)

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request for file %s: %v", filename, err)
		}

		resp, err := d.Client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: error fetching %s: %v", attempt, filename, err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					lastErr = fmt.Errorf("attempt %d: error reading response for %s: %v", attempt, filename, err)
				} else {
					return body, nil // Successful fetch
				}
			} else {
				body, _ := io.ReadAll(resp.Body)
				lastErr = fmt.Errorf("attempt %d: error fetching %s: %s (HTTP %d)", attempt, filename, string(body), resp.StatusCode)
			}
		}

		// Exponential backoff before retrying
		if attempt < maxRetries {
			fmt.Printf("Retrying %s... (attempt %d/%d)\n", filename, attempt, maxRetries)
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	return nil, lastErr // Return last encountered error
}

// DownloadFile downloads a single file with retry logic
func (d *Downloader) DownloadFile(ctx context.Context, filename string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	body, err := d.fetchFileWithRetry(ctx, filename)
	if err != nil {
		results <- err.Error()
		return
	}

	models.InsertMultipleScans(d.DB, body, filename)

	results <- fmt.Sprintf("Successfully downloaded: %s", filename)
}

// DownloadFilesConcurrently downloads multiple files using goroutines with retry logic
func (d *Downloader) DownloadFilesConcurrently(filenames []string) {
	var wg sync.WaitGroup
	results := make(chan string, len(filenames))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, filename := range filenames {
		wg.Add(1)
		go d.DownloadFile(ctx, filename, &wg, results)
	}

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Println(result)
	}
}
