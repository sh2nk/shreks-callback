package iriscallbackapi

// IrisMessage - Iris message data model
type IrisMessage struct {
	ConverstationMessageID int `json:"converstation_message_id"`
	FromID                 int `json:"from_id"`
	Date                   int `json:"date"`
}

// IrisSignal - Iris signal data model
type IrisSignal struct {
	UserID  int    `json:"user_id"`
	Method  string `json:"method"`
	Object  any
	Secret  string      `json:"secret"`
	Message IrisMessage `json:"message"`
}
