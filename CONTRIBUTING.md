# Contributing to Lit

Thank you for your interest in contributing to Lit. This document provides technical instructions for setting up your development environment, executing tests, and submitting pull requests.

---

## Code of Conduct and Standards

Lit maintains high standards for software quality, performance, and architecture. All contributions must adhere to the following principles:

- **Clean Architecture Integrity**: Do not import `cli` or `adapters` inside the `domain` package.
- **Race-Free Concurrency**: All concurrent code must pass `go test -race ./...` without warnings.
- **Zero-Allocation Lookups**: Use `map[string]map[string]struct{}` for symbol set lookups.
- **No Unused Code or Swallowed Errors**: Handle all errors explicitly; avoid silent error suppressions.
- **Tone & Style**: Keep documentation, code comments, and commit messages serious, clear, and professional (no emojis).

---

## Development Setup

### Prerequisites
- Go 1.22 or higher
- GCC or a C compiler (required for CGo bindings of Tree-Sitter grammars)
- Git

### Build Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/chrisanc/lit.git
   cd lit
   ```

2. Download Go module dependencies:
   ```bash
   go mod download
   ```

3. Build the binary:
   ```bash
   go build -o lit .
   ```

---

## Running Tests

Run the full unit test suite:

```bash
go test -v ./tests/...
```

Run race detector testing across all packages:

```bash
go test -race ./...
```

Run benchmarks for AST parsing and scanner concurrency:

```bash
go test -bench=. ./tests/...
```

---

## Adding Support for a New Language

To add support for a new programming language (e.g., TypeScript or C++):

1. **Add Language Grammar Dependency**:
   Add the Tree-Sitter Go binding for the target language to `go.mod`.

2. **Register Built-in Symbols**:
   Update `BuiltinSymbolSet` in `internal/domain/builtins.go` with standard library functions, built-in types, and language keywords.

3. **Implement `NodeManagement` Adapter**:
   Create a new file in `internal/adapters/analysis/languages/` (e.g., `ts.go`). Implement:
   - `ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData, source []byte)`
   - `GetLanguageData() types.LanguageData`
   - `GetVarAppearancesQuery(varPattern string) string`
   - `GetFuncAppearancesQuery(funcPattern string) string`

4. **Wire Language into Analyzer**:
   Update `getLanguage(ext string)` in `internal/adapters/analysis/languages/analyzer.go` to instantiate the new language handler.

5. **Add Unit Test Cases**:
   Add test files and code samples to `tests/` verifying AST parsing, complexity scoring, and variable refactoring.

---

## Commit Guidelines

Use Conventional Commits format for clear git history:

- `feat(analysis): add TypeScript language parser support`
- `fix(domain): resolve race condition in feedback initialization`
- `docs(readme): update CLI usage guide and SARIF exporter docs`
- `perf(builtins): optimize O(1) set lookups`

---

## Submitting Pull Requests

1. Create a feature branch off `main`:
   ```bash
   git checkout -b feat/my-new-feature
   ```
2. Ensure all tests pass cleanly:
   ```bash
   go build ./... && go test -race ./...
   ```
3. Commit your changes and push to your fork.
4. Open a Pull Request detailing the changes, motivation, and verification steps.
