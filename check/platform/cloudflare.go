package platform

import (
	"net/http"

	"log/slog"
)

func CheckCloudflare(httpClient *http.Client) (bool, error) {
    endpoints := []struct {
        url        string
        statusCode int
    }{
        {"https://gstatic.com/generate_204", 204},
        {"https://www.cloudflare.com/cdn-cgi/trace", 200},
        {"https://1.1.1.1", 200},
        {"https://www.cloudflare.com", 200}, // 可选，可能403
    }

    for _, ep := range endpoints {
        success, err := checkCloudflareEndpoint(httpClient, ep.url, ep.statusCode)
        if err == nil && success {
            slog.Debug("Cloudflare检测通过: " + ep.url)
            return true, nil
        }
        slog.Debug("Cloudflare检测失败: " + ep.url)
    }

    return false, nil
}


func checkCloudflareEndpoint(httpClient *http.Client, url string, statusCode int) (bool, error) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	// 添加请求头,模拟正常浏览器访问
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "close")

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Debug(err.Error())
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == statusCode, nil
}
