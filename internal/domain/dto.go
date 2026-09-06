package domain

// the 'dto' folder contains all the data transfer objects used mainly in the configuration set up.

// Config gets the value of the current index on the config file (JSON, etc.).
type Config struct {
	NamingConventionIndex int8   `json:"activeNamingConventionIndex"`
	Alerts                Alerts `json:"alerts"`
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
