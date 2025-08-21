# Nexus Operator

[Русская версия](README.ru.md)

Nexus Operator is a Go-based command-line tool designed for exporting and importing files from/to a Nexus Repository Manager. The program handles various repository formats and supports authentication for protected repositories.

Key Features
- Parallel Processing: Significantly speeds up export and import operations by processing multiple files concurrently.
- File Export: Downloads all artifacts from a specified Nexus repository, preserving the directory structure.
- File Import: Uploads files from a local directory to a Nexus repository.
- Format Support: Maven, npm, Raw, PyPI, NuGet, Helm, Yum, Apt.
- Authentication: Basic Auth for working with protected repositories.
- Safety: A `--dry-run` mode to preview operations without making any actual changes.

Requirements
- Go 1.20 or higher.
- Access to a Nexus Repository Manager instance.

## Installation

From Source:
1. `git clone https://github.com/your-repo/nexus-operator.git && cd nexus-operator`
2. `go build -o nexus-operator`

Using `go install`:
Make sure your `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.
`go install github.com/your-repo/nexus-operator@latest`

## Usage

### Exporting Files:
`./nexus-operator -repo-url=https://nexus.example.com -repo-name=maven-test -action=export -repo-type=maven`

The program will create a directory named after the repository (e.g., `maven-test`) and download all artifacts into it, preserving their structure.

### Importing Files:
`./nexus-operator -repo-url=https://nexus.example.com -repo-name=my-npm-repo -action=import -import-dir=./local-npm-packages -repo-type=npm -username=admin -password=admin123`

The program will upload all supported files from the specified directory to the Nexus repository.

### Dry Run:
`./nexus-operator -action=import -repo-type=maven -dry-run=true ...`

## Command-Line Flags
Flag            | Description                                      | Required | Example Value
------------------|--------------------------------------------------|----------|-----------------------------------
-repo-url         | Base URL of the Nexus Repository Manager         | Yes      | https://nexus.example.com
-repo-name        | Name of the repository                           | Yes      | maven-test
-action           | Action to perform: `export` or `import`          | Yes      | export
-import-dir       | Directory to import files from                   | For `import` | ./local-files
-repo-type        | Repository format: `maven`, `npm`, `raw`, etc.   | Yes      | maven
-username         | Username for Nexus authentication                | No       | admin
-password         | Password for Nexus authentication                | No       | admin123
-dry-run          | Show what would be done, without making changes  | No       | true

## Environment Variables
For convenience in CI/CD environments and for better security, credentials can be provided via environment variables. They have a lower priority than command-line flags.
- NEXUS_USERNAME: Username for authentication.
- NEXUS_PASSWORD: Password for authentication.

## Security
Important: Passing a password via the `-password` command-line flag can be insecure as it may be saved in your shell's history. For production use, consider using environment variables or other secure secret management methods.

## Contributing
Contributions are welcome! If you find a bug or have an idea for an improvement, please open an issue or submit a pull request in our repository.

## License
This project is distributed under the MIT License. See the `LICENSE` file for details.

## Troubleshooting
- "401 Unauthorized" error: Check if your username and password are correct. Ensure the user has the necessary permissions in Nexus.
- "Connection refused" or "timeout" error: Verify that the Nexus URL is accessible. A proxy or VPN configuration might be required.
- "No files found to upload": Make sure the `-repo-type` flag matches the files in the `-import-dir` directory. For example, for `-repo-type=maven`, the directory must contain `.jar` or `.pom` files.

### Author
Ilya Zaissler