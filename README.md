# MDC Rules Copy Tool

[![Chinese README](https://img.shields.io/badge/README-中文-blue.svg)](README_CN.md)

This project generates Cursor MDC (Markdown Cursor) rule files from a structured JSON file containing library information.

**This project helps you make Cursor obey your instructions quickly!**

## Features

- Copy files according to mapping relationships
- Support filtering by rule category and type (e.g., only copy Python-related rules)
- Support listing all available rule categories and types
- Uses configuration files to define mapping relationships
- Alphabetically sorted output
- Embedded default configuration (no need for external config files)
- Automatic .cursor directory creation for Cursor compatibility
- Default target is the current directory

## Usage

### Prerequisites

- Go installed on your system
  - If you don't have Go installed, please visit [https://go.dev/doc/install](https://go.dev/doc/install) to install the latest version
  - Make sure to set up your GO_BIN environment variable correctly to use the `go install` command

### Installation

You can install directly from GitHub:

```bash
go install github.com/zhaolion/cursor-mdc-work@latest
```

After installation, you can use the `cursor-mdc-work` command directly to generate Cursor MDC (Markdown Cursor) rule files.

### Manual Compilation

If you prefer to compile manually:

```bash
go build -o cursor-mdc-work
```

### Examples

1. List all available rule categories and types using the embedded configuration:

```bash
cursor-mdc-work list
```

2. List types within a specific category:

```bash
cursor-mdc-work list --category backend-languages
```

3. Copy specific types of rule files to the current directory:

```bash
cursor-mdc-work copy --category backend-languages --types python,go
```

4. Copy all rule files in a specific category to a specified directory:

```bash
cursor-mdc-work copy --category backend-languages --target /path/to/destination
```

5. Copy all rule files to the current directory:

```bash
cursor-mdc-work copy
```

6. Use a custom configuration file:

```bash
cursor-mdc-work copy --config my-mapping.json
```

### Automatic .cursor Directory

When copying files, the tool will automatically check if the target path contains a `.cursor` directory:

- If the path already includes `.cursor`, files will be copied directly to that location
- If not, a new `.cursor` directory will be created within the target path, and files will be copied there

This ensures compatibility with Cursor's expected configuration structure.

## Configuration File Format

The configuration file format is JSON, as shown in the example below:

```json
{
  "mappings": {
    "frontend-frameworks": {
      "react": [
        "cursor/rules-mdc/react.mdc",
        "cursor/rules-mdc/javascript.mdc",
        "cursor/rules-mdc/typescript.mdc",
        "cursor/rules-mdc/html.mdc",
        "cursor/rules-mdc/css.mdc"
      ]
    },
    "backend-languages": {
      "python": ["cursor/rules-mdc/python.mdc", "cursor/rules-mdc/fastapi.mdc"]
    }
  }
}
```

If no configuration file is specified with the `--config` flag, the tool will use the embedded default configuration. If a configuration file is specified but does not exist, an example configuration file will be created at that location.

## Currently Supported Categories

The tool currently supports the following categories:

> For a complete list of categories and their types, run `cursor-mdc-work list`.

```
All rule categories and types:

ai-ml:
  - ai-agents
  - computer-vision
  - data-science
  - machine-learning
  - nlp

backend-languages:
  - go
  - java
  - node
  - other-backend
  - python
  - rust

build-tools:
  - bundlers
  - dependency-management

cloud-services:
  - aws
  - azure
  - gcp
  - other-cloud

databases:
  - nosql
  - sql

desktop-apps:
  - cross-platform
  - gui-frameworks

devops:
  - ci-cd
  - configuration-management
  - containerization
  - service-deployment

editors:
  - code-editors

frontend-frameworks:
  - angular
  - other-frontend
  - react
  - svelte
  - vue

frontend-ui:
  - ui-components
  - visualization

game-development:
  - game-engines
  - game-related

mobile:
  - flutter
  - other-mobile
  - react-native

testing:
  - api-testing
  - code-quality
  - testing-frameworks

utility-libraries:
  - api-integration
  - auth-security
  - documentation
  - monitoring
```
