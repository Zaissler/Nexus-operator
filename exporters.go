package main

import "strings"

// Exporter определяет контракт для экспортеров разных форматов,
// позволяя обрабатывать специфичные для формата пути.
type Exporter interface {
	// GetLocalPath преобразует путь ассета из Nexus в локальный путь для сохранения.
	GetLocalPath(assetPath string) string
}

// GetExporter возвращает нужную реализацию экспортера по типу репозитория.
func GetExporter(repoType string) Exporter {
	if u, ok := exporters[repoType]; ok {
		return u
	}
	// Для всех остальных типов используется экспортер по умолчанию.
	return &DefaultExporter{}
}

var exporters = map[string]Exporter{
	"npm": &NpmExporter{},
}

// DefaultExporter - реализация по умолчанию, которая не меняет путь.
type DefaultExporter struct{}

func (e *DefaultExporter) GetLocalPath(assetPath string) string {
	return assetPath
}

// NpmExporter - реализация для npm, которая обрабатывает "/-/".
type NpmExporter struct{}

func (e *NpmExporter) GetLocalPath(assetPath string) string {
	return strings.Replace(assetPath, "/-/", "/", 1)
}
