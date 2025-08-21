package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadFileRaw(t *testing.T) {
	// 1. Настраиваем тестовый сервер, который будет имитировать Nexus
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос пришел правильным методом и по правильному пути
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/repository/test-raw/test-file.txt"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Проверяем содержимое файла
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		if string(body) != "hello world" {
			t.Errorf("Expected body 'hello world', got '%s'", string(body))
		}

		w.WriteHeader(http.StatusCreated) // Отвечаем успехом
	}))
	defer server.Close()

	// 2. Создаем временную директорию и файл для "загрузки"
	importDir := t.TempDir()
	filePath := filepath.Join(importDir, "test-file.txt")
	if err := os.WriteFile(filePath, []byte("hello world"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 3. Вызываем нашу функцию и проверяем результат
	err := uploadFileRaw(server.URL, "test-raw", filePath, importDir, "", "", false)
	if err != nil {
		t.Errorf("uploadFileRaw failed: %v", err)
	}
}

func TestUploadFileMaven(t *testing.T) {
	// 1. Настраиваем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что URL для Maven был сформирован правильно
		expectedPath := "/repository/test-maven/com/example/my-app/1.0/my-app-1.0.jar"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	// 2. Создаем временную файловую структуру, имитирующую Maven
	importDir := t.TempDir()
	mavenPath := filepath.Join(importDir, "com", "example", "my-app", "1.0")
	if err := os.MkdirAll(mavenPath, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	filePath := filepath.Join(mavenPath, "my-app-1.0.jar")
	if err := os.WriteFile(filePath, []byte("dummy jar content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 3. Вызываем функцию и проверяем результат
	err := uploadFileMaven(server.URL, "test-maven", filePath, importDir, "", "", false)
	if err != nil {
		t.Errorf("uploadFileMaven failed: %v", err)
	}
}
