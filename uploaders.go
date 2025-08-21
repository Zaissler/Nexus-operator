package main

import (
	"strings"
)

// Uploader определяет контракт для загрузчиков разных форматов.
type Uploader interface {
	// Upload выполняет загрузку файла.
	Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error
	// IsSupported проверяет, подходит ли файл для данного загрузчика.
	IsSupported(filePath string) bool
}

// GetUploader возвращает нужную реализацию загрузчика по типу репозитория.
func GetUploader(repoType string) (Uploader, bool) {
	uploader, ok := uploaders[repoType]
	return uploader, ok
}

// --- Реализации для каждого формата ---

var uploaders = map[string]Uploader{
	"maven": &MavenUploader{},
	"npm":   &NpmUploader{},
	"raw":   &RawUploader{},
	"pypi":  &PypiUploader{},
	"nuget": &NugetUploader{},
	"helm":  &HelmUploader{},
	"yum":   &YumUploader{},
	"apt":   &AptUploader{},
}

type MavenUploader struct{}

func (u *MavenUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileMaven(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *MavenUploader) IsSupported(filePath string) bool {
	return strings.HasSuffix(filePath, ".jar") || strings.HasSuffix(filePath, ".pom")
}

type NpmUploader struct{}

func (u *NpmUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileNpm(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *NpmUploader) IsSupported(filePath string) bool {
	return strings.HasSuffix(filePath, ".tgz")
}

type RawUploader struct{}

func (u *RawUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileRaw(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *RawUploader) IsSupported(filePath string) bool {
	// Raw поддерживает любые файлы
	return true
}

type PypiUploader struct{}

func (u *PypiUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFilePypi(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *PypiUploader) IsSupported(filePath string) bool {
	return strings.HasSuffix(filePath, ".whl") || strings.HasSuffix(filePath, ".tar.gz")
}

type NugetUploader struct{}

func (u *NugetUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileNuget(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *NugetUploader) IsSupported(filePath string) bool {
	return strings.HasSuffix(filePath, ".nupkg")
}

type HelmUploader struct{}

func (u *HelmUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileHelm(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *HelmUploader) IsSupported(filePath string) bool {
	// Helm чарты и npm пакеты имеют одинаковое расширение.
	// Здесь мы доверяем флагу -repo-type, который указал пользователь.
	return strings.HasSuffix(filePath, ".tgz")
}

type YumUploader struct{}

func (u *YumUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileYum(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *YumUploader) IsSupported(filePath string) bool {
	return strings.HasSuffix(filePath, ".rpm")
}

type AptUploader struct{}

func (u *AptUploader) Upload(repoURL, repoName, filePath, importDir, username, password string, dryRun bool) error {
	return uploadFileApt(repoURL, repoName, filePath, importDir, username, password, dryRun)
}
func (u *AptUploader) IsSupported(filePath string) bool {
	return strings.HasSuffix(filePath, ".deb")
}
