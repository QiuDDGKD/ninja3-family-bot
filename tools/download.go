package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadByUrl(url string) (string, error) {
	// 发起 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: HTTP status %d", resp.StatusCode)
	}

	// 从 URL 中提取文件名
	fileName := filepath.Base(url)

	// 创建本地文件
	out, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// 将响应内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("File downloaded successfully: %s\n", fileName)
	return fileName, nil
}
