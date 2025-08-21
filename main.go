package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	repoURL := flag.String("repo-url", "", "Base URL of the Nexus repository (e.g., https://nexus.ac.astralinux.ru)")
	repoName := flag.String("repo-name", "", "Name of the repository (e.g., maven-test)")
	action := flag.String("action", "", "Action to perform: 'export' or 'import'")
	importDir := flag.String("import-dir", "", "Directory to import files from (required for import action)")
	repoType := flag.String("repo-type", "", "Type of repository: 'maven', 'npm', 'raw', 'pypi', 'nuget', 'helm', 'yum', 'apt'")
	username := flag.String("username", "", "Username for Nexus authentication (optional)")
	password := flag.String("password", "", "Password for Nexus authentication (optional)")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run without making any changes")
	numWorkers := flag.Int("workers", 10, "Number of concurrent workers for upload/download")
	flag.Parse()

	// Приоритет у флагов, но если они не заданы, используем переменные окружения.
	// Это удобно для CI/CD.
	if *username == "" {
		*username = os.Getenv("NEXUS_USERNAME")
	}
	if *password == "" {
		*password = os.Getenv("NEXUS_PASSWORD")
	}

	if *repoURL == "" || *repoName == "" || *action == "" || *repoType == "" {
		fmt.Println("Error: Please provide all required flags: -repo-url, -repo-name, -action, and -repo-type.")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	switch *action {
	case "export":
		err := ExportFiles(*repoURL, *repoName, *repoType, *dryRun, *numWorkers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Export failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Export completed successfully.")
	case "import":
		if *importDir == "" {
			fmt.Println("Please provide -import-dir flag for import action.")
			os.Exit(1)
		}
		err := ImportFiles(*repoURL, *repoName, *importDir, *repoType, *username, *password, *dryRun, *numWorkers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Import failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Import completed successfully.")
	default:
		fmt.Println("Invalid action. Use 'export' or 'import'.")
		os.Exit(1)
	}
}
