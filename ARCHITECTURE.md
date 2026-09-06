# Lit Architecture and System Design

This document details the architectural principles, layer responsibilities, concurrency patterns, and AST analysis mechanisms implemented in Lit.

---

## 1. Architectural Strategy: Clean Architecture

Lit adheres strictly to **Clean Architecture** (Hexagonal Architecture / Ports and Adapters). The codebase is divided into isolated layers, guaranteeing that domain logic remains independent of external frameworks, AST parsers, or CLI drivers.

```
                      +---------------------------------------+
                      |               CLI Layer               |
                      |          (internal/cli)               |
                      +-------------------+-------------------+
                                          |
                                          v
                      +-------------------+-------------------+
                      |             Service Layer             |
                      |        (internal/service)             |
                      +-------------------+-------------------+
                                          |
                                          v
      +-----------------------------------+-----------------------------------+
      |                            Domain Layer                               |
      |                         (internal/domain)                             |
      |   - Entities & Models (FunctionData, Config, Alerts, Feedback)        |
      |   - Built-in Symbol Tables (O(1) Set Lookups)                        |
      |   - Port Interfaces (Analyzer, Exporter, ConfigAdapter)               |
      +-----------------------------------+-----------------------------------+
                                          ^
                                          |
                      +-------------------+-------------------+
                      |           Adapters Layer              |
                      |        (internal/adapters)            |
                      |   - Analysis (Tree-Sitter Parsers)    |
                      |   - Config (JSON Persistence)         |
                      |   - Exporter (SARIF, JSON, Markdown)  |
                      +---------------------------------------+
```

### Layer Breakdown

1. **Domain Layer (`internal/domain`)**:
   - Contains pure Go domain entities (`FunctionData`, `Config`, `Alerts`, `Feedback`).
   - Declares domain interface contracts (Ports: `Analyzer`, `Exporter`, `NodeManagement`).
   - Houses `BuiltinSymbolSet` (`map[string]map[string]struct{}`), providing $O(1)$ zero-heap-allocation lookup tables for language built-ins and dunder methods.
   - Zero external imports outside standard Go library packages.

2. **Service Layer (`internal/service`)**:
   - Implements application use cases (`ScanService`).
   - Manages bounded worker pool concurrency (`traverseFiles`) for parallel repository traversal.
   - Houses the custom Myers diff calculation engine (`GenerateUnifiedDiff`) for colorized ANSI diff previews.

3. **Adapters Layer (`internal/adapters`)**:
   - **Analysis Adapter (`adapters/analysis`)**: Implements language-specific Tree-Sitter queries for Python, JavaScript, JSX, Go, and Java. Calculates cyclomatic complexity and executes AST variable refactoring (`FileModifier`).
   - **Config Adapter (`adapters/config`)**: Manages `config.json` serialization and deserialization.
   - **Exporter Adapter (`adapters/exporter`)**: Translates analysis findings into Text, SARIF v2.1.0, JSON, and Markdown formats.

4. **CLI Layer (`internal/cli`)**:
   - Handles Cobra CLI commands (`scan`, `config`), flag parsing, and interactive terminal setup.

---

## 2. AST Analysis Engine and Tree-Sitter Integration

Lit uses `github.com/tree-sitter/go-tree-sitter` for fast, deterministic AST generation and query matching.

### 2.1 Scoped AST Queries vs Wildcard Capture
To eliminate false positives on third-party library calls (e.g., `os.path.exists()`, `json.dumps()`, `console.log()`), Tree-Sitter queries target **local declaration nodes only**:

- **Variables**: `assignment`, `variable_declarator`, `formal_parameters`, `for_statement`.
- **Functions**: `function_definition`, `function_declaration`, `method_definition`.

Member expression properties (`object.property`) are omitted from `@variable.name` capture definitions.

### 2.2 Built-in Protection ($O(1)$ Zero-Allocation Sets)
When an AST capture is matched, the node text is extracted via `node.Node.Utf8Text(source)` and evaluated against `domain.IsBuiltinSymbol(lang, symbol)`:

```go
// Zero-allocation set definition
var BuiltinSymbolSet = map[string]map[string]struct{}{
    "py":   { "print": {}, "len": {}, "__init__": {}, ... },
    "js":   { "console": {}, "document": {}, "window": {}, ... },
    "jsx":  { "useState": {}, "useEffect": {}, "React": {}, ... },
    "go":   { "main": {}, "append": {}, "make": {}, ... },
    "java": { "toString": {}, "equals": {}, "System": {}, ... },
}
```

Because Go `struct{}` values occupy 0 bytes of memory, set membership checking operates in $O(1)$ time with zero GC pressure.

---

## 3. Concurrency and Worker Pool Architecture

Repository traversal uses a bounded worker pool pattern to achieve maximum throughput on multi-core systems without exceeding file descriptor limits:

```
[ Directory Walk / Jobs Producer ] 
              |
              v (jobs channel: buffer 100)
    +---------+---------+---------+
    |         |         |         |
    v         v         v         v
 [Worker 1] [Worker 2] [Worker 3] [Worker 4]  (Default: NumCPU * 2)
    |         |         |         |
    +---------+---------+---------+
              |
              v (sync.Mutex protected maps)
    [ ScanService Data Collector ]
```

- Worker count defaults to $\max(4, \text{NumCPU} \times 2)$.
- Context cancellation (`context.Context`) is propagated across all worker goroutines to support immediate shutdown upon interrupt signals.
- Shared domain structures (`Feedback`, result maps) use `sync.Mutex` and `sync.Once` to ensure data-race-free operations.

---

## 4. Diff Engine and In-Place Refactoring

When `--dry-run` or `--fix` is executed:

1. `FileAnalyzer.FixFile` calls `FileModifier.ModifyVariableName(code)`.
2. Tree-Sitter extracts invalid variable identifiers matching `GetVarAppearancesQuery`.
3. Tokenization (`GetTokens`) splits identifiers into constituent words, converting them into the target naming convention (`refactorVarName`).
4. If `oldName == newName`, no string replacement occurs.
5. If `oldName != newName`, `ModifyVariableName` replaces the token string in memory and records column offset adjustments in a line cache.
6. `GenerateUnifiedDiff` executes Myers diff algorithm over original versus modified line arrays, producing colorized ANSI unified diffs (`@@ -L,C +L,C @@`).

---

## 5. Report Exporter Architecture

The exporter subsystem implements the **Strategy Pattern** via the `Exporter` port interface:

```go
type Exporter interface {
    Export(data map[string][]*domain.FunctionData) ([]byte, error)
}
```

Concrete implementations:
- `TextExporter`: Plain-text terminal output.
- `SarifExporter`: OASIS SARIF v2.1.0 standard JSON schema for GitHub Code Scanning integration.
- `JsonExporter`: Structured JSON representation of scan metrics and alerts.
- `MarkdownExporter`: GitHub-Flavored Markdown summary report.
