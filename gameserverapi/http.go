package gameserverapi

import (
	"bytes"
	"io"
	"net/http"

	"fmt"
)

func ResponsePlain(rw http.ResponseWriter, req *http.Request, r io.Reader) {
	rw.Header().Set("Content-type", "text/plain; charset=utf-8")
	b := new(bytes.Buffer)

	_, err := b.ReadFrom(r)
	if err != nil {
		panic(fmt.Errorf("bytes.(*Buffer).ReadFrom fn error: %s", err.Error()))
	}

	fmt.Fprint(rw, b.String())
}

func ResponseXML(rw http.ResponseWriter, req *http.Request, r io.Reader) {
	rw.Header().Set("Content-type", "application/xml; charset=utf-8")

	b := new(bytes.Buffer)

	_, err := b.ReadFrom(r)
	if err != nil {
		panic(fmt.Errorf("bytes.(*Buffer).ReadFrom fn error: %s", err.Error()))
	}

	fmt.Fprint(rw, b.String())
}

func ResponseJSON(
	rw http.ResponseWriter,
	req *http.Request,
	isSuccess bool,
	responsePayload interface{},
	keyValues KV,
	errs ...*ErrorJSON) {
	b, err := jsonWithDebug(req.Context(), isSuccess, responsePayload, nil)
	if err != nil {
		panic(fmt.Errorf("answer.jsonWithDebug fn error: %s", err.Error()))
	}

	responseRawJSON(rw, req, b)
}

func ResponseJSONWithDebug(
	rw http.ResponseWriter,
	req *http.Request,
	isSuccess bool,
	responsePayload interface{},
	keyValues KV,
	errs ...*ErrorJSON) {
	b, err := jsonWithDebug(
		req.Context(), isSuccess, responsePayload, keyValues, errs...)
	if err != nil {
		panic(fmt.Errorf("answer.jsonWithDebug fn error: %s", err.Error()))
	}

	responseRawJSON(rw, req, b)
}

func responseRawJSON(rw http.ResponseWriter, req *http.Request, r io.Reader) {
	rw.Header().Set("Content-type", "application/json; charset=utf-8")

	b := new(bytes.Buffer)

	_, err := b.ReadFrom(r)
	if err != nil {
		panic(fmt.Errorf("bytes.(*Buffer).ReadFrom fn error: %s", err.Error()))
	}

	fmt.Fprint(rw, b.String())
}
