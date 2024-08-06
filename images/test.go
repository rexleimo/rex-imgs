package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type FileToGithub struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	FilePath string `json:"file_path"`
	Token    string `json:"token"`
}

func (f *FileToGithub) Upload(fileContent []byte) (string, error) {
	repoURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", f.Owner, f.Repo, f.FilePath)

	endcodeContent := base64.StdEncoding.EncodeToString(fileContent)

	paylod := map[string]interface{}{
		"message": "upload file image",
		"content": endcodeContent,
	}

	jsonPayload, err := json.Marshal(paylod)
	if err != nil {
		return "", nil
	}

	headers := map[string]string{
		"Authorization": "token " + f.Token,
		"Content-Type":  "application/vnd.github+json",
	}

	req, err := http.NewRequest("PUT", repoURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file to github, status code: %d", resp.StatusCode)
	}

	var githubResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&githubResponse); err != nil {
		return "", err
	}

	content, ok := githubResponse["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse github response")
	}

	return content["html_url"].(string), nil
}
