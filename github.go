package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FileToGithub struct {
	Name       string
	RepoName   string
	Token      string
	OriginPath string
}

func (f *FileToGithub) Upload(fileContent []byte) error {

	uploadContent := map[string]string{
		"message": "Add file via Golang",
		"content": base64.StdEncoding.EncodeToString(fileContent),
	}

	contentBytes, err := json.Marshal(uploadContent)
	if err != nil {
		return err
	}

	// 创建HTTP请求
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", f.Name, f.RepoName, f.OriginPath)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(contentBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+f.Token)
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to upload file: %s", body)
	}

	fmt.Println("File uploaded successfully")
	return nil

}
