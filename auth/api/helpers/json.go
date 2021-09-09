package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// ReadJSON allows us to parse incoming JSON and control the potential errors that
// result in parsing something invalid. Some of the errors that come from decoding
// invalid JSON can give away too many details/information about our underlying API.
/* For Example - Consider the error message given the wront type (int rather than string) for a 'id' property.
$ curl -d '{"id": 123}' localhost:4000/v1/currentuser
{
    "error": "json: cannot unmarshal number into Go struct field .title of type string"
}
*/
func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific fields, then we include that in our error message to make it
		// easier for the client to debug
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message instead
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// A json.InvalidUnmarshalError will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error to our handler.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}
