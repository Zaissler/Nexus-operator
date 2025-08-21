package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func uploadFileMaven(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	relativePath := strings.TrimPrefix(filePath, importDir+"/")
	if relativePath == filePath {
		return fmt.Errorf("file path does not start with '%s/': %s", importDir, filePath)
	}

	parts := strings.Split(relativePath, "/")
	if len(parts) < 4 {
		return fmt.Errorf("invalid file path for Maven repository: %s", relativePath)
	}

	groupID := strings.Join(parts[:len(parts)-3], "/")
	artifactID := parts[len(parts)-3]
	version := parts[len(parts)-2]
	fileName := parts[len(parts)-1]

	nexusPath := fmt.Sprintf("%s/%s/%s/%s", groupID, artifactID, version, fileName)
	apiURL := fmt.Sprintf("%s/repository/%s/%s", repoURL, repoName, nexusPath)

	resp, err := executeNexusRequest("PUT", apiURL, "application/octet-stream", file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response from Nexus: %s\n", string(responseBody))
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return nil
}

func uploadFileNpm(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	relativePath := strings.TrimPrefix(filePath, importDir+"/")
	if relativePath == filePath {
		return fmt.Errorf("file path does not start with '%s/': %s", importDir, filePath)
	}

	parts := strings.Split(relativePath, "/")
	fileName := parts[len(parts)-1]
	apiURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", repoURL, repoName)

	resp, err := executeMultipartUpload(apiURL, "npm.asset", fileName, file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response from Nexus: %s\n", string(responseBody))
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return nil
}

func uploadFileRaw(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Используем filepath.ToSlash для корректной работы на Windows
	relativePath := filepath.ToSlash(strings.TrimPrefix(filePath, importDir+string(filepath.Separator)))
	if relativePath == filePath {
		return fmt.Errorf("file path does not start with '%s/': %s", importDir, filePath)
	}

	apiURL := fmt.Sprintf("%s/repository/%s/%s", repoURL, repoName, relativePath)

	resp, err := executeNexusRequest("PUT", apiURL, "application/octet-stream", file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to upload file, status: %s", resp.Status)
	}
	return nil
}

func uploadFilePypi(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	apiURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", repoURL, repoName)

	resp, err := executeMultipartUpload(apiURL, "pypi.asset", filepath.Base(filePath), file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// PyPI upload returns 204 No Content on success
	if resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file, status: %s, body: %s", resp.Status, string(responseBody))
	}
	return nil
}

func uploadFileNuget(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// NuGet packages are uploaded to the root of the repository endpoint.
	// The trailing slash is important.
	apiURL := fmt.Sprintf("%s/repository/%s/", repoURL, repoName)

	resp, err := executeNexusRequest("PUT", apiURL, "application/octet-stream", file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// NuGet upload returns 201 Created on success.
	if resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file, status: %s, body: %s", resp.Status, string(responseBody))
	}
	return nil
}

func uploadFileHelm(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	apiURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", repoURL, repoName)

	resp, err := executeMultipartUpload(apiURL, "helm.asset", filepath.Base(filePath), file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file, status: %s, body: %s", resp.Status, string(responseBody))
	}
	return nil
}

func uploadFileYum(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	apiURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", repoURL, repoName)

	resp, err := executeMultipartUpload(apiURL, "yum.asset", filepath.Base(filePath), file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file, status: %s, body: %s", resp.Status, string(responseBody))
	}
	return nil
}

func uploadFileApt(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	if dryRun {
		// В режиме dry-run просто выходим, прогресс-бар покажет инкремент.
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	apiURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", repoURL, repoName)

	resp, err := executeMultipartUpload(apiURL, "apt.asset", filepath.Base(filePath), file, username, password)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file, status: %s, body: %s", resp.Status, string(responseBody))
	}
	return nil
}

// executeMultipartUpload создает и выполняет multipart/form-data запрос.
func executeMultipartUpload(apiURL, assetKey, fileName string, file io.Reader, username, password string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(assetKey, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file to form: %w", err)
	}
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	resp, err := executeNexusRequest("POST", apiURL, writer.FormDataContentType(), body, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to execute multipart request: %w", err)
	}

	return resp, nil
}

func executeNexusRequest(method, url, contentType string, body io.Reader, username, password string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	if username != "" && password != "" {
		auth := username + ":" + password
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}
