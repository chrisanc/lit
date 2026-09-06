package domain

// DTO structs for application configuration and settings.

type Config struct {
	NamingConventionIndex         int8     `json:"activeNamingConventionIndex"`
	VariableNamingConventionIndex int8     `json:"activeVariableNamingConventionIndex,omitempty"`
	FunctionNamingConventionIndex int8     `json:"activeFunctionNamingConventionIndex,omitempty"`
	IgnoredSymbols                []string `json:"ignoredSymbols,omitempty"`
	Alerts                        Alerts   `json:"alerts"`
}

type Alerts struct {
	Parameters FeedbackValues `json:"parameters"`
	Complexity FeedbackValues `json:"complexity"`
	MethodSize FeedbackValues `json:"method-length"`
}

type FeedbackValues struct {
	Info    uint `json:"info"`
	Warning uint `json:"warning"`
	Error   uint `json:"error"`
}

func (c *Config) GetVariableConventionIndex() int8 {
	if c.VariableNamingConventionIndex >= 1 && int(c.VariableNamingConventionIndex) <= len(Conventions) {
		return c.VariableNamingConventionIndex
	}
	if c.NamingConventionIndex >= 1 && int(c.NamingConventionIndex) <= len(Conventions) {
		return c.NamingConventionIndex
	}
	return 1
}

func (c *Config) GetFunctionConventionIndex() int8 {
	if c.FunctionNamingConventionIndex >= 1 && int(c.FunctionNamingConventionIndex) <= len(Conventions) {
		return c.FunctionNamingConventionIndex
	}
	if c.NamingConventionIndex >= 1 && int(c.NamingConventionIndex) <= len(Conventions) {
		return c.NamingConventionIndex
	}
	return 1
}

func (c *Config) GetVariableConvention() string {
	idx := c.GetVariableConventionIndex()
	return Conventions[idx-1]
}

func (c *Config) GetFunctionConvention() string {
	idx := c.GetFunctionConventionIndex()
	return Conventions[idx-1]
}
