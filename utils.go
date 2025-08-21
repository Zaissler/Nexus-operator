package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type downloadTask struct {
	URL      string
	FilePath string
}

type uploadTask struct {
	FilePath string
}

func ExportFiles(repoURL, repoName, repoType string, dryRun bool, numWorkers int) error {
	exportDir := repoName
	err := os.MkdirAll(exportDir, 0755) // Более безопасные права доступа
	if err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	var allAssets []Asset
	continuationToken := ""

	for {
		searchResult, err := fetchAssets(repoURL, repoName, continuationToken)
		if err != nil {
			return fmt.Errorf("error fetching assets: %w", err)
		}

		allAssets = append(allAssets, searchResult.Items...)

		if searchResult.ContinuationToken == "" {
			break
		}

		continuationToken = searchResult.ContinuationToken
	}

	if len(allAssets) == 0 {
		fmt.Println("No assets found in the repository.")
		return nil
	}

	var wg sync.WaitGroup
	total := len(allAssets)
	bar := progressbar.NewOptions(total,
		progressbar.OptionSetDescription("Exporting"),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "=", SaucerHead: ">", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
	)

	// --- Worker Pool для скачивания ---
	tasks := make(chan downloadTask, total)
	results := make(chan error, total)

	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go func() {
			defer wg.Done()
			for task := range tasks {
				err := downloadFile(task.URL, task.FilePath, dryRun)
				results <- err
				bar.Add(1)
			}
		}()
	}

	exporter := GetExporter(repoType)
	for _, asset := range allAssets {
		relativePath := exporter.GetLocalPath(asset.Path)
		filePath := filepath.Join(exportDir, relativePath)
		tasks <- downloadTask{URL: asset.DownloadURL, FilePath: filePath}
	}
	close(tasks)

	wg.Wait()
	close(results)

	failedCount := 0
	for err := range results {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка скачивания: %v\n", err)
			failedCount++
		}
	}

	if dryRun {
		fmt.Printf("[Dry Run] Было бы предпринято %d скачиваний.\n", len(allAssets))
	} else {
		fmt.Printf("Всего обработано файлов: %d, успешно: %d, с ошибками: %d\n", len(allAssets), len(allAssets)-failedCount, failedCount)
	}

	if failedCount > 0 {
		return fmt.Errorf("%d файлов не удалось скачать", failedCount)
	}

	return nil
}

func ImportFiles(repoURL, repoName, importDir, repoType, username, password string, dryRun bool, numWorkers int) error {
	var filesToUpload []string

	uploader, ok := GetUploader(repoType)
	if !ok {
		return fmt.Errorf("неподдерживаемый тип репозитория: %s", repoType)
	}

	// 1. Собираем все файлы для загрузки
	err := filepath.Walk(importDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if uploader.IsSupported(path) {
			filesToUpload = append(filesToUpload, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("ошибка при обходе директории импорта: %w", err)
	}

	if len(filesToUpload) == 0 {
		fmt.Println("Не найдено файлов для загрузки.")
		return nil
	}

	var wg sync.WaitGroup
	total := len(filesToUpload)
	bar := progressbar.NewOptions(total,
		progressbar.OptionSetDescription("Importing"),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "=", SaucerHead: ">", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
	)

	// --- Worker Pool для загрузки ---
	tasks := make(chan uploadTask, total)
	results := make(chan error, total)

	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go func() {
			defer wg.Done()
			for task := range tasks {
				uploadErr := uploader.Upload(repoURL, repoName, task.FilePath, importDir, username, password, dryRun)
				results <- uploadErr
				bar.Add(1)
			}
		}()
	}

	for _, filePath := range filesToUpload {
		tasks <- uploadTask{FilePath: filePath}
	}
	close(tasks)

	wg.Wait()
	close(results)

	failedCount := 0
	for err := range results {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка загрузки: %v\n", err)
			failedCount++
		}
	}

	if dryRun {
		fmt.Printf("[Dry Run] Было бы предпринято %d загрузок.\n", len(filesToUpload))
	} else {
		fmt.Printf("Всего обработано файлов: %d, успешно: %d, с ошибками: %d\n", len(filesToUpload), len(filesToUpload)-failedCount, failedCount)
	}

	if failedCount > 0 {
		return fmt.Errorf("%d файлов не удалось загрузить", failedCount)
	}

	return nil
}
