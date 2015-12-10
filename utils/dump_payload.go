package utils

import (
	"encoding/json"
	"net/http"
)

type header struct {
	key   string
	value string
}

type payloadDumper struct {
	code    int
	headers []header
}

type Dumper interface {
	Dump(w http.ResponseWriter, payload interface{}) error
	Code(code int) Dumper
	ContentType(typ string) Dumper
	AddHeader(key, val string) Dumper
}

func (py payloadDumper) Dump(w http.ResponseWriter, payload interface{}) error {
	for _, h := range py.headers {
		w.Header().Add(h.key, h.value)
	}
	w.WriteHeader(py.code)

	return json.NewEncoder(w).Encode(payload)
}

func (py payloadDumper) Code(code int) Dumper {
	py.code = code
	return py
}

func (py payloadDumper) ContentType(typ string) Dumper {
	py.headers = append(py.headers, header{"ContentType", typ})
	return py
}

func (py payloadDumper) AddHeader(key, val string) Dumper {
	py.headers = append(py.headers, header{key, val})
	return py
}

var JSON = payloadDumper{}.Code(http.StatusOK).ContentType("application/json")
var Default = JSON
var DumpJSON = JSON.Dump
