package gameserverapi

import (
	"bytes"
	"io"
	"net/http"

	"fmt"
)

func ResponsePlain(w http.ResponseWriter, r *http.Request, rd io.Reader) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	b := new(bytes.Buffer)

	_, err := b.ReadFrom(rd)
	if err != nil {
		panic(fmt.Errorf("bytes.(*Buffer).ReadFrom fn error: %s", err.Error()))
	}

	fmt.Fprint(w, b.String())
}

func ResponseXML(w http.ResponseWriter, r *http.Request, rd io.Reader) {
	w.Header().Set("Content-type", "application/xml; charset=utf-8")

	b := new(bytes.Buffer)

	_, err := b.ReadFrom(rd)
	if err != nil {
		panic(fmt.Errorf("bytes.(*Buffer).ReadFrom fn error: %s", err.Error()))
	}

	fmt.Fprint(w, b.String())
}

func ResponseJSON(
	w http.ResponseWriter,
	r *http.Request,
	isSuccess bool,
	responsePayload interface{},
	keyValues KV,
	errs ...*ErrorJSON) {
	b, err := jsonWithDebug(r.Context(), isSuccess, responsePayload, nil)
	if err != nil {
		panic(fmt.Errorf("answer.jsonWithDebug fn error: %s", err.Error()))
	}

	responseRawJSON(w, r, b)
}

func ResponseJSONWithDebug(
	w http.ResponseWriter,
	r *http.Request,
	isSuccess bool,
	responsePayload interface{},
	keyValues KV,
	errs ...*ErrorJSON) {
	b, err := jsonWithDebug(
		r.Context(), isSuccess, responsePayload, keyValues, errs...)
	if err != nil {
		panic(fmt.Errorf("answer.jsonWithDebug fn error: %s", err.Error()))
	}

	responseRawJSON(w, r, b)
}

func responseRawJSON(w http.ResponseWriter, r *http.Request, rd io.Reader) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")

	b := new(bytes.Buffer)

	_, err := b.ReadFrom(rd)
	if err != nil {
		panic(fmt.Errorf("bytes.(*Buffer).ReadFrom fn error: %s", err.Error()))
	}

	fmt.Fprint(w, b.String())
}
