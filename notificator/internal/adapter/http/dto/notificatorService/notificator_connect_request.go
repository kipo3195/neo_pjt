package notificatorService

import "encoding/json"

type NotificatorConnectRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
