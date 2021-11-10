package ports

type Payload struct {
	Status string `json:"status"`
}

type ErrorOut struct {
	Message   string                 `json:"message"`
	Operation string                 `json:"operation"`
	ID        string                 `json:"id"`
	Context   map[string]interface{} `json:"context"`
}

type ErrorPayload map[string][]string

type Common struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
