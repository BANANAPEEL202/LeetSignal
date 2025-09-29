package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://leetcode-api-pied.vercel.app"

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

type TotalSubmission struct {
	Difficulty  string `json:"difficulty"`
	Count       int    `json:"count"`
	Submissions int    `json:"submissions"`
}

type SubmitStats struct {
	TotalSubmissionNum []TotalSubmission `json:"acSubmissionNum"`
}

type UserStats struct {
	SubmitStats SubmitStats `json:"submitStats"`
}

func (c *Client) GetNumSolved(username string) (int, error) {
	url := fmt.Sprintf("%s/user/%s", baseURL, username)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result UserStats
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Find the "All" entry
	for _, t := range result.SubmitStats.TotalSubmissionNum {
		if t.Difficulty == "All" {
			return t.Count, nil
		}
	}

	return 0, fmt.Errorf("total submissions not found")
}

func (c *Client) GetMostRecentAcceptedSubmission(username string) (*Submission, error) {
	// 1. Fetch most recent accepted submission
	submissionURL := fmt.Sprintf("%s/user/%s/submissions", baseURL, username)
	resp, err := c.httpClient.Get(submissionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch most recent submission: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var submissions []Submission
	if err := json.NewDecoder(resp.Body).Decode(&submissions); err != nil {
		return nil, fmt.Errorf("failed to decode submission response: %w", err)
	}

	if len(submissions) == 0 {
		return nil, nil // no accepted submissions
	}

	accepted_sub := &Submission{}
	for _, sub := range submissions {
		if sub.StatusDisplay == "Accepted" {
			accepted_sub = &sub
			break
		}
	}

	if accepted_sub.TitleSlug == "" {
		return accepted_sub, nil
	}

	// 2. Fetch difficulty using titleSlug
	selectURL := fmt.Sprintf("%s/problem/%s", baseURL, accepted_sub.TitleSlug)
	detailResp, err := c.httpClient.Get(selectURL)
	if err != nil {
		return accepted_sub, fmt.Errorf("failed to fetch question details: %w", err)
	}
	defer detailResp.Body.Close()

	if detailResp.StatusCode != http.StatusOK {
		return accepted_sub, fmt.Errorf("select API returned status %d", detailResp.StatusCode)
	}

	var detail struct {
		Difficulty string `json:"difficulty"`
		ID         string `json:"questionId"`
	}
	if err := json.NewDecoder(detailResp.Body).Decode(&detail); err != nil {
		return accepted_sub, fmt.Errorf("failed to decode question detail: %w", err)
	}

	accepted_sub.Difficulty = detail.Difficulty
	accepted_sub.ID = detail.ID
	return accepted_sub, nil
}
