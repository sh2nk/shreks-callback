package iriscallbackapi

import (
	"encoding/json"
	"io"
)

type signalRaw struct {
	IrisSignal
	Object json.RawMessage `json:"object"`
}

func unmarshalObject[O any](s signalRaw) (IrisSignal, error) {
	var obj O
	if err := json.Unmarshal(s.Object, &obj); err != nil {
		return IrisSignal{}, err
	}
	return IrisSignal{
		UserID:  s.UserID,
		Method:  s.Method,
		Object:  obj,
		Secret:  s.Secret,
		Message: s.Message,
	}, nil
}

// UnmarshalSignal - parse Iris signal data
func UnmarshalSignal(r io.Reader) (IrisSignal, error) {
	var sigRaw signalRaw
	if err := json.NewDecoder(r).Decode(&sigRaw); err != nil {
		return IrisSignal{}, err
	}

	switch sigRaw.Method {
	case "banExpired":
		return unmarshalObject[BanExpired](sigRaw)

	case "addUser":
		return unmarshalObject[AddUser](sigRaw)

	case "subscribeSignals":
		return unmarshalObject[SubscribeSignals](sigRaw)

	case "deleteMessages":
		return unmarshalObject[DeleteMessages](sigRaw)

	case "deleteMessagesFromUser":
		return unmarshalObject[DeleteMessagesFromUser](sigRaw)

	case "ignoreMessages":
		return unmarshalObject[IgnoreMessages](sigRaw)

	case "printBookmark":
		return unmarshalObject[PrintBookmark](sigRaw)

	case "forbiddenLinks":
		return unmarshalObject[ForbiddenLinks](sigRaw)

	case "sendSignal":
		return unmarshalObject[SendSignal](sigRaw)

	case "sendMySignal":
		return unmarshalObject[SendMySignal](sigRaw)

	case "hireApi":
		return unmarshalObject[HireAPI](sigRaw)

	case "toGroup":
		return unmarshalObject[ToGroup](sigRaw)

	case "banGetReason":
		return unmarshalObject[BanGetReason](sigRaw)

	default:
		return unmarshalObject[struct{}](sigRaw)
	}
}
