package server

type ErrorResponse struct {
	Error string `json:"error"`
}

type HttpError struct {
	Code int
	Resp ErrorResponse
}
