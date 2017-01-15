package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpHandlerErrJson(t *testing.T) {
	assert := assert.New(t)
	body := strings.NewReader("")
	req := httptest.NewRequest("POST", "/", body)
	res := httptest.NewRecorder()
	handlerGenerateArchiveAndPng(res, req)
	assert.Equal(http.StatusInternalServerError, res.Result().StatusCode, "")
}

func TestHttpHandlerSuccJson(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test in short mode")
	}
	assert := assert.New(t)
	body := strings.NewReader(`{
    "name_ua": "Коломийський",
    "class": "провулок",
    "num": "22/1а",
		"left": "31",
		"right": "33"
}`)
	req := httptest.NewRequest("POST", "/", body)
	res := httptest.NewRecorder()
	go startLoop()
	handlerGenerateArchiveAndPng(res, req)
	assert.NotEqual(http.StatusInternalServerError, res.Result().StatusCode, "")
}
