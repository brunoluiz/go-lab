package app

type Envelope struct {
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}
