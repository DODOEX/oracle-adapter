package bridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Copied from the example to use as a test fixture
type CryptoCompare struct{}

func (cc *CryptoCompare) Run(h *Helper) (interface{}, error) {
	r := make(map[string]interface{})
	err := h.HTTPCall(
		http.MethodGet,
		"https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,JPY,EUR",
		&r,
	)
	return r, err
}

func (cc *CryptoCompare) Opts() *Opts {
	return &Opts{
		Name:   "CryptoCompare",
		Lambda: true,
	}
}

func TestParseInterface_Map(t *testing.T) {
	p := map[string]interface{}{
		"alice": "bob",
		"carl":  "dennis",
	}
	json, err := ParseInterface(&p)
	assert.Nil(t, err)

	assert.Equal(t, json.Get("alice").String(), "bob")
	assert.Equal(t, json.Get("carl").String(), "dennis")
}

func TestParseInterface_String(t *testing.T) {
	p := "hello world"
	json, err := ParseInterface(&p)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", json.String())
}

type HelloWorld struct{}

func (tb *HelloWorld) Run(h *Helper) (interface{}, error) {
	return `{ "key": "hello world" }`, nil
}

func (tb *HelloWorld) Opts() *Opts {
	return &Opts{}
}

func TestNewServer_HelloWorld(t *testing.T) {
	b := &HelloWorld{}
	s := NewServer(b)
	assert.Nil(t, s.ldaBridge)
	assert.Equal(t, b, s.pathMap["/"])
}

type LambdaPath struct{}

func (tb *LambdaPath) Run(h *Helper) (interface{}, error) {
	return `{ "key": "hello world" }`, nil
}

func (tb *LambdaPath) Opts() *Opts {
	return &Opts{
		Lambda: true,
		Path:   "/path",
	}
}

func TestNewServer_LambdaPath(t *testing.T) {
	b := &LambdaPath{}
	s := NewServer(b)
	assert.Equal(t, b, s.ldaBridge)
	assert.Equal(t, b, s.pathMap["/path"])
}

func TestNewServer_Nil(t *testing.T) {
	s := NewServer()
	assert.Nil(t, s.ldaBridge)
	assert.Len(t, s.pathMap, 0)
}

func TestServer_Mux(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]interface{}
	}{
		{
			"input contains only id",
			map[string]interface{}{
				"id": "1234",
			},
		},
		{
			"input contains only jobRunId",
			map[string]interface{}{
				"jobRunId": "1234",
			},
		},
		{
			"input contains id and jobRunId",
			map[string]interface{}{
				"id":       "1234",
				"jobRunId": "1234",
			},
		},
	}

	for _, test := range tests {
		b := &HelloWorld{}
		mux := NewServer(b).Mux()

		pb, err := json.Marshal(test.in)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(pb))
		assert.Nil(t, err)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		body, err := ioutil.ReadAll(rr.Body)
		assert.Nil(t, err)
		json, err := Parse(body)
		assert.Nil(t, err)

		assert.Equal(t, "1234", json.Get("jobRunId").String())
		assert.Equal(t, "completed", json.Get("status").String())
	}
}

func TestServer_Mux_InvalidJSON(t *testing.T) {
	b := &HelloWorld{}
	mux := NewServer(b).Mux()

	p := `{ not json }`
	pb, err := json.Marshal(p)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(pb))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	body, err := ioutil.ReadAll(rr.Body)
	assert.Nil(t, err)
	json, err := Parse(body)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "errored", json.Get("status").String())
}

func TestServer_Mux_BadPath(t *testing.T) {
	b := &HelloWorld{}
	mux := NewServer(b).Mux()

	p := map[string]interface{}{
		"id": "1234",
	}
	pb, err := json.Marshal(p)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "/invalid", bytes.NewReader(pb))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	body, err := ioutil.ReadAll(rr.Body)
	assert.Nil(t, err)
	json, err := Parse(body)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "errored", json.Get("status").String())
}

func TestServer_Mux_BadMethod(t *testing.T) {
	b := &HelloWorld{}
	mux := NewServer(b).Mux()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	body, err := ioutil.ReadAll(rr.Body)
	assert.Nil(t, err)
	json, err := Parse(body)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "errored", json.Get("status").String())
}

type ReturnError struct{}

func (re *ReturnError) Run(h *Helper) (interface{}, error) {
	return `{}`, errors.New("error")
}

func (re *ReturnError) Opts() *Opts {
	return &Opts{}
}

func TestServer_Mux_ReturnError(t *testing.T) {
	b := &ReturnError{}
	mux := NewServer(b).Mux()

	p := map[string]interface{}{
		"id": "1234",
	}
	pb, err := json.Marshal(p)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(pb))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	body, err := ioutil.ReadAll(rr.Body)
	assert.Nil(t, err)
	json, err := Parse(body)
	assert.Nil(t, err)

	assert.Equal(t, "error", json.Get("error").String())
	assert.Equal(t, "errored", json.Get("status").String())
}

func TestServer_Mux_CryptoCompare(t *testing.T) {
	mux := NewServer(&CryptoCompare{}).Mux()

	p := map[string]interface{}{
		"id": "1234",
	}
	pb, err := json.Marshal(p)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(pb))
	assert.Nil(t, err)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	body, err := ioutil.ReadAll(rr.Body)
	assert.Nil(t, err)
	json, err := Parse(body)
	assert.Nil(t, err)

	assert.Equal(t, "1234", json.Get("jobRunId").String())

	data := json.Get("data").Map()
	_, ok := data["USD"]
	assert.True(t, ok)
	_, ok = data["JPY"]
	assert.True(t, ok)
	_, ok = data["EUR"]
	assert.True(t, ok)
}

func TestServer_Lambda_CryptoCompare(t *testing.T) {
	s := NewServer(&CryptoCompare{})

	r := &Result{}
	r.ID = "1234"

	obj, err := s.Lambda(r)
	assert.Nil(t, err)
	json, err := ParseInterface(obj)
	assert.Nil(t, err)

	assert.Equal(t, "1234", json.Get("jobRunId").String())

	data := json.Get("data").Map()
	_, ok := data["USD"]
	assert.True(t, ok)
	_, ok = data["JPY"]
	assert.True(t, ok)
	_, ok = data["EUR"]
	assert.True(t, ok)
}

func TestAuth_Header(t *testing.T) {
	a := NewAuth(AuthHeader, "API-KEY", "key")
	req, err := http.NewRequest(http.MethodGet, "http://test", nil)
	assert.Nil(t, err)
	a.Authenticate(req)

	assert.Equal(t, "key", req.Header.Get("API-KEY"))
}

// Copied from the example to use as a test fixture
type Timeout struct{}

func (t *Timeout) Run(h *Helper) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var r interface{}
	err := h.HTTPCallWithContext(
		ctx,
		http.MethodGet,
		"https://httpstat.us/200?sleep=4000",
		&r,
	)

	return r, err
}

func (t *Timeout) Opts() *Opts {
	return &Opts{
		Name:   "Timeout",
		Lambda: true,
	}
}

func TestServer_Lambda_Timeout(t *testing.T) {
	s := NewServer(&Timeout{})

	r := &Result{}
	r.ID = "1234"

	start := time.Now()

	_, err := s.Lambda(r)
	assert.Nil(t, err)

	diff := time.Since(start)
	assert.Less(t, int64(diff), int64(2*time.Second))
}
