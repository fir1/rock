package server

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
	Status       int    `json:"-"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}
