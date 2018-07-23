package render

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

const (
	// DefaultFindHeaderIndex defines the maximum number of characters
	// to go through to find a generic XML header.
	DefaultFindHeaderIndex = 100
	// DefaultPrettyPrintXMLindent defines the number of spaces to pretty print a xml
	DefaultPrettyPrintXMLindent = "    "
	// DefaultPrettyPrintXMLPrefix defines the number of spaces to pretty print a xml
	DefaultPrettyPrintXMLPrefix = "  "
)

// PrettyPrintXML set XML encoding indent to DefaultPrettyPrintJSONIdent
func PrettyPrintXML() func(*XML) {
	return func(x *XML) {
		x.indentPrefix = DefaultPrettyPrintXMLPrefix
		x.indentValue = DefaultPrettyPrintXMLindent
	}
}

// XML encode the response as "application/xml" content type
// and implement the Renderer and APIRenderer interface.
type XML struct {
	w            http.ResponseWriter
	indentPrefix string
	indentValue  string
}

// NewXML returns a new XML responder instance.
func NewXML(w http.ResponseWriter, opts ...func(*XML)) *XML {
	j := &XML{w: w}
	for _, o := range opts {
		o(j)
	}
	return j
}

// Response marshals 'v' to XML, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func (x *XML) Response(code int, v interface{}) {
	b, err := xml.MarshalIndent(v, x.indentPrefix, x.indentValue)
	if err != nil {
		http.Error(x.w, err.Error(), http.StatusInternalServerError)
		return
	}

	x.w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	x.w.WriteHeader(code)

	// Try to find <?xml header in first DefaultFindHeaderIndex bytes
	// (just in case there're some XML comments).
	findHeaderUntil := len(b)
	if findHeaderUntil > DefaultFindHeaderIndex {
		findHeaderUntil = DefaultFindHeaderIndex
	}
	if !bytes.Contains(b[:findHeaderUntil], []byte("<?xml")) {
		// No header found. Print it out first.
		x.w.Write([]byte(xml.Header))
	}

	x.w.Write(b)
}

// Send sends a XML-encoded v in the body of a request with the 200 status code.
func (x *XML) Send(v interface{}) {
	x.Response(http.StatusOK, v)
}

// Created sends a XML-encoded v in the body of a request with the 201 status code.
func (x *XML) Created(v interface{}) {
	x.Response(http.StatusCreated, v)
}

// NoContent sends a v without no content with the 204 status code.
func (x *XML) NoContent() {
	x.w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a XML-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func (x *XML) BadRequest(err error) {
	s := http.StatusBadRequest
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	x.Response(http.StatusBadRequest, message)
}

// NotFound sends a XML-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func (x *XML) NotFound(err error) {
	s := http.StatusNotFound
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	x.Response(http.StatusNotFound, message)
}

// MethodNotAllowed sends a XML-encoded error response in the body of a request with the 405 status code.
// The response will contains the status 405 and error "Method Not Allowed".
func (x *XML) MethodNotAllowed(err error) {
	s := http.StatusMethodNotAllowed
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	x.Response(http.StatusMethodNotAllowed, message)
}

// InternalServerError sends a XML-encoded error response in the body of a request with the 500 status code.
// The response will contains the status 500 and error "Internal Server Error".
func (x *XML) InternalServerError(err error) {
	s := http.StatusInternalServerError
	message := NewHTTPError(err.Error(), http.StatusText(s), s)
	x.Response(http.StatusInternalServerError, message)
}