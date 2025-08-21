package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Asset struct {
	DownloadURL string `json:"downloadUrl"`
	Path        string `json:"path"`
}

type SearchResult struct {
	Items             []Asset `json:"items"`
	ContinuationToken string  `json:"continuationToken"`
}

func fetchAssets(repoURL, repoName, continuationToken string) (SearchResult, error) {
	apiURL := fmt.Sprintf("%s/service/rest/v1/search/assets?repository=%s", repoURL, repoName)
	if continuationToken != "" {
		apiURL += "&continuationToken=" + continuationToken
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to fetch assets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResult{}, fmt.Errorf("failed to fetch assets: %s", resp.Status)
	}

	var searchResult SearchResult
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to decode response: %v", err)
	}

	return searchResult, nil
}
