package client

import (
	"fmt"
	"testing"
)

// TestClientWithRealAPI tests both API methods with a real username
func TestClientWithRealAPI(t *testing.T) {
	// Use a known LeetCode username - you can change this to any valid username
	testUsername := "banana_peel202" // This is a known valid username

	client := NewClient()

	fmt.Printf("=== Testing LeetCode API Client ===\n")
	fmt.Printf("Username: %s\n", testUsername)
	fmt.Printf("Base URL: %s\n\n", baseURL)

	// Test GetNumSolved
	fmt.Printf("1. Testing GetNumSolved...\n")
	solvedCount, err := client.GetNumSolved(testUsername)
	if err != nil {
		fmt.Printf("❌ GetNumSolved failed: %v\n", err)
		t.Errorf("GetNumSolved failed: %v", err)
	} else {
		fmt.Printf("✅ GetNumSolved successful!\n")
		fmt.Printf("   Solved problems: %d\n", solvedCount)
	}
	fmt.Printf("\n")

	// Test GetMostRecentAcceptedSubmission
	fmt.Printf("2. Testing GetMostRecentAcceptedSubmission...\n")
	submission, err := client.GetMostRecentAcceptedSubmission(testUsername)
	if err != nil {
		fmt.Printf("❌ GetMostRecentAcceptedSubmission failed: %v\n", err)
		t.Errorf("GetMostRecentAcceptedSubmission failed: %v", err)
	} else {
		fmt.Printf("✅ GetMostRecentAcceptedSubmission successful!\n")
		if submission == nil {
			fmt.Printf("   No accepted submissions found\n")
		} else {
			fmt.Printf("   Title: %s\n", submission.Title)
			fmt.Printf("   Title Slug: %s\n", submission.TitleSlug)
			fmt.Printf("   Difficulty: %s\n", submission.Difficulty)
			fmt.Printf("   Language: %s\n", submission.Lang)
			fmt.Printf("   Status: %s\n", submission.StatusDisplay)
			fmt.Printf("   Timestamp: %s\n", submission.Timestamp)
			fmt.Printf("   Problem ID: %s\n", submission.ID)
		}
	}

	fmt.Printf("\n=== Test Complete ===\n")
}
