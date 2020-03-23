package api

type AuthRequiredMessage struct{}

func (AuthRequiredMessage) Type() string { return "auth_required" }

type AuthMessage struct {
	AccessToken string `json:"access_token,omitempty"`
	ApiPassword string `json:"api_password,omitempty"`
}

func (AuthMessage) Type() string { return "auth" }

type AuthInvalidMessage struct {
	Message string `json:"message"`
}

func (AuthInvalidMessage) Type() string { return "auth_invalid" }

type AuthOkMessage struct{}

func (AuthOkMessage) Type() string { return "auth_ok" }

func init() {
	RegisterMessageType(AuthRequiredMessage{}, AuthMessage{}, AuthInvalidMessage{}, AuthOkMessage{})
}
