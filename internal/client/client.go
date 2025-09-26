package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://alfa-leetcode-api.onrender.com"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type solvedResponse struct {
	SolvedProblem int `json:"solvedProblem"`
}

type ACSubmissionResponse struct {
	Count      int          `json:"count"`
	Submission []Submission `json:"submission"`
}

type Submission struct {
	Title         string `json:"title"`
	TitleSlug     string `json:"titleSlug"`
	Timestamp     string `json:"timestamp"` // UNIX timestamp as string
	StatusDisplay string `json:"statusDisplay"`
	Lang          string `json:"lang"`
	Difficulty    string `json:"difficulty"`
	ID            string `json:"id"`
}

func (c *Client) GetNumSolved(username string) (int, error) {
	url := fmt.Sprintf("%s/%s/solved", baseURL, username)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result solvedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.SolvedProblem, nil
}

func (c *Client) GetMostRecentAcceptedSubmission(username string) (*Submission, error) {
	// 1. Fetch most recent accepted submission
	submissionURL := fmt.Sprintf("%s/%s/acSubmission?limit=1", baseURL, username)
	resp, err := c.httpClient.Get(submissionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch most recent submission: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var acResp struct {
		Count      int          `json:"count"`
		Submission []Submission `json:"submission"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&acResp); err != nil {
		return nil, fmt.Errorf("failed to decode submission response: %w", err)
	}

	if len(acResp.Submission) == 0 {
		return nil, nil // no accepted submissions
	}

	sub := &acResp.Submission[0]

	if sub.TitleSlug == "" {
		return sub, nil
	}

	// 2. Fetch difficulty using titleSlug
	selectURL := fmt.Sprintf("%s/select?titleSlug=%s", baseURL, sub.TitleSlug)
	detailResp, err := c.httpClient.Get(selectURL)
	if err != nil {
		return sub, fmt.Errorf("failed to fetch question details: %w", err)
	}
	defer detailResp.Body.Close()

	if detailResp.StatusCode != http.StatusOK {
		return sub, fmt.Errorf("select API returned status %d", detailResp.StatusCode)
	}

	var detail struct {
		Difficulty string `json:"difficulty"`
		ID         string `json:"questionId"`
	}
	if err := json.NewDecoder(detailResp.Body).Decode(&detail); err != nil {
		return sub, fmt.Errorf("failed to decode question detail: %w", err)
	}

	sub.Difficulty = detail.Difficulty
	sub.ID = detail.ID
	return sub, nil
}
