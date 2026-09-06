# Lit: High-Performance Code Analysis and Linting CLI

Lit is an advanced, high-performance static code analysis and refactoring command-line tool written in Go. Powered by Tree-Sitter Abstract Syntax Tree (AST) parsing and a concurrent worker pool, Lit enables software teams to analyze codebases, enforce naming conventions, detect architectural code smells, and preview automated refactoring changes safely.

---

## Key Features

- **AST-Powered Multi-Language Parsing**: Accurate syntactic analysis for Python, JavaScript, JSX, Go, and Java using native Tree-Sitter grammars.
- **Zero-Memory O(1) Built-in Symbol Protection**: Intelligent protection for built-in functions (`print`, `len`, `console.log`, `append`, `make`, etc.) and external library calls (`os.path.exists`, `json.dumps`), preventing unintended code breakage.
- **Dual Naming Conventions**: Independent rules for variable declarations versus function/method signatures (`camelCase`, `CamelCase`, `snake_case`, `CamelCase/camelCase`).
- **Automated Refactoring and Preview**:
  - `--dry-run`: Generates colorized, ANSI-formatted unified Git diffs showing proposed variable renamings without modifying files on disk.
  - `--fix`: Automatically applies safe in-place variable renamings across the repository.
- **Multi-Format Export Engine**: Export scan reports in Text, SARIF (Static Analysis Results Interchange Format), JSON, or Markdown for GitHub Actions and CI/CD pipelines.
- **Metrics and Code Smell Detection**: Evaluates Cyclomatic Complexity, Method Size (LOC), and Parameter Count against configurable thresholds.
- **Repository Composition Analysis**: `--loc` flag calculates lines of code and language distribution metrics.

---

## Supported Languages

| Language | Extension | AST Grammar Engine | Built-in Protection |
| :--- | :--- | :--- | :--- |
| **Go** | `.go` | `tree-sitter-go` | Yes |
| **Python** | `.py` | `tree-sitter-python` | Yes (includes dunders) |
| **JavaScript** | `.js` | `tree-sitter-javascript` | Yes |
| **JSX / React** | `.jsx` | `tree-sitter-javascript` | Yes (includes React hooks) |
| **Java** | `.java` | `tree-sitter-java` | Yes |

---

## Installation

### Requirements
- Go 1.22 or higher
- GCC or C compiler (required for Tree-Sitter C bindings)

### Build from Source
Clone the repository and build the binary:

```bash
git clone https://github.com/chrisanc/lit.git
cd lit
go build -o lit .
```

To make `lit` globally accessible, copy the binary to your PATH or create a symbol link:

```bash
sudo mv lit /usr/local/bin/
```

---

## Usage Guide

### 1. Basic Code Analysis
Scan the active repository using configured threshold alerts and naming rules:

```bash
lit scan
```

### 2. Preview Variable Refactorings (Dry Run)
Preview all proposed variable name changes as ANSI-colored unified diffs without altering files on disk:

```bash
lit scan --dry-run
```

### 3. Apply In-Place Variable Refactorings
Automatically refactor variable names across the repository to match the active convention:

```bash
lit scan --fix
```

### 4. Exporting Reports for CI/CD Pipelines
Export findings to SARIF for GitHub Security Code Scanning integration:

```bash
lit scan --format sarif --output report.sarif
```

Export findings to JSON or Markdown:

```bash
lit scan --format json --output report.json
lit scan --format markdown --output report.md
```

### 5. Repository Language Statistics
Calculate total lines of code and percentage breakdown by language:

```bash
lit scan --loc
```

### 6. Interactive Configuration Setup
Launch the terminal UI to configure variable and function naming conventions as well as alert thresholds:

```bash
lit config
```

---

## Configuration (`config.json`)

Lit stores repository configuration in `config.json` alongside the executable. Below is an annotated example configuration:

```json
{
  "activeNamingConventionIndex": 3,
  "activeVariableNamingConventionIndex": 1,
  "activeFunctionNamingConventionIndex": 4,
  "ignoredSymbols": ["vendor_", "_unused", "TEMP_"],
  "alerts": {
    "parameters": {
      "info": 5,
      "warning": 8,
      "error": 10
    },
    "complexity": {
      "info": 10,
      "warning": 15,
      "error": 20
    },
    "method-length": {
      "info": 120,
      "warning": 150,
      "error": 180
    }
  }
}
```

### Naming Convention Index Mapping
- `1`: `camelCase` (`LowerCamelCase`)
- `2`: `CamelCase` (`UpperCamelCase`)
- `3`: `CamelCase/camelCase` (Go-style: `CamelCase` for exported symbols, `camelCase` for unexported)
- `4`: `snake_case` (`SnakeCase`)

---

## Architecture Overview

Lit is architected around **Clean Architecture / Hexagonal Architecture** principles, enforcing strict decoupling between domain models, port interfaces, and adapter implementations:

```
lit/
├── cmd/               # CLI Entry Point
├── internal/
│   ├── domain/        # Pure Domain Entities, Models, Rules, and Built-in Tables
│   ├── service/       # Use Cases: Scanner Service, Worker Pool, Diff Engine
│   ├── adapters/
│   │   ├── analysis/  # Tree-Sitter AST Parsers and Language Rules
│   │   ├── config/    # JSON Configuration Persistence Adapter
│   │   └── exporter/  # Exporters: Text, SARIF, JSON, Markdown
│   └── cli/           # Cobra Command Handlers & Terminal UI
└── tests/             # Unit and Integration Test Suite
```

For comprehensive technical details on internal design, data flow, and concurrency model, refer to [ARCHITECTURE.md](ARCHITECTURE.md).

---

## Contributing

Contributions are welcome. Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for build instructions, code quality standards, and guidance on adding support for new programming languages.

---

## License

Distributed under the MIT License. See `LICENSE` for more information.
