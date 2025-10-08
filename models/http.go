package models

type ServerBlankOk struct {
	Ok bool `json:"ok"`
}

type ServerErrorResponse struct {
	Error string `json:"error"`
}

type RequestError struct {
	Code    int
	Message string
}
