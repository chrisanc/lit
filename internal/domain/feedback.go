package domain

// Feedback is used to send messages to the user when a something is >= value
type Feedback struct {
	configAdapter *Config
	messages      map[string][]Message
}

type Message struct {
	MinValue uint
	Message  string
}

func NewFeedback(configAdapter *Config) *Feedback {
	return &Feedback{
		configAdapter: configAdapter,
		messages:      make(map[string][]Message),
	}
}

func (f *Feedback) GetMessages() map[string][]Message {
	if len(f.messages) < 1 {
		f.setMessages()
	}

	return f.messages
}

func (f *Feedback) setMessages() {
	f.messages = map[string][]Message{
		"parameters": {
			Message{MinValue: f.configAdapter.Alerts.Parameters.Info, Message: "   INFO: The function takes more parameters than recommended.\n"},
			Message{MinValue: f.configAdapter.Alerts.Parameters.Warning, Message: "   WARNING: The function is dangerous and has many responsibilities.\n    You may want to make it cleaner/readable.\n"},
			Message{MinValue: f.configAdapter.Alerts.Parameters.Error, Message: "   POTENTIAL ERROR: The function takes too much parameters itself. You should simplify it.\n"},
		},
		"complexity": {
			Message{MinValue: f.configAdapter.Alerts.Complexity.Info, Message: "   INFO: The function it's a little bit complex.\n"},
			Message{MinValue: f.configAdapter.Alerts.Complexity.Warning, Message: "   WARNING: The function it's very complex itself. You may simplify or break it down.\n"},
			Message{MinValue: f.configAdapter.Alerts.Complexity.Error, Message: "   POTENTIAL ERROR: The function it's too complex and unreadable, taking too much decisions.\n    You must modify it and make it cleaner!.\n"},
		},
		"size": {
			Message{MinValue: f.configAdapter.Alerts.MethodSize.Info, Message: "   INFO: The functions might be long. You might want to make it smaller.\n"},
			Message{MinValue: f.configAdapter.Alerts.MethodSize.Warning, Message: "   WARNING: The function is longer than recommended. You should consider a refactor.\n"},
			Message{MinValue: f.configAdapter.Alerts.MethodSize.Error, Message: "   POTENTIAL ERROR: The function is too long. You must refactor now!.\n"},
		},
	}
}
