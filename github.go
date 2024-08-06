package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FileToGithub struct {
	Name       string
	RepoName   string
	Token      string
	OriginPath string
}

func (f *FileToGithub) Upload(fileContent []byte) (string, error) {

	uploadContent := map[string]string{
		"message": "Add file via Golang",
		"content": base64.StdEncoding.EncodeToString(fileContent),
		"branch":  "main",
	}

	contentBytes, err := json.Marshal(uploadContent)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", f.Name, f.RepoName, f.OriginPath)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(contentBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "token "+f.Token)
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to upload file: %s", body)
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to upload file: %s", body)
	}

	fileURL, ok := result["content"].(map[string]interface{})["html_url"].(string)
	if !ok {
		return "", errors.New("Failed to get file URL from GitHub response")
	}

	return fileURL, nil

}
