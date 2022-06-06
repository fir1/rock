package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

/*
Don’t have to repeat yourself every time you respond to user, instead you can use some helper functions.
*/
func (s *Service) respond(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	v, ok := data.(ErrorResponse)
	if ok {
		w.WriteHeader(v.Status)
	} else {
		w.WriteHeader(status)
	}

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Could not encode in json", status)
		}
	}
}

// it does not read to the memory, instead it will read it to the given 'v' interface.
func (s *Service) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// it reads to the memory.
func (s *Service) readRequestBody(r *http.Request) ([]byte, error) {
	// Read the content
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			err := errors.New("could not read request body")
			return nil, err
		}
	}
	return bodyBytes, nil
}

// will place the body bytes back to the request body which could be read in subsequent calls on Handlers
// for example, you have more than 1 middleware and each of them need to read the body. If the first middleware read the body
// the second one won't be able to read it, unless you put the request body back.
func (s *Service) restoreRequestBody(r *http.Request, bodyBytes []byte) {
	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}
