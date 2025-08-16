package platform

import (
	"io"
	"net/http"
	"strings"
)

func CheckGitHub(httpClient *http.Client) (bool, error) {
	checks := []struct {
		url     string
		keyword string
	}{
		{"https://api.github.com", "current_user_url"},
		{"https://github.com", "github"},
		{"https://raw.githubusercontent.com", "githubusercontent"},
		{"https://githubusercontent.com", "githubusercontent"},
	}

	success := 0
	var lastErr error

	for _, c := range checks {
		ok, err := checkGitHubURL(httpClient, c.url, c.keyword)
		if ok {
			success++
		} else {
			lastErr = err // 记录最后一个错误
		}
	}
	// 至少两个子域可访问，认为 GitHub 可用
	return success >= 2, lastErr
}

func checkGitHubURL(httpClient *http.Client, url string, keyword string) (bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/122.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if strings.Contains(strings.ToLower(string(body)), strings.ToLower(keyword)) {
		return true, nil
	}
	return false, nil
}
