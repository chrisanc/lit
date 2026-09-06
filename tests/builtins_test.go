package tests

import (
	"CLI_App/internal/domain"
	"testing"
)

func TestBuiltinSymbols(t *testing.T) {
	tests := []struct {
		lang     string
		symbol   string
		expected bool
	}{
		// Python
		{"py", "__init__", true},
		{".py", "print", true},
		{"py", "len", true},
		{"py", "my_custom_var", false},
		// JavaScript
		{"js", "console", true},
		{".js", "document", true},
		{"js", "customFunction", false},
		// JSX
		{"jsx", "useState", true},
		{".jsx", "React", true},
		{"jsx", "myComponent", false},
		// Go
		{"go", "main", true},
		{".go", "append", true},
		{"go", "myStruct", false},
		// Java
		{"java", "toString", true},
		{".java", "System", true},
		{"java", "myField", false},
	}

	for _, tt := range tests {
		got := domain.IsBuiltinSymbol(tt.lang, tt.symbol)
		if got != tt.expected {
			t.Errorf("IsBuiltinSymbol(%q, %q) = %v; want %v", tt.lang, tt.symbol, got, tt.expected)
		}
	}
}

func TestIgnoredSymbols(t *testing.T) {
	ignored := []string{"_unused", "vendor_", "INTERNAL", "", "  "}

	tests := []struct {
		symbol   string
		expected bool
	}{
		{"_unused", true},
		{"_unusedVar", true},
		{"vendor_lib", true},
		{"INTERNAL_DATA", true},
		{"normalVar", false},
	}

	for _, tt := range tests {
		got := domain.IsIgnoredSymbol(tt.symbol, ignored)
		if got != tt.expected {
			t.Errorf("IsIgnoredSymbol(%q) = %v; want %v", tt.symbol, got, tt.expected)
		}
	}
}

func TestDualNamingConventionsConfig(t *testing.T) {
	// Fallback to legacy NamingConventionIndex if dual indices aren't set
	legacyCfg := &domain.Config{
		NamingConventionIndex: 4, // snake_case (index 4 -> domain.Conventions[3])
	}
	if legacyCfg.GetVariableConvention() != domain.Conventions[3] {
		t.Errorf("expected variable convention %q, got %q", domain.Conventions[3], legacyCfg.GetVariableConvention())
	}
	if legacyCfg.GetFunctionConvention() != domain.Conventions[3] {
		t.Errorf("expected function convention %q, got %q", domain.Conventions[3], legacyCfg.GetFunctionConvention())
	}

	// Dual naming conventions explicit configuration
	dualCfg := &domain.Config{
		NamingConventionIndex:         1,
		VariableNamingConventionIndex: 1, // camelCase (index 1 -> domain.Conventions[0])
		FunctionNamingConventionIndex: 4, // snake_case (index 4 -> domain.Conventions[3])
	}
	if dualCfg.GetVariableConvention() != domain.Conventions[0] {
		t.Errorf("expected variable convention %q, got %q", domain.Conventions[0], dualCfg.GetVariableConvention())
	}
	if dualCfg.GetFunctionConvention() != domain.Conventions[3] {
		t.Errorf("expected function convention %q, got %q", domain.Conventions[3], dualCfg.GetFunctionConvention())
	}
}
