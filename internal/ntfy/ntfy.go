package ntfy

import (
	"bytes"
	"fmt"
	"leetsignal/internal/client"
	"leetsignal/internal/config"
	"net/http"
	"strings"
)

func SendLeetSignal(cfg config.Config, profile string, submission client.Submission) error {
	msg := fmt.Sprintf("%s just solved: \n%s. %s (%s)", profile, submission.ID, submission.Title, strings.ToLower(submission.Difficulty))
	return SendNtfy(cfg, "Leetcode Signal! 🎉", "1", msg, "https://leetcode.com/problems/"+submission.TitleSlug)
}

func SendErrorAlert(cfg config.Config, errMsg string) error {
	return SendNtfy(cfg, "Error", "3", errMsg, "")
}

func SendNtfy(cfg config.Config, title string, priority string, message string, link string) error {
	url := fmt.Sprintf("https://ntfy.sh/%s", cfg.NtfyTopic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return err
	}
	req.Header.Set("Title", title)
	req.Header.Set("Priority", priority)
	if link != "" {
		req.Header.Set("Click", link)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
