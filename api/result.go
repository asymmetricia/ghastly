package api

type ResultMessage struct {
	Id      int
	Success bool
	Result  interface{}
	Error   ResultError `json:"error,omitempty"`
}

type ResultError struct {
	Code    string
	Message string
}

func (ResultMessage) Type() string { return "result" }

func init() { RegisterMessageType(ResultMessage{}) }
