package helpers

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

func SendJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// Envelope the data in a wrapping tag to prevent subtle JSON security vulnerability.
	// https://haacked.com/archive/2008/11/20/anatomy-of-a-subtle-json-vulnerability.aspx/
	// Example of outcome:
	// 	{
	//		"data": {
	//			...data fields
	//		}
	//	}
	encoded := envelope{"data": data}

	res, err := json.Marshal(encoded)
	if err != nil {
		return err
	}

	// Add a newline to make it easier to view in terminals
	res = append(res, '\n')

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
	return nil
}
