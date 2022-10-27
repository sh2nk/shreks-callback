package iriscallbackapi

// Ping - contains "ping" signal data
type Ping struct {
}

// BanExpired - contains "banExpired" signal data
type BanExpired struct {
	UserID int    `json:"user_id"`
	Chat   string `json:"chat"`
	Reason string `json:"reason"`
}

// AddUser - contains "addUser" signal data
type AddUser struct {
	UserID int    `json:"user_id"`
	Chat   string `json:"chat"`
	Source string `json:"source"`
}

// SubscribeSignals - contains "subscribeSignals" signal data
type SubscribeSignals struct {
	Chat                   string `json:"chat"`
	ConverstationMessageID int    `json:"converstation_message_id"`
	Text                   string `json:"text"`
	FromID                 int    `json:"from_id"`
}

// DeleteMessages - contains "deleteMessages" signal data
type DeleteMessages struct {
	Chat    string `json:"chat"`
	LocalID []int  `json:"local_ids"`
	IsSpam  bool   `json:"is_spam"`
}

// DeleteMessagesFromUser - contains "deleteMessagesFromUser" signal data
type DeleteMessagesFromUser struct {
	Chat   string `json:"chat"`
	UserID int    `json:"user_id"`
	Amount int    `json:"amount"`
	IsSpam bool   `json:"is_spam"`
}

// IgnoreMessages - contains "ignoreMessages" signal data
type IgnoreMessages struct {
	Chat    string `json:"chat"`
	LocalID []int  `json:"local_ids"`
}

// PrintBookmark - contains "printBookmark" signal data
type PrintBookmark struct {
	Chat                   string `json:"chat"`
	ConverstationMessageID []int  `json:"converstation_message_id"`
	Description            string `json:"description"`
}

// ForbiddenLinks - contains "forbiddenLinks" signal data
type ForbiddenLinks struct {
	IgnoreMessages
}

// SendSignal - contains "sendSignal" signal data
type SendSignal struct {
	Chat                   string `json:"chat"`
	FromID                 int    `json:"from_id"`
	Value                  string `json:"value"`
	ConverstationMessageID int    `json:"converstation_message_id"`
}

// SendMySignal - contains "sendMySignal" signal data
type SendMySignal struct {
	SendSignal
}

// HireAPI - contains "hireApi" signal data
type HireAPI struct {
	Chat  string `json:"chat"`
	Price int    `json:"price"`
}

// ToGroup - contains "toGroup" signal data
type ToGroup struct {
	Chat    string `json:"chat"`
	GroupID int    `json:"group_id"`
	LocalID int    `json:"local_id"`
}

// BanGetReason - contains "banGetReason" signal data
type BanGetReason struct {
	Chat    string `json:"chat"`
	LocalID int    `json:"local_id"`
}
